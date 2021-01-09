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
	utils.CheckErr(err)
	smtpConfig := smtp.Config{
		Hostname: os.Getenv("SMTP_HOSTNAME"),
		Port:     smtpPort,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
	smMail := smtp.Init(&smtpConfig)
	// Test smtp service
	smMail.Test()
	// Consume service
	q.Consume(smMail)
}
