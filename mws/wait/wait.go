// Package wait provides utilities for waiting for resource state changes.
package wait

import (
	"context"
	"fmt"
	"time"

	commonerrors "go.mws.cloud/go-sdk/mws/errors"
	"go.mws.cloud/go-sdk/pkg/clock"
)

const (
	// DefaultTimeout is the default timeout for wait operations.
	DefaultTimeout = 5 * time.Minute
	// DefaultRetryInterval is the default interval between waiter retries.
	DefaultRetryInterval = time.Second
	// DefaultNotFoundAllowed is the default number of not found errors allowed
	// before giving up.
	DefaultNotFoundAllowed = 3
)

// Callback is a function that is called to check the current resource state.
type Callback[T any] func(context.Context) (res T, stop bool, err error)

// Waiter is a resource state waiter. It periodically calls the provided
// callback until it returns a result with a stop condition or an error.
//
// If the callback returns an error, waiter returns that error immediately. As
// an exception, the first few errors with status [commonerrors.NotFound] are
// allowed because sometimes the operation state may not be on all replicas yet.
//
// If the callback returns a stop condition, waiter returns callback result
// immediately.
//
// If timeout is exceeded before the callback returns result with stop
// condition, waiter returns an error.
type Waiter[T any] struct {
	callback      Callback[T]
	timeout       time.Duration
	retryInterval time.Duration
	clock         clock.Clock
}

// NewWaiter creates a new waiter.
func NewWaiter[T any](callback Callback[T], opts ...WaiterOption) *Waiter[T] {
	config := waiterConfig{
		timeout:         DefaultTimeout,
		retryInterval:   DefaultRetryInterval,
		notFoundAllowed: DefaultNotFoundAllowed,
	}

	for _, option := range opts {
		option(&config)
	}

	if config.notFoundAllowed > 0 {
		callback = skipNotFoundCallbackWrapper(callback, config.notFoundAllowed)
	}

	if config.clock == nil {
		config.clock = clock.NewReal()
	}

	return &Waiter[T]{
		callback:      callback,
		timeout:       config.timeout,
		retryInterval: config.retryInterval,
		clock:         config.clock,
	}
}

// Wait waits for resource to reach a certain state.
func (w *Waiter[T]) Wait(ctx context.Context) (T, error) {
	ctx, cancel := w.clock.WithTimeout(ctx, w.timeout)
	defer cancel()

	ticker := w.clock.NewTicker(w.retryInterval)

	for {
		select {
		case <-ticker.Chan():
			response, stop, err := w.callback(ctx)
			if err != nil {
				return response, err
			}
			if stop {
				return response, nil
			}
		case <-ctx.Done():
			var response T
			return response, fmt.Errorf("operation canceled: %w", ctx.Err())
		}
	}
}

func skipNotFoundCallbackWrapper[T any](callback Callback[T], notFoundAllowed int) Callback[T] {
	return func(ctx context.Context) (T, bool, error) {
		response, stop, err := callback(ctx)
		if !stop && commonerrors.IsAPIErrorNotFoundStatus(err) && notFoundAllowed > 0 {
			notFoundAllowed--
			return response, false, nil
		}

		return response, stop, err
	}
}
