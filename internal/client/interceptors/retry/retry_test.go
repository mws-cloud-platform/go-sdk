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
	}{
		{
			Name:     "success",
			Invoker:  okInvoker,
			Response: okResponse{},
		},
		{
			Name:    "invoker error no retry",
			Invoker: errInvoker,
			Error:   errInvoke,
		},
		{
			Name:     "api error no retry",
			Invoker:  okInvoker,
			Response: errResponse{err: errResponses, code: 500},
		},
		{
			Name:     "retry error",
			Options:  []retry.Option{retry.WithRetryer(errRetryer{})},
			Invoker:  okInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Error:    errRetry,
		},
		{
			Name:     "context canceled",
			Ctx:      doneCtx,
			Options:  []retry.Option{retry.WithRetryer(delayedRetryer{})},
			Invoker:  okInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Error:    context.Canceled,
		},
		{
			Name:     "api error retry",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{})},
			Invoker:  okInvoker,
			Response: errResponse{err: errResponses, code: 500},
		},
		{
			Name:     "api error retry fail",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{}), retry.WithFailOnExhaustedAttempts(true)},
			Invoker:  okInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Error:    errors.RetryAttemptsExhaustedError{},
		},
		{
			Name:     "invoker error retry",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{})},
			Invoker:  errInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Error:    errInvoke,
		},
		{
			Name:     "invoker error retry fail",
			Options:  []retry.Option{retry.WithRetryer(okRetryer{}), retry.WithFailOnExhaustedAttempts(true)},
			Invoker:  errInvoker,
			Response: errResponse{err: errResponses, code: 500},
			Error:    errors.RetryAttemptsExhaustedError{},
		},
		{
			Name:     "retry with policy",
			Options:  []retry.Option{retry.WithRetryer(errRetryer{}), retry.WithFailOnExhaustedAttempts(true)},
			Invoker:  okInvoker,
			Response: errResponse{err: &mwserrors.APIError{RetryPolicy: &mwserrors.RetryPolicy{}}, code: 500},
			Error:    errors.RetryAttemptsExhaustedError{},
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			interceptor := retry.New(v.Options...)
			ctx := t.Context()
			if v.Ctx != nil {
				ctx = v.Ctx()
			}

			err := interceptor(ctx, nil, v.Response, v.Invoker)
			if v.Error != nil {
				require.ErrorIs(t, err, v.Error)
			} else {
				require.NoError(t, err)
			}
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
