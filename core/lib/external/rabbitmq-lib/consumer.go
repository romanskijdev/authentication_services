package rabbitmqlib

import (
	"fmt"
	"log"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

const (
	ConsumerCountParam = 3
)

// ConsumerRabitMq представляет потребителя сообщений.
type ConsumerRabitMq struct {
	Conn        *ConnectionRabitMq // Соединение RabbitMQ.
	WorkerCount int                // Количество воркеров.
	channels    []*amqp091.Channel // Пул каналов.
	channelLock sync.Mutex         // Mutex для защиты доступа к пулу каналов.
	currentChan int                // Индекс текущего канала для round-robin.
}

func newConsumer(conn *ConnectionRabitMq, workerCount int) *ConsumerRabitMq {
	consumer := &ConsumerRabitMq{
		Conn:        conn,
		WorkerCount: workerCount,
		channels:    make([]*amqp091.Channel, workerCount),
	}

	// Инициализация пула каналов.
	for i := 0; i < workerCount; i++ {
		ch, err := conn.getChannel()
		if err != nil {
			log.Printf("🔴 error: failed to create channel %d: %v", i, err)
			continue
		}
		consumer.channels[i] = ch
	}
	return consumer
}

// Consume начинает потребление сообщений из очереди.
// Для каждого доставленного сообщения вызывается handler.
func (c *ConsumerRabitMq) Consume(
	queue string, // Имя очереди.
	consumerTag string, // Тег потребителя.
	handler func(amqp091.Delivery), // Функция обработки сообщения.
) error {
	// Запускаем воркеров для каждого канала.
	for i := 0; i < c.WorkerCount; i++ {
		go func(workerID int) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered in consumer worker %d: %v", workerID, r)
				}
			}()

			// Получаем канал из пула.
			ch := c.getChannel(workerID)
			if ch == nil {
				log.Printf("🔴 error: channel is nil for worker %d", workerID)
				return
			}

			// Настраиваем потребление.
			msgs, err := ch.Consume(
				queue,
				fmt.Sprintf("%s-%d", consumerTag, workerID),
				false, // auto-ack отключен
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				log.Printf("🔴 error: failed to start consuming on worker %d: %v", workerID, err)
				return
			}

			// Обрабатываем сообщения.
			log.Printf("✅ Started consumer worker %d", workerID)
			for d := range msgs {
				safeHandler(d, handler)
			}
			log.Printf("🔴 Consumer worker %d stopped", workerID)
		}(i)
	}
	return nil
}

// getChannel возвращает канал из пула, используя round-robin.
func (c *ConsumerRabitMq) getChannel(workerID int) *amqp091.Channel {
	c.channelLock.Lock()
	defer c.channelLock.Unlock()

	channelIndex := workerID % len(c.channels)
	if c.channels[channelIndex] == nil {
		log.Printf("🔴 error: channel %d is nil", channelIndex)
		return nil
	}
	return c.channels[channelIndex]
}

// safeHandler вызывает handler и обрабатывает возможную панику.
func safeHandler(d amqp091.Delivery, handler func(amqp091.Delivery)) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in handler: %v", r)
		}
	}()
	handler(d)
}

// ConsumeSimple запускает потребление сообщений из очереди и передаёт в handler только тело сообщения.
// Handler должен вернуть error, если обработка не удалась.
func (c *ConsumerRabitMq) ConsumeSimple(
	queue string, // Имя очереди.
	consumerTag string, // Тег потребителя.
	handler func([]byte) error, // Функция обработки (тело сообщения) с возвратом ошибки.
) error {
	return c.Consume(queue, consumerTag, func(d amqp091.Delivery) {
		if err := handler(d.Body); err != nil {
			log.Printf("🔴 Ошибка обработки сообщения: %v", err)
			if nackErr := d.Nack(false, true); nackErr != nil {
				log.Printf("🔴 Ошибка отклонения сообщения: %v", nackErr)
			}
			return
		}

		if err := d.Ack(false); err != nil {
			log.Printf("🔴 Ошибка подтверждения сообщения: %v", err)
		}
	})
}
