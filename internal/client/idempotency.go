package client

import (
	"context"

	"github.com/google/uuid"
)

type idempotencyKeyCtxKey struct{}

// IdempotencyKeyFromContext retrieves idempotency key from the context.
func IdempotencyKeyFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v, ok := ctx.Value(idempotencyKeyCtxKey{}).(string)
	if !ok {
		return ""
	}
	return v
}

// WithIdempotencyKey adds given idempotency key to the context.
func WithIdempotencyKey(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, idempotencyKeyCtxKey{}, id)
}

// IdempotencyKeyInjector is a client interceptor that generates and injects an
// idempotency key into request context.
func IdempotencyKeyInjector(
	ctx context.Context,
	request any,
	response APIResp,
	invoker Invoker,
) error {
	if IdempotencyKeyFromContext(ctx) == "" {
		idempotencyKey := uuid.NewString()
		ctx = WithIdempotencyKey(ctx, idempotencyKey)
	}
	return invoker(ctx, request, response)
}
