package rabbitMQService

import "github.com/streadway/amqp"

func Producer(ch *amqp.Channel, name string) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(name, false, false, false, false, nil)
	FailOnError(err, "Failed to declare a queue from producer")
	return q, err
}
