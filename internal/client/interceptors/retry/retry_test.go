package retry_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"

	commonclient "go.mws.cloud/go-sdk/internal/client"
	"go.mws.cloud/go-sdk/internal/client/errors"
	"go.mws.cloud/go-sdk/internal/client/interceptors/retry"
	mwserrors "go.mws.cloud/go-sdk/mws/errors"
)

func TestRetry(t *testing.T) {
	for _, v := range []struct {
		Name     string
		Ctx      func() context.Context
		Options  []retry.Option
		Response commonclient.APIResp
		Invoker  commonclient.Invoker
		Error    error
		Attempts int
	}{
		{
			Name:     "success",
			Invoker:  okInvoker,
			Response: okResponse{},
			Attempts: 1,
		},
		{
			Name:     "invoker error no retry",
			Invoker:  errInvoker,
			Error:    errInvoke,
			Attempts: 1,
		},
		{
			Name:     "api error no retry",
			Invoker:  okInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Attempts: 1,
		},
		{
			Name:     "retry error",
			Options:  []retry.Option{retry.WithRetryer(errRetryer{})},
			Invoker:  okInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Error:    errRetry,
			Attempts: 1,
		},
		{
			Name:     "context canceled",
			Ctx:      doneCtx,
			Options:  []retry.Option{retry.WithRetryer(delayedRetryer{})},
			Invoker:  okInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Error:    context.Canceled,
			Attempts: 1,
		},
		{
			Name:     "api error retry",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{})},
			Invoker:  okInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Attempts: 3,
		},
		{
			Name:     "api error retry fail",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{}), retry.WithFailOnExhaustedAttempts(true)},
			Invoker:  okInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Error:    errors.RetryAttemptsExhaustedError{},
			Attempts: 3,
		},
		{
			Name:     "invoker error retry",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{})},
			Invoker:  errInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Error:    errInvoke,
			Attempts: 3,
		},
		{
			Name:     "invoker error retry fail",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{}), retry.WithFailOnExhaustedAttempts(true)},
			Invoker:  errInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Error:    errors.RetryAttemptsExhaustedError{},
			Attempts: 3,
		},
		{
			Name:     "retry with policy",
			Options:  []retry.Option{retry.WithRetryer(errRetryer{}), retry.WithFailOnExhaustedAttempts(true)},
			Invoker:  okInvoker,
			Response: errResponse{err: &mwserrors.APIError{RetryPolicy: &mwserrors.RetryPolicy{}}, code: 500},
			Error:    errors.RetryAttemptsExhaustedError{},
			Attempts: 1,
		},
		{
			Name:     "custom max attempts",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{}), retry.WithMaxAttempts(5)},
			Invoker:  errInvoker,
			Error:    errInvoke,
			Attempts: 5,
		},
		{
			Name:     "max attempts of one disables retries",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{}), retry.WithMaxAttempts(1)},
			Invoker:  errInvoker,
			Error:    errInvoke,
			Attempts: 1,
		},
		{
			Name:     "non-positive max attempts ignored",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{}), retry.WithMaxAttempts(-1)},
			Invoker:  errInvoker,
			Error:    errInvoke,
			Attempts: 3,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			calls := 0
			invoker := func(ctx context.Context, req any, resp commonclient.APIResp) error {
				calls++
				return v.Invoker(ctx, req, resp)
			}

			interceptor := retry.New(v.Options...)
			ctx := t.Context()
			if v.Ctx != nil {
				ctx = v.Ctx()
			}

			err := interceptor(ctx, nil, v.Response, invoker)
			if v.Error != nil {
				require.ErrorIs(t, err, v.Error)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, v.Attempts, calls)
		})
	}
}

const (
	errInvoke    = consterr.Error("invoker")
	errResponses = consterr.Error("response")
	errRetry     = consterr.Error("retryer")
)

type okRetryer struct{}

func (okRetryer) RetryDelay(int, error) (time.Duration, error) {
	return 0, nil
}

type errRetryer struct{}

func (errRetryer) RetryDelay(int, error) (time.Duration, error) {
	return 0, errRetry
}

type delayedRetryer struct{}

func (delayedRetryer) RetryDelay(int, error) (time.Duration, error) {
	return time.Minute, nil
}

func okInvoker(context.Context, any, commonclient.APIResp) error {
	return nil
}

func errInvoker(context.Context, any, commonclient.APIResp) error {
	return errInvoke
}

type okResponse struct{}

func (okResponse) GetErr() error {
	return nil
}

func (okResponse) GetCode() int {
	return 200
}

type errResponse struct {
	err  error
	code int
}

func (r errResponse) GetErr() error {
	return r.err
}

func (r errResponse) GetCode() int {
	return r.code
}

func doneCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
