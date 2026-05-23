package mailer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mailgun/mailgun-go/v5"
)

// sendMailWithMailgun sends mail with Mailgun.
func sendMailWithMailgun(ctx context.Context, domain string, apiKey string, message Message) error {
	mailgunClient := mailgun.NewMailgun(apiKey)

	mailgunMessage := mailgun.NewMessage(domain, message.From.Address, message.Subject, message.BodyPlainText, message.To.Address)
	mailgunMessage.SetHTML(message.BodyHTML)

	var messageID string
	response, err := mailgunClient.Send(ctx, mailgunMessage)
	if err != nil {
		return err
	}
	messageID = response.ID

	slog.InfoContext(ctx, fmt.Sprintf("E-mail message: Message ID: %s", messageID))
	return nil
}
