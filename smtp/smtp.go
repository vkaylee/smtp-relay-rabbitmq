package smtp

import (
	"errors"
	"github.com/vleedev/smtp-relay-rabbitmq/utils"
	"gopkg.in/gomail.v2"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)
/*
Mail is a generic struct type for representing a mail send request.
*/
type Mail struct {
	message			*gomail.Message
	client			*gomail.Dialer
	emailRegex		*regexp.Regexp
	defaultEmail	string
}

type MailTemplate struct {
	From       	string		`json:"from"`
	To         	[]string	`json:"to"`
	Subject    	string		`json:"subject"`
	BodyType   	string		`json:"body_type"`
	Body       	string		`json:"body"`
	Attachment 	[]string	`json:"attachment"`
}
type Config struct {
	Hostname		string
	Port			int
	Username		string
	Password		string
	DefaultEmail	string
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
		emailRegex: regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"),
		defaultEmail: c.DefaultEmail,
	}
}

func (m *Mail) NewMail(mailTemp MailTemplate) error {
	if strings.EqualFold(mailTemp.From, "") {
		mailTemp.From = m.defaultEmail
	}
	if !m.isEmailValid(mailTemp.From) {
		return errors.New("the form email is not valid")
	}
	m.message.SetHeader("From", mailTemp.From)

	// Set default destination email if no input
	if len(mailTemp.To) == 0 {
		mailTemp.To = []string{m.defaultEmail}
	}

	for _, toEmail := range mailTemp.To {
		if !m.isEmailValid(toEmail) {
			return errors.New("the destination email is not valid")
		}
	}
	m.message.SetHeader("To", mailTemp.To...)
	if strings.EqualFold(mailTemp.Subject, "") {
		return errors.New("the email must have a subject")
	}
	m.message.SetHeader("Subject", mailTemp.Subject)
	if strings.EqualFold(mailTemp.BodyType, "") {
		return errors.New("the email must have a body type")
	}
	if strings.EqualFold(mailTemp.Body, "") {
		return errors.New("the email must have a body contents")
	}
	m.message.SetBody(mailTemp.BodyType, mailTemp.Body)
	if len(mailTemp.Attachment) > 0 {
		for _, attachmentUrl := range mailTemp.Attachment {
			if fileName, err := utils.DownloadFile(attachmentUrl); err != nil {
				utils.ErrPrintln(err)
			} else {
				m.message.Attach(fileName)
				// Auto delete file after 1 minute
				go func() {
					time.Sleep(time.Minute)
					err := os.Remove(fileName)  // remove a single file
					if err != nil {
						utils.ErrPrintln(err)
					}
				}()
			}
		}
	}
	return nil
}
// isEmailValid checks if the email provided passes the required structure
// and length test. It also checks the domain has a valid MX record.
func (m *Mail) isEmailValid(email string) bool {
	if strings.EqualFold(email, "") {
		return false
	}
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	if !m.emailRegex.MatchString(email) {
		return false
	}
	parts := strings.Split(email, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}
/*
MailSend sends email with settings configured by envs.
*/
func (m *Mail) Send() error {
	return m.client.DialAndSend(m.message)
}