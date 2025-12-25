// Package fakeclock provides a fake clock implementation for testing purposes.
package fakeclock

import (
	"context"
	"time"

	"github.com/jonboulle/clockwork"

	"go.mws.cloud/go-sdk/pkg/clock"
)

// StartTime is the default start time for the fake clock.
var StartTime = time.Date(2025, time.January, 9, 10, 20, 0, 0, time.UTC)

type innerClock = clockwork.FakeClock

// FakeClock is a fake clock implementation that can be used in tests to control time flow.
type FakeClock interface {
	clock.Clock

	// Advance advances [FakeClock] to a new point in time, ensuring waiters and
	// blockers are notified appropriately before returning.
	Advance(d time.Duration)

	// BlockUntilContext blocks until the [FakeClock] has the given number of waiters
	// or the context is cancelled.
	BlockUntilContext(ctx context.Context, n int) error
}

// NewFake creates a fake clock.
//
// If no start time option is provided, it is
// initialized with [StartTime] constant.
func NewFake(opts ...Option) FakeClock {
	fc := &fakeClock{}
	fc.applyOptions(opts)
	return fc
}

type fakeClock struct {
	*innerClock
	deltaAfterNow time.Duration
}

func (fc *fakeClock) applyOptions(opts []Option) {
	for _, opt := range opts {
		opt(fc)
	}

	if fc.innerClock == nil {
		fc.innerClock = clockwork.NewFakeClockAt(StartTime)
	}
}

func (fc *fakeClock) NewTicker(d time.Duration) clock.Ticker { return fc.innerClock.NewTicker(d) }
func (fc *fakeClock) NewTimer(d time.Duration) clock.Timer   { return fc.innerClock.NewTimer(d) }
func (fc *fakeClock) Tick(d time.Duration) <-chan time.Time  { return fc.NewTicker(d).Chan() }

func (fc *fakeClock) WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return clockwork.WithTimeout(ctx, fc.innerClock, d)
}

func (fc *fakeClock) WithDeadline(ctx context.Context, t time.Time) (context.Context, context.CancelFunc) {
	return clockwork.WithDeadline(ctx, fc.innerClock, t)
}

func (fc *fakeClock) Now() time.Time {
	now := fc.innerClock.Now()
	if fc.deltaAfterNow != 0 {
		fc.Advance(fc.deltaAfterNow)
	}
	return now
}
