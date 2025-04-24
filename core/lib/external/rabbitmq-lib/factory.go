package rabbitmqlib

// NewConsumer создаёт нового ConsumerRabitMq.
func (c *ConnectionRabitMq) NewConsumer() *ConsumerRabitMq {
	return newConsumer(c, ConsumerCountParam)
}

// NewPublisher создаёт нового PublisherRabitMq.
func (c *ConnectionRabitMq) NewPublisher() (*PublisherRabitMq, error) {
	return newPublisher(c, PublisherCountParam)
}
