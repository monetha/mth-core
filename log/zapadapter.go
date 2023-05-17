package log

import (
	"fmt"

	"go.temporal.io/sdk/log"
	"go.uber.org/zap"
)

// ZapAdapter is a log.Logger implementation that wraps a Zap logger.
type ZapAdapter struct {
	zl *zap.Logger
}

// NewZapAdapter creates a new ZapAdapter.
func NewZapAdapter(zapLogger *zap.Logger) *ZapAdapter {
	return &ZapAdapter{
		zl: zapLogger.WithOptions(zap.AddCallerSkip(1)),
	}
}

// ZapLogger returns the underlying Zap logger.
func (log *ZapAdapter) fields(keyvals []interface{}) []zap.Field {
	if len(keyvals)%2 != 0 {
		return []zap.Field{zap.Error(fmt.Errorf("odd number of keyvals pairs: %v", keyvals))}
	}

	var fields []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keyvals[i])
		}
		fields = append(fields, zap.Any(key, keyvals[i+1]))
	}

	return fields
}

// Debug logs a message at debug level.
func (log *ZapAdapter) Debug(msg string, keyvals ...interface{}) {
	log.zl.Debug(msg, log.fields(keyvals)...)
}

// Info logs a message at info level.
func (log *ZapAdapter) Info(msg string, keyvals ...interface{}) {
	log.zl.Info(msg, log.fields(keyvals)...)
}

// Warn logs a message at warn level.
func (log *ZapAdapter) Warn(msg string, keyvals ...interface{}) {
	log.zl.Warn(msg, log.fields(keyvals)...)
}

// Error logs a message at error level.
func (log *ZapAdapter) Error(msg string, keyvals ...interface{}) {
	log.zl.Error(msg, log.fields(keyvals)...)
}

// With returns a child logger with the provided keyvals.
func (log *ZapAdapter) With(keyvals ...interface{}) log.Logger {
	return &ZapAdapter{zl: log.zl.With(log.fields(keyvals)...)}
}
