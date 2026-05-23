package mailer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sendgrid/sendgrid-go"
	sendgridmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// sendMailWithSendgrid sends mail with Sendgrid.
func sendMailWithSendgrid(ctx context.Context, apiKey string, message Message) error {
	slog.InfoContext(ctx, fmt.Sprintf("Sendgrid from address: %s", message.From.Address))
	client := sendgrid.NewSendClient(apiKey)
	sendgridMessage := sendgridmail.NewSingleEmail(sendgridmail.NewEmail(message.From.Name, message.From.Address), message.Subject, sendgridmail.NewEmail(message.To.Name, message.To.Address), message.BodyPlainText, message.BodyHTML)
	response, err := client.SendWithContext(ctx, sendgridMessage)
	if err != nil {
		return err
	}
	if response.StatusCode != 202 {
		return fmt.Errorf("error sending email: status code %d, body: %q", response.StatusCode, response.Body)
	}
	return nil
}
