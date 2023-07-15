package log

import (
	"context"

	"go.uber.org/zap"
)

type loggerKeyType int

const keyLogger loggerKeyType = iota

func NewCtx(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, keyLogger, l)
}

func FromCtx(ctx context.Context) *Logger {
	l, ok := ctx.Value(keyLogger).(*Logger)
	if ok {
		return l
	}
	return zap.L()
}
