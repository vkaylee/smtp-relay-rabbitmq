package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vleedev/smtp-relay-rabbitmq/smtp"
	"github.com/vleedev/smtp-relay-rabbitmq/utils"
	"log"
	"strings"
)
const acceptedContentType = "application/json"

type Queue struct {
	Connection	*amqp.Connection
	queueInfo	amqp.Queue
	Channel 	*amqp.Channel
}
func Init(queueName string, amqpUrl string) *Queue {
	conn, err := amqp.Dial(amqpUrl)
	utils.ErrFatal(err)

	ch, err := conn.Channel()
	utils.ErrFatal(err)

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	utils.ErrFatal(err)
	return &Queue{
		Connection: conn,
		queueInfo:  q,
		Channel: ch,
	}
}

func (q *Queue) messageChan() <-chan amqp.Delivery {
	msgs, err := q.Channel.Consume(
		q.queueInfo.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.ErrFatal(err)
	return msgs
}

func (q *Queue) Consume(m *smtp.Mail)  {
	forever := make(chan bool)
	go func() {
		for d := range q.messageChan() {

			log.Println("Received: 1 email from queue.")

			if strings.EqualFold(acceptedContentType, d.ContentType) {
				var email smtp.MailTemplate
				err := json.Unmarshal(d.Body, &email)
				utils.ErrPrintln(err)


				if err := m.NewMail(email); err != nil {
					utils.ErrPrintln(err)
				} else if err := m.Send(); err != nil {
					utils.ErrPrintln(err)
				} else {
					log.Println("Sent: 1 email from queue.")
				}
			} else {
				utils.ErrPrintln(errors.New(fmt.Sprintf("the queue message content-type must be in %s", acceptedContentType)))
			}
		}
	}()
	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func (q *Queue) Send(info interface{}) {
	data, err := json.Marshal(info)
	utils.ErrFatal(err)
	err = q.Channel.Publish(
		"",     // exchange
		q.queueInfo.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: acceptedContentType,
			Body:        data,
		})
	utils.ErrFatal(err)
}