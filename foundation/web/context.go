package web

import (
	"context"
	"net/http"
)

type ctxKey int

const (
	writerKey ctxKey = iota + 1
)

func setWriter(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, writerKey, w)
}

// GetWriter returns the underlying writer for the request.
func GetWriter(ctx context.Context) http.ResponseWriter {
	v, ok := ctx.Value(writerKey).(http.ResponseWriter)
	if !ok {
		return nil
	}

	return v
}

func Set(ctx context.Context, key, val any) context.Context {
	return context.WithValue(ctx, key, val)
}

func Get[T any](ctx context.Context, key any) T {
	return ctx.Value(key).(T)
}
