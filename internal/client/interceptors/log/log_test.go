package log

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.mws.cloud/go-sdk/internal/client"
	clientattribute "go.mws.cloud/go-sdk/internal/client/interceptors/attribute"
	mwsinternalerrors "go.mws.cloud/go-sdk/internal/errors"
	"go.mws.cloud/go-sdk/pkg/clock/fakeclock"
	commonmodel "go.mws.cloud/go-sdk/service/common/model"
)

func TestLog(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/log/golden"),
		golden.WithRecreateOnUpdate())

	for _, tc := range []struct {
		name     string
		opts     []Option
		response client.APIResp
		invoker  client.Invoker
	}{
		{
			name:     "ok",
			response: &helloResponse{},
			invoker: func(_ context.Context, _ any, resp client.APIResp) error {
				respPtr := resp.(*helloResponse)
				*respPtr = helloResponse{
					Code:        200,
					Response200: &helloResponse200{},
				}
				return nil
			},
		},
		{
			name:     "ok_with_trace",
			opts:     []Option{WithTrace()},
			response: &helloResponse{},
			invoker: func(_ context.Context, _ any, resp client.APIResp) error {
				respPtr := resp.(*helloResponse)
				*respPtr = helloResponse{
					Code:        200,
					Response200: &helloResponse200{},
				}
				return nil
			},
		},
		{
			name: "transport_error",
			invoker: func(context.Context, any, client.APIResp) error {
				return errors.New("error")
			},
		},
		{
			name:     "with_api_error_not_found",
			response: &helloResponse{},
			invoker: func(_ context.Context, _ any, resp client.APIResp) error {
				respPtr := resp.(*helloResponse)
				*respPtr = helloResponse{
					Code: 404,
					Response404: &commonmodel.ApiError{
						Code: commonmodel.ApiErrorCode_NOT_FOUND,
					},
				}
				return nil
			},
		},
		{
			name:     "with_api_error_already_exists",
			response: &helloResponse{},
			invoker: func(_ context.Context, _ any, resp client.APIResp) error {
				respPtr := resp.(*helloResponse)
				*respPtr = helloResponse{
					Code: 409,
					Response409: &commonmodel.ApiError{
						Code: commonmodel.ApiErrorCode_ALREADY_EXISTS,
					},
				}
				return nil
			},
		},
		{
			name:     "with_api_error_internal",
			response: &helloResponse{},
			invoker: func(_ context.Context, _ any, resp client.APIResp) error {
				respPtr := resp.(*helloResponse)
				*respPtr = helloResponse{
					Code: 500,
					Response500: &commonmodel.ApiError{
						Code: commonmodel.ApiErrorCode_INTERNAL,
						Details: map[string]json.RawMessage{
							"foo": []byte(`"bar"`),
						},
					},
				}
				return nil
			},
		},
		{
			name:     "with_api_error_status_warn",
			response: &helloResponse{},
			invoker: func(_ context.Context, _ any, resp client.APIResp) error {
				respPtr := resp.(*helloResponse)
				*respPtr = helloResponse{
					Code: 111,
					ResponseBaseError: &commonmodel.BaseError{
						Code: ptr.Get("invalid_code"),
					},
				}
				return nil
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actual, logger := newTestLogger()
			clock := fakeclock.NewFake(fakeclock.WithAdvancingNowBy(512 * time.Millisecond))
			opts := []Option{withClock(clock)}
			opts = append(opts, tc.opts...)
			f := New(logger, opts...)

			ctx := clientattribute.WithContext(t.Context(), testAttrs()...)
			_ = f(ctx,
				request{Text: "trace"},
				tc.response,
				tc.invoker,
			)

			dir.String(t, tc.name+".jsonl", actual.String())
		})
	}
}

func TestLog_span_context(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/log_span_context/golden"),
		golden.WithRecreateOnUpdate())

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1},
		SpanID:     trace.SpanID{1},
		TraceFlags: trace.FlagsSampled,
		Remote:     false,
	})

	ctx := trace.ContextWithSpanContext(t.Context(), spanCtx)
	actual, logger := newTestLogger()
	clock := fakeclock.NewFake(fakeclock.WithAdvancingNowBy(256 * time.Millisecond))
	f := New(logger, withClock(clock))

	ctx = clientattribute.WithContext(ctx, testAttrs()...)

	err := f(ctx, request{}, &response{}, ok)
	require.NoError(t, err)

	dir.String(t, "expected.jsonl", actual.String())
}

func TestLog_panic(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/log_panic/golden"),
		golden.WithRecreateOnUpdate())

	actual, logger := newTestLogger()
	stack := func(s string) zap.Field { return zap.String(s, "stub") }
	clock := fakeclock.NewFake(fakeclock.WithAdvancingNowBy(1024 * time.Millisecond))
	f := New(logger, withClock(clock), withStack(stack))

	ctx := clientattribute.WithContext(t.Context(), testAttrs()...)

	require.Panics(t, func() {
		_ = f(ctx,
			request{Text: "trace"},
			nil,
			func(context.Context, any, client.APIResp) error {
				panic("oops!")
			},
		)
	})
	dir.String(t, "expected.jsonl", actual.String())
}

func ok(context.Context, any, client.APIResp) error {
	return nil
}

type helloResponse struct {
	Code              int
	Response200       *helloResponse200
	Response404       *commonmodel.ApiError
	Response409       *commonmodel.ApiError
	Response500       *commonmodel.ApiError
	ResponseBaseError *commonmodel.BaseError
}

func (r *helloResponse) GetCode() int {
	return r.Code
}

func (r *helloResponse) GetErr() error {
	if r.Response404 != nil {
		return mwsinternalerrors.WrapAPIGenError(r.Code, r.Response404)
	}
	if r.Response409 != nil {
		return mwsinternalerrors.WrapAPIGenError(r.Code, r.Response409)
	}
	if r.Response500 != nil {
		return mwsinternalerrors.WrapAPIGenError(r.Code, r.Response500)
	}
	if r.ResponseBaseError != nil {
		return mwsinternalerrors.WrapAPIGenError(r.Code, r.ResponseBaseError)
	}
	return nil
}

type helloResponse200 struct {
}

type request struct {
	Text string `json:"text"`
}

type response struct {
	Text string `json:"text"`
	code int
	err  error
}

func (r response) GetCode() int {
	return r.code
}

func (r response) GetErr() error {
	return r.err
}

func newTestLogger() (out *strings.Builder, logger *zap.Logger) {
	buf := &strings.Builder{}
	enc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        zapcore.OmitKey,
		LevelKey:       "level",
		NameKey:        "component",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	core := zapcore.NewCore(enc, zapcore.AddSync(buf), zapcore.DebugLevel)
	return buf, zap.New(core)
}

func testAttrs() []clientattribute.KeyValue {
	return []clientattribute.KeyValue{
		clientattribute.Bool("bool", true),
		clientattribute.Int64("int", 42),
		clientattribute.Float64("float", 128),
		clientattribute.String("string", "foo"),
		clientattribute.BoolSlice("bools", []bool{true, false, true}),
		clientattribute.Int64Slice("ints", []int64{1, 2, 3}),
		clientattribute.Float64Slice("floats", []float64{1, 2}),
		clientattribute.StringSlice("strings", []string{"a", "b"}),
		{},               // empty invalid
		{Key: "invalid"}, // invalid with key
	}
}
