package main

import (
	"github.com/vleedev/smtp-relay-rabbitmq/queue"
	"github.com/vleedev/smtp-relay-rabbitmq/smtp"
	"github.com/vleedev/smtp-relay-rabbitmq/utils"
	"os"
	"strconv"
)

func main() {
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
	smMail.Test()
	// Test send queue
	mailTemp := smtp.MailTemplate{
		Subject:    "smtp-relay-rabbitmq queue",
		BodyType:   "text/html",
		Body:       "<b>This email is from the queue</b>",
		Attachment: []string{"https://i.imgur.com/UbUQWHO.jpeg"},
	}
	q.Send(mailTemp)
	// Consume service
	q.Consume(smMail)
}
