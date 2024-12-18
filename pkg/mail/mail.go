package mail

import (
	"log"
	"os"

	mail "github.com/wneessen/go-mail"
)

func newClient() (*mail.Client, error) {
	client, err := mail.NewClient("smtp.zoho.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername((os.Getenv("SERVER_EMAIL"))),
		mail.WithPassword(os.Getenv("SERVER_EMAIL_PASSWORD")),
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func newMessage(to string) (*mail.Msg, error) {
	message := mail.NewMsg()
	if err := message.From(os.Getenv("SERVER_EMAIL")); err != nil {
		return nil, err
	}

	if err := message.To(to); err != nil {
		return nil, err
	}

	return message, nil
}

func SendPlainTextEmail(to string, subject string, body string) {
	client, err := newClient()

	if err != nil {
		log.Fatal(err)
	}

	message, _ := newMessage(to)
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	if err := client.DialAndSend(message); err != nil {
		log.Fatal(err)
	}
}

func SendHTMLEmail(to string, subject string, body string) {
	client, err := newClient()

	if err != nil {
		log.Fatal(err)
	}

	message, _ := newMessage(to)
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, body)

	if err := client.DialAndSend(message); err != nil {
		log.Fatal(err)
	}
}
