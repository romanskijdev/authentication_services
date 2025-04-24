package rabbitmqlib

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var (
	clientInstance *ConnectionRabitMq
	once           sync.Once
)

// ConnectParams содержит параметры для установления подключения к RabbitMQ.
type ConnectParams struct {
	Host     string // Адрес хоста RabbitMQ.
	Port     int    // Порт подключения.
	Username string // Имя пользователя.
	Password string // Пароль подключения.
}

// ConnectionRabitMq оборачивает соединение RabbitMQ.
type ConnectionRabitMq struct {
	Conn *amqp091.Connection
}

// ConnectRabitMq устанавливает соединение с RabbitMQ, используя параметры из ConnectParams.
// При неудачных попытках выполняется серия из 3-х попыток с интервалом 30 секунд.
// После Dial выполняется пинг-проверка (попытка открыть и закрыть канал).
// В случае успеха выводится сообщение с эмодзи.
func ConnectRabitMq(params ConnectParams) (*ConnectionRabitMq, error) {
	var err error
	once.Do(func() {
		const maxAttempts = 3
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%d/", params.Username, params.Password, params.Host, params.Port)
			conn, dialErr := amqp091.Dial(amqpURL)
			if dialErr != nil {
				err = fmt.Errorf("🔴 error: dial failed: %w", dialErr)
			} else {
				// Ping-проверка: открываем канал и сразу его закрываем.
				ch, pingErr := conn.Channel()
				if pingErr != nil {
					err = fmt.Errorf("🔴 error: ping check failed: %w", pingErr)
					_ = conn.Close()
				} else {
					_ = ch.Close()
					log.Println("✅ Successfully connected to RabbitMQ and ping check passed!")
					clientInstance = &ConnectionRabitMq{Conn: conn}
					err = nil
					break
				}
			}
			log.Printf("⏰ Attempt %d/%d: failed to connect to RabbitMQ: %v", attempt, maxAttempts, err)
			if attempt < maxAttempts {
				log.Printf("⏳ Retrying in 30 seconds...")
				time.Sleep(30 * time.Second)
			}
		}
		if clientInstance == nil && err != nil {
			err = fmt.Errorf("🔴 error: failed to connect after %d attempts: %w", maxAttempts, err)
		}
	})
	return clientInstance, err
}

// Close закрывает соединение.
func (c *ConnectionRabitMq) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

// getChannel возвращает новый канал из соединения.
func (c *ConnectionRabitMq) getChannel() (*amqp091.Channel, error) {
	if c.Conn == nil {
		return nil, fmt.Errorf("connection is nil")
	}
	return c.Conn.Channel()
}
