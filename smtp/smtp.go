package smtp

import (
	"gopkg.in/gomail.v2"
	"log"
)
/*
Mail is a generic struct type for representing a mail send request.
*/
type Mail struct {
	message		*gomail.Message
	client		*gomail.Dialer
}

type MailTemplate struct {
	From       	string
	To         	[]string
	Subject    	string
	BodyType   	string
	Body       	string
	Attachment 	[]string
}
type Config struct {
	Hostname	string
	Port		int
	Username	string
	Password	string
}
func Init(c *Config) *Mail {
	return &Mail{
		message: gomail.NewMessage(),
		client: gomail.NewDialer(
			c.Hostname,
			c.Port,
			c.Username,
			c.Password,
		),
	}
}

func (m *Mail) NewMail(mailTemp MailTemplate) {
	m.message.SetHeader("From", mailTemp.From)
	m.message.SetHeader("To", mailTemp.To...)
	m.message.SetHeader("Subject", mailTemp.Subject)
	m.message.SetBody(mailTemp.BodyType, mailTemp.Body)
	if len(mailTemp.Attachment) > 0 {
		for _, attachment := range mailTemp.Attachment {
			m.message.Attach(attachment)
		}
	}
}
/*
MailSend sends email with settings configured by envs.
*/
func (m *Mail) Send() error {
	return m.client.DialAndSend(m.message)
}
func (m *Mail) Test() {
	m.NewMail(MailTemplate{
		From:       "test@vlee.dev",
		To:         []string{"tuanlm1989@gmail.com"},
		Subject:    "Welcome to smtp-relay-rabbitmq",
		BodyType:   "text/html",
		Body:       "<html><body><p>This one is a test email from smtp-relay-rabbitmq</p></body></html>",
		Attachment: nil,
	})
	if err := m.Send(); err != nil {
		log.Fatalln("Please check your smtp configuration!")
	}
}
