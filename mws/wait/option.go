package wait

import (
	"time"

	"go.mws.cloud/go-sdk/pkg/clock"
)

// WaiterOption represents a function for configuring a waiter.
type WaiterOption func(w *waiterConfig)

// WithTimeout sets the waiter timeout.
func WithTimeout(timeout time.Duration) WaiterOption {
	return func(w *waiterConfig) {
		w.timeout = timeout
	}
}

// WithRetryInterval sets the waiter retry interval.
func WithRetryInterval(retryInterval time.Duration) WaiterOption {
	return func(w *waiterConfig) {
		w.retryInterval = retryInterval
	}
}

// WithNotFoundAllowed sets the number of times the waiter can get the not found
// error before giving up.
func WithNotFoundAllowed(notFoundAllowed int) WaiterOption {
	return func(w *waiterConfig) {
		w.notFoundAllowed = notFoundAllowed
	}
}

// WithClock sets the clock used by the waiter.
func WithClock(clock clock.Clock) WaiterOption {
	return func(w *waiterConfig) {
		w.clock = clock
	}
}

type waiterConfig struct {
	timeout         time.Duration
	retryInterval   time.Duration
	notFoundAllowed int
	clock           clock.Clock
}
