package rabbitmqlib

import (
	errm "authentication_service/core/errmodule"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const (
	PublisherCountParam = 3
)

// PublisherRabitMq представляет объект для публикации сообщений с использованием пула каналов.
type PublisherRabitMq struct {
	Conn *ConnectionRabitMq    // Соединение RabbitMQ.
	pool chan *amqp091.Channel // Пул каналов.
}

func newPublisher(conn *ConnectionRabitMq, poolSize int) (*PublisherRabitMq, error) {
	p := &PublisherRabitMq{
		Conn: conn,
		pool: make(chan *amqp091.Channel, poolSize),
	}
	// Инициализация пула каналов.
	for i := 0; i < poolSize; i++ {
		ch, err := conn.getChannel()
		if err != nil {
			closePublisherPool(p.pool)
			return nil, fmt.Errorf("🔴 error: failed to create channel: %w", err)
		}
		if err := ch.Confirm(false); err != nil {
			_ = ch.Close()
			closePublisherPool(p.pool)
			return nil, fmt.Errorf("🔴 error: failed to enable publisher confirms: %w", err)
		}
		p.pool <- ch
	}
	return p, nil
}

// closePublisherPool закрывает все каналы из пула.
func closePublisherPool(pool chan *amqp091.Channel) {
	close(pool)
	for ch := range pool {
		_ = ch.Close()
	}
}

// PublishJSON отправляет данные, сериализованные в JSON, в указанный exchange с routingKey.
// При ошибке публикации канал переинициализируется.
func (p *PublisherRabitMq) PublishJSON(exchange string, routingKey string, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("🔴 error: failed to marshal data: %w", err)
	}

	var ch *amqp091.Channel
	select {
	case ch = <-p.pool:
		// Канал успешно получен.
	case <-time.After(5 * time.Second):
		return fmt.Errorf("🔴 error: timeout while receiving a channel")
	}

	// Пытаемся опубликовать сообщение максимум два раза.
	var pubErr error
	for attempt := 0; attempt < 2; attempt++ {
		pub := amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "application/json",
			Body:         body,
		}
		pubErr = ch.Publish(exchange, routingKey, false, false, pub)
		if pubErr == nil {
			break
		}
		log.Printf("🔴 error: failed to publish message, reinitializing channel (attempt %d): %v", attempt+1, pubErr)
		newCh, dialErr := p.Conn.getChannel()
		if dialErr != nil {
			pubErr = fmt.Errorf("🔴 error: failed to create a new channel: %w", dialErr)
			break
		}
		if confirmErr := newCh.Confirm(false); confirmErr != nil {
			_ = newCh.Close()
			pubErr = fmt.Errorf("🔴 error: failed to enable publisher confirms on new channel: %w", confirmErr)
			break
		}
		ch = newCh
	}
	// Возвращаем канал в пул.
	p.pool <- ch

	if pubErr != nil {
		return fmt.Errorf("🔴 error: failed to publish message after channel reinitialization: %w", pubErr)
	}
	return nil
}

// PublishMessage отправляет одно сообщение в указанную очередь
func PublishMessage(connection *ConnectionRabitMq, exchange, routingKey string, message interface{}) *errm.Error {
	// Создаем publisher
	publisher, err := connection.NewPublisher()
	if err != nil {
		return errm.NewError("failed to create new publisher", err)
	}

	// Публикуем сообщение
	if err := publisher.PublishJSON(exchange, routingKey, message); err != nil {
		return errm.NewError("failed to publish message", err)
	}

	return nil
}
