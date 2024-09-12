package logger

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
)

type contextKey string

const logContextKey contextKey = "logContext"

type LogContext struct {
	mu     sync.RWMutex
	values map[string]interface{}
}

func NewLogContext() *LogContext {
	return &LogContext{
		values: make(map[string]interface{}),
	}
}

func NewCtxFromLogContext(ctx context.Context, lc *LogContext) context.Context {
	return context.WithValue(ctx, logContextKey, lc)
}

func (lc *LogContext) Set(key string, value interface{}) *LogContext {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.values[key] = value
	return lc
}

func (lc *LogContext) Get(key string) (interface{}, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	value, ok := lc.values[key]
	return value, ok
}

func (lc *LogContext) Clone() *LogContext {

	lc.mu.RLock()

	defer lc.mu.RUnlock()

	clone := NewLogContext()

	for key, value := range lc.values {
		clone.values[key] = value
	}
	return clone
}

func WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, NewLogContext())
}

func GetLogContext(ctx context.Context) *LogContext {
	lc, ok := ctx.Value(logContextKey).(*LogContext)
	if !ok {
		return nil
	}
	return lc
}

func logWithContext(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	lc := GetLogContext(ctx)

	if lc == nil {
		return event
	}

	lc.mu.RLock()

	defer lc.mu.RUnlock()

	for key, value := range lc.values {
		event = event.Interface(key, value)
	}

	return event
}

func TraceWithCtx(ctx context.Context) *zerolog.Event {
	return logWithContext(ctx, log.Trace())
}

func DebugWithCtx(ctx context.Context) *zerolog.Event {
	return logWithContext(ctx, log.Debug())
}

func InfoWithCtx(ctx context.Context) *zerolog.Event {
	return logWithContext(ctx, log.Info())
}

func WarnWithCtx(ctx context.Context) *zerolog.Event {
	return logWithContext(ctx, log.Warn())
}

func ErrorWithCtx(ctx context.Context) *zerolog.Event {
	return logWithContext(ctx, log.Error())
}

func FatalWithCtx(ctx context.Context) *zerolog.Event {
	return logWithContext(ctx, log.Fatal())
}
