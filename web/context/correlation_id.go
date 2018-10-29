package web

import "context"

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type correlationIDContextKey int

const (
	// CorrelationIDContextKey is the context key name of the correlation id
	contextKey correlationIDContextKey = iota
)

// CorrelationID returnc correlation-id from request context
func CorrelationID(ctx context.Context) string {
	return ctx.Value(contextKey).(string)
}

// WithCorrelationID returns a new Context carrying correlationID
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, contextKey, correlationID)
}
