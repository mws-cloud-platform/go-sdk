package backoff

import (
	"time"
)

// LinearJitter is an implementation of Backoff interface that provides
// backoff delays with jitter. Delay grows linearly.
type LinearJitter struct {
	Base
}

var _ Backoff = (*LinearJitter)(nil)

// NewLinearJitter returns a configured LinearJitter.
// Arguments are the same as for NewBase.
func NewLinearJitter(base, limit time.Duration, options ...Option) LinearJitter {
	return LinearJitter{NewBase(base, limit, options...)}
}

// RetryDelay returns delay before next retry based on the attempt number.
// Formula is: min((attempt + 1) * baseDelay * (1 + [0.0..1.0)/2), maxDelay).
func (b LinearJitter) RetryDelay(attempt int) time.Duration {
	return b.Limit(b.Jitter(float64(attempt+1) * float64(b.base)))
}
