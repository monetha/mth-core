package log

import (
	"fmt"

	"go.temporal.io/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapAdapterInterface is a log.Logger implementation that wraps a Zap logger.
type ZapAdapterInterface interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	GetZap() *zap.Logger
}

type ZapAdapter struct {
	Zl *zap.Logger
}

// NewZapAdapter creates a new ZapAdapter.
func NewZapAdapter(zapLogger *zap.Logger) *ZapAdapter {
	return &ZapAdapter{
		Zl: zapLogger.WithOptions(zap.AddCallerSkip(1)),
	}
}

// ZapLogger returns the underlying Zap logger.
func (log *ZapAdapter) fields(keyvals []interface{}) []zap.Field {

	// If keyvals are zap fields - expand them to key-value pairs
	expKV := make([]interface{}, 0)
	for _, v := range keyvals {
		if zv, ok := v.(zapcore.Field); ok {
			expKV = append(expKV, zv.Key, zv.Interface)
		} else {
			expKV = append(expKV, v)
		}
	}

	if len(expKV)%2 != 0 {
		return []zap.Field{zap.Error(fmt.Errorf("odd number of keyvals pairs: %v", expKV))}
	}

	var fields []zap.Field
	for i := 0; i < len(expKV); i += 2 {
		key, ok := expKV[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", expKV[i])
		}
		fields = append(fields, zap.Any(key, expKV[i+1]))
	}

	return fields
}

// Debug logs a message at debug level.
func (log *ZapAdapter) Debug(msg string, keyvals ...interface{}) {
	log.Zl.Debug(msg, log.fields(keyvals)...)
}

// Info logs a message at info level.
func (log *ZapAdapter) Info(msg string, keyvals ...interface{}) {
	log.Zl.Info(msg, log.fields(keyvals)...)
}

// Warn logs a message at warn level.
func (log *ZapAdapter) Warn(msg string, keyvals ...interface{}) {
	log.Zl.Warn(msg, log.fields(keyvals)...)
}

// Error logs a message at error level.
func (log *ZapAdapter) Error(msg string, keyvals ...interface{}) {
	log.Zl.Error(msg, log.fields(keyvals)...)
}

// With returns a child logger with the provided keyvals.
func (log *ZapAdapter) With(keyvals ...interface{}) log.Logger {
	return &ZapAdapter{Zl: log.Zl.With(log.fields(keyvals)...)}
}

// GetZap returns a child logger with the provided keyvals.
func (log *ZapAdapter) GetZap() *zap.Logger {
	return log.Zl
}
