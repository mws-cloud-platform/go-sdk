package client_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"go.mws.cloud/go-sdk/internal/client"
)

func TestErrorWrapper(t *testing.T) {
	ctx := t.Context()
	err := client.ErrorWrapper(ctx, nil, &response{code: http.StatusOK}, noopInvoke)
	require.NoError(t, err)
}

func TestErrorWrapper_error(t *testing.T) {
	traceID := trace.TraceID(uuid.New())
	expected := client.Error{
		Err:       io.EOF,
		RequestID: uuid.NewString(),
		TraceID:   traceID.String(),
	}

	ctx := t.Context()
	ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  [8]byte{1},
	}))
	ctx = client.WithRequestID(ctx, expected.RequestID)

	actual := client.ErrorWrapper(ctx, nil, nil, func(context.Context, any, client.APIResp) error {
		return expected.Err
	})
	require.ErrorIs(t, actual, expected)
}

func TestErrorWrapper_responseError(t *testing.T) {
	traceID := trace.TraceID(uuid.New())
	expected := client.Error{
		Err:       io.EOF,
		RequestID: uuid.NewString(),
		TraceID:   traceID.String(),
	}

	ctx := t.Context()
	ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  [8]byte{1},
	}))
	ctx = client.WithRequestID(ctx, expected.RequestID)

	resp := response{err: io.EOF}

	err := client.ErrorWrapper(ctx, nil, &resp, noopInvoke)
	require.NoError(t, err)

	responseError := resp.GetErr()
	require.ErrorIs(t, responseError, expected)
}

type response struct {
	code         int
	err          error
	errorWrapper func(err error) error
}

func (r *response) GetCode() int {
	return r.code
}

func (r *response) GetErr() error {
	if r.errorWrapper != nil {
		return r.errorWrapper(r.err)
	}
	return r.err
}

func (r *response) SetErrorWrapper(f func(err error) error) {
	r.errorWrapper = f
}
