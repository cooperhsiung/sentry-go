package notice

import (
	"gopkg.in/gomail.v2"
	"log"
)

func SendEmail(subject string, message string, receivers []string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "sentry_user@163.com")
	m.SetHeader("To", receivers...)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	d := gomail.NewDialer("smtp.163.com", 465, "sentry_user@163.com", "xjp12345")

	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}
}
