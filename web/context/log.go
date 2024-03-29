package web

import (
	"context"

	"github.com/monetha/mth-core/log"
)

// NewLogger returns a request-scoped logger. Use this to log info/errors
// Will automatically log request correlation id.
func NewLogger(ctx context.Context) *log.Logger {
	correlationID := CorrelationID(ctx)
	correlationIDField := log.CorrelationID(correlationID)
	return log.With(correlationIDField)
}
