package main

import (
	"github.com/vleedev/smtp-relay-rabbitmq/queue"
	"os"
)

func main() {
	q := queue.Init(
		os.Getenv("QUEUE_NAME"),
		os.Getenv("RABBITMQ_URL"),
		)
	defer q.Connection.Close()
	defer q.Channel.Close()
	q.Consume()
}
