package backoff

import (
	"time"
)

// NewFixed returns a Func that returns fixed delay for every attempt.
func NewFixed(delay time.Duration) Func {
	return func(int) time.Duration {
		return delay
	}
}
