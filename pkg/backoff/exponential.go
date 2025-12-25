package backoff

import (
	"math"
	"time"
)

const (
	exponentialBase = 2
)

// ExponentialJitter is an implementation of Backoff interface that provides
// backoff delays with jitter. Delay grows exponentially.
type ExponentialJitter struct {
	Base
}

var _ Backoff = (*ExponentialJitter)(nil)

// NewExponentialJitter returns a configured ExponentialJitter.
// Arguments are the same as for NewBase.
func NewExponentialJitter(base, limit time.Duration, options ...Option) ExponentialJitter {
	return ExponentialJitter{NewBase(base, limit, options...)}
}

// RetryDelay returns delay before next retry based on the attempt number.
// Formula is: min(exponentialBase^attempt * baseDelay * (1 + [0.0..1.0)/2), maxDelay).
func (b ExponentialJitter) RetryDelay(attempt int) time.Duration {
	return b.Limit(b.Jitter(math.Pow(exponentialBase, float64(attempt)) * float64(b.base)))
}
