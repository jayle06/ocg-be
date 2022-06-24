package rabbitMQService

import "github.com/streadway/amqp"

func Receiver(ch *amqp.Channel, queue amqp.Queue) (result <-chan amqp.Delivery) {
	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	FailOnError(err, "Failed to register a consumer")
	return msgs
}
