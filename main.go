package main

import (
	"errors"
	"fmt"
	"github.com/vleedev/smtp-relay-rabbitmq/queue"
	"github.com/vleedev/smtp-relay-rabbitmq/smtp"
	"github.com/vleedev/smtp-relay-rabbitmq/utils"
	"os"
	"strconv"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		utils.ErrFatal(err)
	}
	// Init queue
	q := queue.Init(
		os.Getenv("QUEUE_NAME"),
		os.Getenv("RABBITMQ_URL"),
		)
	defer q.Connection.Close()
	defer q.Channel.Close()
	// Init email client
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	utils.ErrFatal(err)
	smtpConfig := smtp.Config{
		Hostname: os.Getenv("SMTP_HOSTNAME"),
		Port:     smtpPort,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		DefaultEmail: os.Getenv("SMTP_DEFAULT_EMAIL"),
	}
	smMail := smtp.Init(&smtpConfig)
	// Test smtp service
	_ = smMail.NewMail(smtp.MailTemplate{
		Subject:    "Welcome to smtp-relay-rabbitmq",
		BodyType:   "text/html",
		Body:       fmt.Sprintf("<html><body><p>This one is a test email from smtp-relay-rabbitmq<br />Hostname: %s</p></body></html>", hostname),
		Attachment: nil,
	})
	if err := smMail.Send(); err != nil {
		utils.ErrFatal(errors.New("please check your smtp configuration"))
	}
	// Test send queue
	mailTemp := smtp.MailTemplate{
		Subject:    "smtp-relay-rabbitmq queue",
		BodyType:   "text/html",
		Body:       fmt.Sprintf("<html><body><p><b>This email is from the queue</b><br />Hostname: %s</p></body></html>", hostname),
		Attachment: []string{"https://i.imgur.com/UbUQWHO.jpeg"},
	}
	q.Send(mailTemp)
	// Consume service
	q.Consume(smMail)
}
