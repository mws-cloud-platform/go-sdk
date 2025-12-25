package client

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/trace"
)

// Error is a client error.
type Error struct {
	Err       error
	RequestID string
	TraceID   string
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) Is(err error) bool {
	return errors.As(err, &Error{})
}

type ErrorWrapperSetter interface {
	SetErrorWrapper(f func(error) error)
}

// ErrorWrapper is a client interceptor that wraps errors into [Error].
func ErrorWrapper(
	ctx context.Context,
	request any,
	response APIResp,
	invoker Invoker,
) error {
	err := invoker(ctx, request, response)
	if err == nil && response.GetErr() == nil {
		return nil
	}
	if err == nil {
		if setter, ok := response.(ErrorWrapperSetter); ok {
			setter.SetErrorWrapper(func(err error) error {
				return wrapError(ctx, err)
			})
		}
		return nil
	}

	return wrapError(ctx, err)
}

func wrapError(ctx context.Context, err error) error {
	wrapped := Error{Err: err, RequestID: RequestIDFromContext(ctx)}
	if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
		wrapped.TraceID = spanCtx.TraceID().String()
	}
	return wrapped
}
