package rabbitMQService

import "log"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
}
