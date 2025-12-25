// Package clock provides an interface for time-related operations, allowing for easier testing and mocking.
// It starts with the definition of the [Clock] interface, which includes methods for time manipulation,
// such as creating timers and tickers, getting the current time, and managing context timeouts and deadlines.
// The [Ticker] and [Timer] interfaces are also defined to encapsulate periodic and one-time time events, respectively.
//
// The package is used across the codebase to provide a consistent way to handle time-related operations both in
// production and during testing.
package clock

import (
	"context"
	"time"
)

// Clock is an interface that provides methods for time-related operations. For unification purposes it also
// includes methods for working with [context.Context] timeouts and deadlines.
type Clock interface {
	After(d time.Duration) <-chan time.Time
	NewTicker(d time.Duration) Ticker
	NewTimer(d time.Duration) Timer
	Now() time.Time
	Since(t time.Time) time.Duration
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc)
	WithDeadline(ctx context.Context, t time.Time) (context.Context, context.CancelFunc)
}

// Ticker is an interface that provides methods for periodic time ticks.
type Ticker = interface {
	Chan() <-chan time.Time
	Reset(d time.Duration)
	Stop()
}

// Timer is an interface that provides methods for one-time time events.
type Timer = interface {
	Chan() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}
