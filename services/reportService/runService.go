package reportService

import (
	"fmt"
	"github.com/final-project/services/rabbitMQService"
	"github.com/robfig/cron/v3"
	"strconv"
	"time"
)

func RunReportService() {
	conn := rabbitMQService.Connect()
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	q, err := rabbitMQService.Producer(ch, "order")
	rabbitMQService.FailOnError(err, "can not open queue")

	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 0h1m0s", func() {
		for _, id := range getOrdersPerMinute() {
			sendMessage(q, ch, id)
		}
	})

	c.AddFunc("0 0 19 * * * ", func() {
		msgs := rabbitMQService.Receiver(ch, q)
		i := 1
		var data [][5]string
		go func() {
			for d := range msgs {
				id := string(d.Body)
				//fmt.Println(id)
				order := getDataByID(id)
				fmt.Println("order", order)
				for _, value := range order {
					value[0] = strconv.Itoa(i)
					data = append(data, value)
					i++
				}
			}
		}()
		time.Sleep(1 * time.Second)
		sendEmailToAdmin(parseDataToCSV(data))
	})
	c.Start()
	defer c.Stop()
	forever := make(chan bool)
	<-forever
}
