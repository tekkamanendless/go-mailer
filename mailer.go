package mailer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/cenkalti/backoff/v5"
)

// Type is the type of mailer.
type Type string

const (
	TypeMailgun  Type = "mailgun"
	TypeSendgrid Type = "sendgrid"
)

// ErrInvalidType is returned when an invalid mailer type is provided.
var ErrInvalidType = errors.New("invalid mailer type")

// Mailer is the mailer client.
type Mailer struct {
	config     Config
	mailerType Type
}

// Config is the configuration for the mailer.
type Config struct {
	APIKey         string                // The API key for the mailer.
	Domain         string                // The domain for the mailer (only used for Mailgun).
	Retry          bool                  // Whether to retry sending the mail.
	BackoffOptions []backoff.RetryOption // The backoff options for the mailer.  If Retry is true, then these options will be used to retry the send.
}

// New creates a new mailer.
func New(mailerType Type, options ...Option) *Mailer {
	m := &Mailer{
		mailerType: mailerType,
		config:     Config{},
	}

	for _, option := range options {
		option(&m.config)
	}

	return m
}

// SendMail sends a mail.
func (m *Mailer) SendMail(ctx context.Context, message Message) error {
	if m.mailerType == TypeMailgun && m.config.Domain == "" {
		return fmt.Errorf("mailer is not set up; missing domain")
	}
	if m.config.APIKey == "" {
		return fmt.Errorf("mailer is not set up; missing API key")
	}

	if IsDummyMode(ctx) {
		slog.InfoContext(ctx, "Dummy mode is enabled; not actually sending e-mail.")
		return sendMailWithDummy(ctx, m.config.APIKey, message)
	}

	if !m.config.Retry {
		return m.sendMailOnce(ctx, message)
	}

	_, err := backoff.Retry(ctx, func() (bool, error) {
		return true, m.sendMailOnce(ctx, message)
	}, m.config.BackoffOptions...)
	return err
}

// sendMailOnce sends a mail once; it does not attempt any retries.
func (m *Mailer) sendMailOnce(ctx context.Context, message Message) error {
	switch m.mailerType {
	case TypeMailgun:
		return sendMailWithMailgun(ctx, m.config.Domain, m.config.APIKey, message)
	case TypeSendgrid:
		return sendMailWithSendgrid(ctx, m.config.APIKey, message)
	default:
		return ErrInvalidType
	}
}
