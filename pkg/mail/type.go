package mail

import (
	mail "github.com/wneessen/go-mail"
)

type Message struct {
	To          string
	Subject     string
	Body        string
	ContentType mail.ContentType
}

const PlainTextEmail mail.ContentType = mail.TypeTextPlain
const HTMLTextEmail mail.ContentType = mail.TypeTextHTML
