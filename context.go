package mailer

import "context"

// contextValue is the type for the context value.
type contextValue string

const (
	// contextValueDummyMode is the context value for dummy mode.
	contextValueDummyMode contextValue = "dummy-mode"
)

// WithDummyMode returns a new context with dummy mode enabled.
func WithDummyMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextValueDummyMode, true)
}

// IsDummyMode returns true if dummy mode is enabled.
func IsDummyMode(ctx context.Context) bool {
	return ctx.Value(contextValueDummyMode) == true
}
