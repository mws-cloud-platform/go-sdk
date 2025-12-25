// Package backoff provides backoff mechanisms for retrying operations.
package backoff

import (
	"fmt"
	"math/rand/v2"
	"time"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

const (
	MinDelay = 20 * time.Millisecond
	MaxDelay = 10 * time.Second

	ErrUnknownScale = consterr.Error("unknown scale")
	ErrInvalidScale = consterr.Error("invalid scale")
)

// Scale defines timeout increase formula.
// Fixed – fixed, timeout doesn't change.
// Linear – t, 2t, 3t, ..., n*t.
// Progressive – t, t^2, t^3, ..., t^n.
// Exponential – t, t*2, t*2^2, ..., t*2^n.
type Scale int

const (
	ScaleInvalid Scale = iota - 1
	ScaleFixed
	ScaleLinear
	ScaleProgressive
	ScaleExponential
)

var (
	scaleToString = map[Scale]string{
		ScaleInvalid:     "INVALID",
		ScaleFixed:       "FIXED",
		ScaleLinear:      "LINEAR",
		ScaleProgressive: "PROGRESSIVE",
		ScaleExponential: "EXPONENTIAL",
	}
)

// Backoff represents a backoff mechanism that returns delay before next retry happens.
type Backoff interface {
	RetryDelay(attempt int) time.Duration
}

// Func is an implementation of Backoff interface suitable for simple scenarios.
type Func func(int) time.Duration

func (f Func) RetryDelay(attempt int) time.Duration {
	return f(attempt)
}

// New returns backoff implementation based on the Scale type and basic configuration.
// It returns an error when an invalid Scale is passed.
func New(scale Scale, base, limit time.Duration) (Backoff, error) {
	switch scale {
	case ScaleFixed:
		return NewFixed(base), nil
	case ScaleLinear:
		return NewLinearJitter(base, limit), nil
	case ScaleProgressive:
		return NewProgressiveJitter(base, limit), nil
	case ScaleExponential:
		return NewExponentialJitter(base, limit), nil
	case ScaleInvalid:
		return nil, ErrInvalidScale
	default:
		return nil, fmt.Errorf("%w: %d (%s)", ErrUnknownScale, scale, scaleToString[scale])
	}
}

// Base is a base entity suitable for most backoff implementations.
// It is expected to be embedded, not to be used directly.
type Base struct {
	base  time.Duration
	limit time.Duration
	rand  func() float64
}

// NewBase returns a configured Base.
// base – starting value of delay between retries.
// limit – maximum value of delay between retries.
func NewBase(base, limit time.Duration, options ...Option) Base {
	b := Base{
		base:  base,
		limit: limit,
		rand:  rand.Float64,
	}
	for _, option := range options {
		option(&b)
	}

	return b
}

func (b Base) Jitter(base float64) float64 {
	return base * (1 + b.rand()/2)
}

func (b Base) Limit(delay float64) time.Duration {
	return time.Duration(min(delay, float64(b.limit)))
}

// Option is a functional option type that extends Base functionality
type Option func(*Base)

// WithRand returns an option that overrides default randomizer function
func WithRand(rand func() float64) Option {
	return func(e *Base) {
		e.rand = rand
	}
}
