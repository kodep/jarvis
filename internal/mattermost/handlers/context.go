package handlers

import (
	"context"

	"go.uber.org/zap"
)

type (
	ctxLoggerKeyType string
)

const (
	ctxLoggerKey ctxLoggerKeyType = "logger"
)

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

func GetLogger(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxLoggerKey).(*zap.Logger); ok {
		return logger
	}

	return nil
}
