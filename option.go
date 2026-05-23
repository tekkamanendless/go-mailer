package mailer

import "github.com/cenkalti/backoff/v5"

// Option can be used to configure the mailer.
type Option func(*Config)

// WithAPIKey sets the API key for the mailer.
func WithAPIKey(apiKey string) Option {
	return func(c *Config) {
		c.APIKey = apiKey
	}
}

// WithConfig sets the configuration for the mailer.
//
// This fully replaces the existing configuration.
func WithConfig(config Config) Option {
	return func(c *Config) {
		*c = config
	}
}

// WithDomain sets the domain for the mailer.
func WithDomain(domain string) Option {
	return func(c *Config) {
		c.Domain = domain
	}
}

// WithRetry sets the retry configuration for the mailer.
func WithRetry(retry bool, options ...backoff.RetryOption) Option {
	return func(c *Config) {
		c.Retry = retry
		c.BackoffOptions = options
	}
}
