package mail

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	safeHTML "filespace/pkg/util/html"

	mail "github.com/wneessen/go-mail"
)

var (
	client     *mail.Client
	clientOnce sync.Once
)

func getClient() (*mail.Client, error) {
	var err error

	clientOnce.Do(func() {
		client, err = mail.NewClient("smtp.sendgrid.net",
			mail.WithPort(2525),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername("apikey"),
			mail.WithPassword(os.Getenv("SERVER_EMAIL_PASSWORD")),
		)
		if err != nil {
			log.Printf("Failed to initialize SMTP client: %v", err)
		}
	})

	return client, err
}

func newMessage(options *Message) (*mail.Msg, error) {
	message := mail.NewMsg()
	if err := message.From(os.Getenv("SERVER_EMAIL")); err != nil {
		return nil, err
	}

	if err := message.To(options.To); err != nil {
		return nil, err
	}

	if options.Subject == "" {
		return nil, errors.New("subject is required, got empty string")
	}

	if options.Body == "" {
		return nil, errors.New("body is required, got empty string")
	}

	if options.ContentType != mail.ContentType(mail.TypeTextPlain) && options.ContentType != mail.ContentType(mail.TypeTextHTML) {
		return nil, fmt.Errorf("invalid content type. expected %s or %s, got %s", mail.TypeTextPlain, mail.TypeTextHTML, options.ContentType)
	}

	message.Subject(options.Subject)
	message.SetBodyString(mail.ContentType(options.ContentType), options.Body)

	return message, nil
}

func SendPlainTextEmail(options *Message) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	if options.ContentType != mail.TypeTextPlain {
		return errors.New("invalid content type. expected TypeTextPlain")
	}

	message, err := newMessage(options)
	if err != nil {
		return err
	}

	if err := client.DialAndSend(message); err != nil {
		fmt.Println("Failed to send email.", err)
		client = nil
	}

	return nil
}

func SendHTMLEmail(options *Message) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	if options.ContentType != mail.TypeTextHTML {
		return errors.New("invalid content type. expected TypeTextHTML")
	}

	body, err := safeHTML.Check(options.Body)
	if err != nil {
		return fmt.Errorf("invalid HTML content: %v", err)
	}

	options.Body = body

	message, err := newMessage(options)
	if err != nil {
		return err
	}

	if err := client.DialAndSend(message); err != nil {
		fmt.Println("Falied to send email.", err)
		client = nil
	}

	return nil
}
