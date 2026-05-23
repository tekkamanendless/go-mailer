package mailer_test

import (
	"context"
	"net/mail"
	"testing"

	"github.com/tekkamanendless/go-mailer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMailer(t *testing.T) {
	ctx := mailer.WithDummyMode(context.Background())

	apiKey := "api-key-1"
	client := mailer.New(mailer.TypeSendgrid, mailer.WithAPIKey(apiKey))
	require.NotNil(t, client)

	assert.Equal(t, 0, len(mailer.Dummy().Outbox(apiKey, "sender@example.com")))
	assert.Nil(t, mailer.Dummy().LastMessageInOutbox(apiKey, "sender@example.com"))
	assert.Equal(t, 0, len(mailer.Dummy().Inbox(apiKey, "receiver@example.com")))
	assert.Nil(t, mailer.Dummy().LastMessageInInbox(apiKey, "receiver@example.com"))

	err := client.SendMail(ctx, mailer.Message{From: mail.Address{Address: "sender@example.com"}, To: mail.Address{Address: "receiver@example.com"}, Subject: "Test1"})
	assert.Nil(t, err)

	assert.Equal(t, 1, len(mailer.Dummy().Outbox(apiKey, "sender@example.com")))
	if message := mailer.Dummy().LastMessageInOutbox(apiKey, "sender@example.com"); assert.NotNil(t, message) {
		assert.Equal(t, "Test1", message.Subject)
	}
	assert.Equal(t, 1, len(mailer.Dummy().Inbox(apiKey, "receiver@example.com")))
	if message := mailer.Dummy().LastMessageInInbox(apiKey, "receiver@example.com"); assert.NotNil(t, message) {
		assert.Equal(t, "Test1", message.Subject)
	}
}
