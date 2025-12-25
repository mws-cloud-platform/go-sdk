package clock

import (
	"context"
	"time"
)

// NewReal function returns a [Clock] that simply delegates calls to the actual time
// package; it should be used by packages in production.
func NewReal() Clock { return realClock{} }

type realClock struct{}

func (realClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
func (realClock) NewTicker(d time.Duration) Ticker       { return realTicker{time.NewTicker(d)} }
func (realClock) NewTimer(d time.Duration) Timer         { return realTimer{time.NewTimer(d)} }
func (realClock) Now() time.Time                         { return time.Now() }
func (realClock) Since(t time.Time) time.Duration        { return time.Since(t) }
func (realClock) Sleep(d time.Duration)                  { time.Sleep(d) }
func (realClock) Tick(d time.Duration) <-chan time.Time  { return time.Tick(d) }
func (realClock) WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}

func (realClock) WithDeadline(ctx context.Context, t time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(ctx, t)
}

type realTicker struct{ *time.Ticker }

func (r realTicker) Chan() <-chan time.Time { return r.C }

type realTimer struct{ *time.Timer }

func (r realTimer) Chan() <-chan time.Time { return r.C }
