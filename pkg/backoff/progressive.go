package backoff

import (
	"math"
	"time"
)

// ProgressiveJitter is an implementation of Backoff interface that provides
// backoff delays with jitter. Delay grows in progression.
type ProgressiveJitter struct {
	Base
}

var _ Backoff = (*ProgressiveJitter)(nil)

// NewProgressiveJitter returns a configured ProgressiveJitter.
// Arguments are the same as for NewBase. Using base value of 1ms is not recommended, see RetryDelay.
func NewProgressiveJitter(base, limit time.Duration, options ...Option) ProgressiveJitter {
	return ProgressiveJitter{NewBase(base, limit, options...)}
}

// RetryDelay returns delay before next retry based on the attempt number.
// Formula is: min((baseDelay/ms)^(attempt + 1) * (1 + [0.0..1.0)/2) * ms, maxDelay).
func (b ProgressiveJitter) RetryDelay(attempt int) time.Duration {
	return b.Limit(b.Jitter(math.Pow(float64(b.base/time.Millisecond), float64(attempt)+1)) * float64(time.Millisecond))
}
