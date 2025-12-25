package retry

import (
	"context"
)

type attemptCtxKey struct{}

func AttemptFromContext(ctx context.Context) int {
	if ctx == nil {
		return 0
	}
	v, ok := ctx.Value(attemptCtxKey{}).(int)
	if !ok {
		return 0
	}
	return v
}

func WithAttempt(ctx context.Context, attempt int) context.Context {
	return context.WithValue(ctx, attemptCtxKey{}, attempt)
}
