package rabbitmqlib

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// DeclareQueue декларирует exchange и очередь, а затем выполняет binding.
func (c *ConnectionRabitMq) DeclareQueue(queueName, exchange, bindingKey string) error {
	ch, err := c.getChannel()
	if err != nil {
		return fmt.Errorf("🔴 error: failed to create channel for declaration: %w", err)
	}
	defer func() {
		if err := ch.Close(); err != nil {
			logrus.Errorf("🔴 error: failed to close declaration channel: %v", err)
		} else {
			logrus.Info("⚪ Declaration channel successfully closed")
		}
	}()

	if err = ch.ExchangeDeclare(
		exchange, // Имя exchange.
		"direct", // Тип exchange.
		true,     // Durable.
		false,    // Auto-deleted.
		false,    // Internal.
		false,    // No-wait.
		nil,      // Arguments.
	); err != nil {
		return fmt.Errorf("🔴 error: failed to declare exchange: %w", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // Имя очереди.
		true,      // Durable.
		false,     // Delete when unused.
		false,     // Exclusive.
		false,     // No-wait.
		nil,       // Arguments.
	)
	if err != nil {
		return fmt.Errorf("🔴 error: failed to declare queue: %w", err)
	}

	if err = ch.QueueBind(
		q.Name,     // Имя очереди.
		bindingKey, // Routing key.
		exchange,   // Exchange.
		false, nil,
	); err != nil {
		return fmt.Errorf("🔴 error: failed to bind queue: %w", err)
	}

	// TODO: Логирование успешного выполнения
	logrus.Infof("⚪ Successfully declared and bound queue: queueName=%s, exchange=%s, bindingKey=%s", queueName, exchange, bindingKey)
	return nil
}
