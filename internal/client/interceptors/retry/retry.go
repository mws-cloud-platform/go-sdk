package retry

import (
	"context"
	"errors"
	"time"

	commonclient "go.mws.cloud/go-sdk/internal/client"
	clienterrors "go.mws.cloud/go-sdk/internal/client/errors"
	mwserrors "go.mws.cloud/go-sdk/mws/errors"
	"go.mws.cloud/go-sdk/pkg/backoff"
)

const (
	defaultAttempts = 3
)

func New(options ...Option) commonclient.Interceptor {
	cfg := defaultConfig()
	for _, option := range options {
		option.apply(cfg)
	}

	i := &retryInvoker{
		retryer:                 cfg.retryer,
		failOnExhaustedAttempts: cfg.failOnExhaustedAttempts,
	}
	return i.interceptor
}

type retryInvoker struct {
	retryer                 Retryer
	failOnExhaustedAttempts bool
}

func (i *retryInvoker) interceptor(ctx context.Context, request any, response commonclient.APIResp, invoker commonclient.Invoker) error {
	ctx = i.checkContext(ctx)

	var (
		err, invokerErr error
		attempt         int
		maxAttempts     = defaultAttempts
		retryer         = i.retryer
	)
	for attempt < maxAttempts {
		ctx = WithAttempt(ctx, attempt)
		if invokerErr = invoker(ctx, request, response); invokerErr != nil {
			err = invokerErr
		} else {
			err = response.GetErr()
		}
		if err == nil {
			return nil
		}

		var apiError *mwserrors.APIError
		if attempt == 0 && errors.As(err, &apiError) && apiError.RetryPolicy != nil {
			retryer = policyRetryer{
				backoff: newBackoff(apiError.RetryPolicy),
				err:     apiError,
			}
			maxAttempts = apiError.RetryPolicy.RetryCount
		}

		delay, retryErr := retryer.RetryDelay(attempt, err)
		if errors.Is(retryErr, clienterrors.ErrNoRetry) {
			return invokerErr
		}
		if retryErr != nil {
			return retryErr
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			attempt++
		}
	}

	if !i.failOnExhaustedAttempts {
		return invokerErr
	}
	return clienterrors.NewRetryAttemptsExhaustedError(maxAttempts, err)
}

func (i *retryInvoker) checkContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return ctx
}

type Retryer interface {
	RetryDelay(attempt int, err error) (time.Duration, error)
}

type NoRetry struct{}

func (NoRetry) RetryDelay(int, error) (time.Duration, error) {
	return 0, clienterrors.ErrNoRetry
}

type policyRetryer struct {
	backoff backoff.Backoff
	err     *mwserrors.APIError
}

func (r policyRetryer) RetryDelay(attempt int, err error) (time.Duration, error) {
	var apiError *mwserrors.APIError
	if !errors.As(err, &apiError) || apiError.Code != r.err.Code || apiError.Status != r.err.Status {
		return 0, clienterrors.ErrNoRetry
	}
	return r.backoff.RetryDelay(attempt), nil
}

func newBackoff(r *mwserrors.RetryPolicy) backoff.Backoff {
	b, _ := backoff.New(r.RetryTimeoutScale, r.RetryTimeout, r.MaxRetryTimeout)
	return b
}
