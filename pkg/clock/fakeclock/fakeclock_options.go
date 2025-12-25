package fakeclock

import (
	"time"

	"github.com/jonboulle/clockwork"
)

// DefaultFakeClockDelta is the default duration by which
// the fake clock will advance time on each call to clock.Clock.Now.
const DefaultFakeClockDelta = 5 * time.Millisecond

// Option is a constructor option for the [FakeClock] implementation.
type Option func(fc *fakeClock)

// WithAdvancingNow is a fake clock constructor option
// which makes each clock.Clock.Now call automatically
// advance time by [DefaultFakeClockDelta].
func WithAdvancingNow() Option {
	return WithAdvancingNowBy(DefaultFakeClockDelta)
}

// WithAdvancingNowBy is a fake clock constructor option
// which makes each clock.Clock.Now call automatically
// advance time by d.
func WithAdvancingNowBy(d time.Duration) Option {
	return func(fc *fakeClock) {
		fc.deltaAfterNow = d
	}
}

// WithStartAt is a fake clock constructor option
// which sets starting time to the provided one.
func WithStartAt(t time.Time) Option {
	return func(fc *fakeClock) {
		fc.innerClock = clockwork.NewFakeClockAt(t)
	}
}

// WithStartAtCurrentTime is a fake clock constructor option
// which sets starting time to the current system time.
func WithStartAtCurrentTime() Option {
	return func(fc *fakeClock) {
		fc.innerClock = clockwork.NewFakeClock()
	}
}
