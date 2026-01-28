package log

import (
	"context"
	"errors"

	"go.mws.cloud/util-toolset/pkg/utils/zaputil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.mws.cloud/go-sdk/internal/client"
	clientattribute "go.mws.cloud/go-sdk/internal/client/interceptors/attribute"
	zapattributeadapter "go.mws.cloud/go-sdk/internal/client/interceptors/attribute/zap"
	"go.mws.cloud/go-sdk/internal/client/interceptors/retry"
	mwserrors "go.mws.cloud/go-sdk/mws/errors"
	"go.mws.cloud/go-sdk/pkg/clock"
)

// New creates a new log interceptor.
func New(log *zap.Logger, opts ...Option) client.Interceptor {
	o := options{
		clock: clock.NewReal(),
		stack: func(key string) zap.Field {
			return zap.StackSkip(key, 1)
		},
	}
	for _, opt := range opts {
		opt(&o)
	}

	return func(
		ctx context.Context,
		request any,
		response client.APIResp,
		invoker client.Invoker,
	) (err error) {
		log := log
		start := o.clock.Now()
		fields := fieldsFromContext(ctx)

		defer func() {
			fields = append(fields, zap.Duration("client_req_duration", o.clock.Since(start)))
			if o.trace {
				fields = append(fields, zap.Any("client_req_data", request))
			}
			log = log.With(fields...)

			if r := recover(); r != nil {
				log.Error("panic occurred during request",
					o.stack("stacktrace"),
					zap.Any("reason", r))
				panic(r)
			}
			if err != nil {
				log.Error("request failed", zap.Error(err))
				return
			}

			log = log.With(zap.Int("client_resp_code", response.GetCode()))
			if o.trace {
				log = log.With(zap.Any("client_resp_data", response))
			}

			if respErr := response.GetErr(); respErr != nil {
				logResponseError(log, respErr)
				return
			}
			log.Info("request succeeded")
		}()

		return invoker(ctx, request, response)
	}
}

func logResponseError(log *zap.Logger, respErr error) {
	log = log.With(zap.Error(respErr))

	var target *mwserrors.APIError
	if errors.As(respErr, &target) {
		log = log.With(zap.String("client_resp_api_error_status", target.Status.String()))
		if details := target.Details; details != nil {
			log = log.With(zap.Any("client_resp_api_error_details", details))
		}

		var lvl zapcore.Level
		switch target.Status {
		case mwserrors.NotFound:
			lvl = zap.InfoLevel
		case mwserrors.AlreadyExists:
			lvl = zap.WarnLevel
		default:
			lvl = zap.ErrorLevel
		}

		log.Log(lvl, "got api error response")
		return
	}

	log.Error("got error response")
}

// Option is a log interceptor option.
type Option func(*options)

// WithTrace enables request and response data tracing.
func WithTrace() Option {
	return func(o *options) {
		o.trace = true
	}
}

func withClock(clock clock.Clock) Option {
	return func(o *options) {
		o.clock = clock
	}
}

func withStack(stack func(string) zap.Field) Option {
	return func(o *options) {
		o.stack = stack
	}
}

type options struct {
	trace bool
	clock clock.Clock
	stack func(string) zap.Field
}

func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := zapattributeadapter.Convert(clientattribute.FromContext(ctx))
	fields = append(fields, zap.Int("client_req_attempt", retry.AttemptFromContext(ctx)))
	fields = append(fields, zaputil.Ctx(ctx))
	return fields
}
