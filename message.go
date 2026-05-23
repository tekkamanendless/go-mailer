package mailer

import "net/mail"

// Message is a piece of mail.
type Message struct {
	From          mail.Address
	To            mail.Address
	Subject       string
	BodyPlainText string
	BodyHTML      string
}
