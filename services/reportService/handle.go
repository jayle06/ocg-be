package reportService

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/final-project/database"
	"github.com/final-project/services/rabbitMQService"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func getOrdersPerMinute() []string {
	db := database.Connect()
	defer db.Close()
	var id []string
	now := time.Now().Add(-time.Minute)
	rows, err := db.Query("SELECT id FROM orders WHERE created_at >= ?", now)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var i string
		rows.Scan(&i)
		id = append(id, i)
	}
	return id
}
func sendMessage(queue amqp.Queue, ch *amqp.Channel, message string) {
	err := ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	rabbitMQService.FailOnError(err, "Failed to publish message")
}

func getAdminEmail() []string {
	db := database.Connect()
	defer db.Close()

	rows, _ := db.Query("SELECT email FROM users u JOIN roles r on u.role_id = r.id WHERE r.name = 'ADMIN'")
	var emails []string
	for rows.Next() {
		var email string
		rows.Scan(&email)
		emails = append(emails, email)
	}
	return emails
}

func getDataByID(id string) [][5]string {
	db := database.Connect()
	defer db.Close()
	rows, _ := db.Query("SELECT p.name, p.price, o.created_at, ot.quantity "+
		"FROM orders o  "+
		"JOIN order_items ot ON o.id = ot.order_id "+
		"JOIN products p ON p.id = ot.product_id "+
		"WHERE o.id = ?", id)
	var rowData [5]string
	for i := 0; i < 5; i++ {
		rowData[i] = " "
	}
	var rowsData [][5]string
	for rows.Next() {
		var name, price, created, quantity string
		rows.Scan(&name, &price, &created, &quantity)
		rowData[0] = ""
		rowData[1] = name
		rowData[2] = price
		rowData[3] = created
		rowData[4] = quantity
		rowsData = append(rowsData, rowData)
		fmt.Println(rowsData)
	}
	return rowsData
}

func parseDataToCSV(data [][5]string) string {
	filename := "uploads/report/report" + time.Now().Format("2006-Jan-02") + ".csv"
	csvFile, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}

	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Write([]string{"STT", "tên sản phẩm", "Giá tiền", "Thời gian", "số lượng"})
	for _, empRow := range data {
		csvwriter.Write(empRow[:])
	}
	csvwriter.Flush()
	csvFile.Close()
	return filename
}

func sendEmailToAdmin(fileName string) {
	m := mail.NewV3Mail()

	from := mail.NewEmail("Nam Dương", "nhuphongkiwi@gmail.com")
	content := mail.NewContent("text/html", "<p> Daily Report : %day%</p>")

	m.SetFrom(from)
	m.AddContent(content)

	personalization := mail.NewPersonalization()

	for _, email := range getAdminEmail() {
		to := mail.NewEmail("Admin Kayle Shop", email)
		personalization.AddTos(to)
	}

	personalization.SetSubstitution("%day%", time.Now().Format("2006-Jan-12"))
	personalization.Subject = "Daily Report" + time.Now().Format("2006-Jan-21")

	m.AddPersonalizations(personalization)

	a_csv := mail.NewAttachment()
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(dat))
	a_csv.SetContent(encoded)
	a_csv.SetType("text/csv")
	a_csv.SetFilename(fileName)
	a_csv.SetDisposition("attachment")
	m.AddAttachment(a_csv)
	request := sendgrid.GetRequest("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", "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
