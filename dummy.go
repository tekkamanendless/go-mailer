package mailer

import (
	"context"
)

// sendMailWithDummy sends a dummy mail.
func sendMailWithDummy(_ context.Context, apiKey string, message Message) error {
	return dummyMailServer.Send(apiKey, message)
}
