package mailer_test

import (
	"fmt"
	"net/mail"
	"testing"

	"github.com/tekkamanendless/go-mailer"

	"github.com/stretchr/testify/assert"
)

func TestDummyMailServer(t *testing.T) {
	t.Run("First function call", func(t *testing.T) {
		t.Run("Send", func(t *testing.T) {
			mailServer := mailer.DummyMailServer{}
			err := mailServer.Send("api-key-1", mailer.Message{From: mail.Address{Address: "sender@example.com"}, To: mail.Address{Address: "receiver@example.com"}, Subject: "Test1"})
			assert.Nil(t, err)
		})
		t.Run("Outbox", func(t *testing.T) {
			mailServer := mailer.DummyMailServer{}
			assert.Nil(t, mailServer.LastMessageInOutbox("api-key-1", "sender@example.com"))
		})
		t.Run("Inbox", func(t *testing.T) {
			mailServer := mailer.DummyMailServer{}
			assert.Nil(t, mailServer.LastMessageInInbox("api-key-1", "receiver@example.com"))
		})
	})
	t.Run("Multiple API keys and multiple messages", func(t *testing.T) {
		apiKey1 := "api-key-1"
		apiKey2 := "api-key-2"

		mailServer := mailer.DummyMailServer{}

		assert.Equal(t, 0, len(mailServer.Outbox(apiKey1, "sender@example.com")))
		assert.Nil(t, mailServer.LastMessageInOutbox(apiKey1, "sender@example.com"))
		assert.Equal(t, 0, len(mailServer.Inbox(apiKey1, "receiver@example.com")))
		assert.Nil(t, mailServer.LastMessageInInbox(apiKey1, "receiver@example.com"))

		assert.Equal(t, 0, len(mailServer.Outbox(apiKey2, "sender@example.com")))
		assert.Nil(t, mailServer.LastMessageInOutbox(apiKey2, "sender@example.com"))
		assert.Equal(t, 0, len(mailServer.Inbox(apiKey2, "receiver@example.com")))
		assert.Nil(t, mailServer.LastMessageInInbox(apiKey2, "receiver@example.com"))

		err := mailServer.Send("", mailer.Message{From: mail.Address{Address: "sender@example.com"}, To: mail.Address{Address: "receiver@example.com"}, Subject: "Test1"})
		assert.NotNil(t, err)

		err = mailServer.Send(apiKey1, mailer.Message{From: mail.Address{Address: "sender@example.com"}, To: mail.Address{Address: "receiver@example.com"}, Subject: "Test1"})
		assert.Nil(t, err)

		assert.Equal(t, 1, len(mailServer.Outbox(apiKey1, "sender@example.com")))
		if message := mailServer.LastMessageInOutbox(apiKey1, "sender@example.com"); assert.NotNil(t, message) {
			assert.Equal(t, "Test1", message.Subject)
		}
		assert.Equal(t, 1, len(mailServer.Inbox(apiKey1, "receiver@example.com")))
		if message := mailServer.LastMessageInInbox(apiKey1, "receiver@example.com"); assert.NotNil(t, message) {
			assert.Equal(t, "Test1", message.Subject)
		}

		assert.Equal(t, 0, len(mailServer.Outbox(apiKey2, "sender@example.com")))
		assert.Nil(t, mailServer.LastMessageInOutbox(apiKey2, "sender@example.com"))
		assert.Equal(t, 0, len(mailServer.Inbox(apiKey2, "receiver@example.com")))
		assert.Nil(t, mailServer.LastMessageInInbox(apiKey2, "receiver@example.com"))

		err = mailServer.Send(apiKey1, mailer.Message{From: mail.Address{Address: "sender@example.com"}, To: mail.Address{Address: "other-receiver@example.com"}, Subject: "Test2"})
		assert.Nil(t, err)

		assert.Equal(t, 2, len(mailServer.Outbox(apiKey1, "sender@example.com")))
		if message := mailServer.LastMessageInOutbox(apiKey1, "sender@example.com"); assert.NotNil(t, message) {
			assert.Equal(t, "Test2", message.Subject)
		}
		assert.Equal(t, 1, len(mailServer.Inbox(apiKey1, "receiver@example.com")))
		assert.NotNil(t, mailServer.LastMessageInInbox(apiKey1, "receiver@example.com"))
		if message := mailServer.LastMessageInInbox(apiKey1, "receiver@example.com"); assert.NotNil(t, message) {
			assert.Equal(t, "Test1", message.Subject)
		}

		assert.Equal(t, 0, len(mailServer.Outbox(apiKey2, "sender@example.com")))
		assert.Nil(t, mailServer.LastMessageInOutbox(apiKey2, "sender@example.com"))
		assert.Equal(t, 0, len(mailServer.Inbox(apiKey2, "receiver@example.com")))
		assert.Nil(t, mailServer.LastMessageInInbox(apiKey2, "receiver@example.com"))
	})
}

func TestDummy(t *testing.T) {
	t.Run("Multiple calls return the same instance", func(t *testing.T) {
		dummy1 := mailer.Dummy()
		dummy2 := mailer.Dummy()
		assert.Equal(t, dummy1, dummy2)
		assert.Equal(t, fmt.Sprintf("%p", dummy1), fmt.Sprintf("%p", dummy2))
	})
}
