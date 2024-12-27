package util

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	// correlationIDCtxKey is the key for the correlation id
	correlationIDCtxKey contextKey = "correlation_id"
)

// Info wraps the *zap.Logger instance in order to prefix each log with a unique request correlation id
func Info(logger *zap.Logger, ctx context.Context, msg, key, val string) {
	correlationId, ok := ctx.Value(correlationIDCtxKey).(string) // to prevent panic
	_ = ok

	logger.
		With(zap.String(string(correlationIDCtxKey), correlationId)).
		Info(msg, zap.String(key, val))
}

// Error wraps the *zap.Logger instance in order to prefix each log with a unique request correlation id
func Error(logger *zap.Logger, ctx context.Context, msg string, err error) {
	correlationId, ok := ctx.Value(correlationIDCtxKey).(string)
	_ = ok

	logger.
		With(zap.String(string(correlationIDCtxKey), correlationId)).
		Error(msg, zap.Error(err))
}

// Warn wraps the *zap.Logger instance in order to prefix each log with a unique request correlation id
func Warn(logger *zap.Logger, ctx context.Context, msg, key, val string) {
	correlationId, ok := ctx.Value(correlationIDCtxKey).(string)
	_ = ok

	logger.
		With(zap.String(string(correlationIDCtxKey), correlationId)).
		Warn(msg, zap.String(key, val))
}
