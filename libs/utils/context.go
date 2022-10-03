package utils

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

// RemoteAddress is the key to store the remote ip address in the context
type RemoteAddress struct {
	address string
}

// CtxLogger is the key to store the logger object in contexts
type CtxLogger struct {
	logger string
}

var addressID = RemoteAddress{address: "remote-address"}
var contextLogger = CtxLogger{logger: "ctx-logger"}

// WithLogger adds the logger to the current context
func WithLogger(ctx context.Context, logger zap.Logger) context.Context {
	return context.WithValue(ctx, contextLogger, logger)
}

// WithRemoteAddress adds the request remote address to the current context
func WithRemoteAddress(ctx context.Context, remoteAddress string) context.Context {
	return context.WithValue(ctx, addressID, remoteAddress)
}

// CreateRequestContext adds contextual information to the default request
func CreateRequestContext(r http.Request) context.Context {
	ctx := r.Context()
	ctx = WithRemoteAddress(ctx, r.RemoteAddr)

	logger := NewLogger()
	logger = logger.With(
		zap.String("remote", r.RemoteAddr),
	)

	ctx = WithLogger(ctx, *logger)

	return ctx
}
