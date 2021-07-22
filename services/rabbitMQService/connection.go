package rabbitMQService

import "github.com/streadway/amqp"

const (
	user = "guest"
	pass = "guest"
	host = "localhost"
	port = "5672"
)

func Connect() *amqp.Connection {
	conn, err := amqp.Dial("amqp://" + user + ":" + pass + "@" + host + ":" + port)
	FailOnError(err, "Fail to connect to RabbitMQ")
	return conn
}
