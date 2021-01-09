package queue

import (
	"github.com/streadway/amqp"
	"github.com/vleedev/smtp-relay-rabbitmq/smtp"
	"github.com/vleedev/smtp-relay-rabbitmq/utils"
	"log"
)
type Queue struct {
	Connection	*amqp.Connection
	queueInfo	amqp.Queue
	Channel 	*amqp.Channel
}
func Init(queueName string, amqpUrl string) *Queue {
	conn, err := amqp.Dial(amqpUrl)
	utils.CheckErr(err)

	ch, err := conn.Channel()
	utils.CheckErr(err)

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	utils.CheckErr(err)
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
	utils.CheckErr(err)
	return msgs
}

func (q *Queue) Consume(m *smtp.Mail)  {
	forever := make(chan bool)
	go func() {
		for d := range q.messageChan() {
			log.Println(string(d.Body))
		}
	}()
	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func (q *Queue) Send(info string) {
	err := q.Channel.Publish(
		"",     // exchange
		q.queueInfo.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(info),
		})
	utils.CheckErr(err)
}