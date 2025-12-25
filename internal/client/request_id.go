package client

import (
	"context"

	"github.com/google/uuid"
)

type requestIDCtxKey struct{}

// RequestIDFromContext retrieves request ID from the context.
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v, ok := ctx.Value(requestIDCtxKey{}).(string)
	if !ok {
		return ""
	}
	return v
}

// WithRequestID adds given request ID to the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey{}, id)
}

// RequestIDInjector is a client interceptor that generates and injects a
// request ID into request context.
func RequestIDInjector(
	ctx context.Context,
	request any,
	response APIResp,
	invoker Invoker,
) error {
	requestID := uuid.NewString()
	ctx = WithRequestID(ctx, requestID)
	return invoker(ctx, request, response)
}
