// Package retry provides implementations for retry mechanisms.
package retry

import (
	"errors"
	"net/http"
	"time"

	clienterrors "go.mws.cloud/go-sdk/internal/client/errors"
	mwserrors "go.mws.cloud/go-sdk/mws/errors"
	"go.mws.cloud/go-sdk/pkg/backoff"
)

// Retryer represents a retry mechanism that determines whether an operation
// should be retried on error and, if so, with what delay.
type Retryer interface {
	RetryDelay(attempt int, err error) (time.Duration, error)
}

// NoRetry is a retryer that always decides to not retry errors.
type NoRetry struct{}

// RetryDelay always returns an error.
func (NoRetry) RetryDelay(int, error) (time.Duration, error) {
	return 0, clienterrors.ErrNoRetry
}

const (
	httpStatusWebServerDown      = 521
	httpStatusConnectionTimedOut = 522
	httpStatusTimeoutOccurred    = 524
)

var defaultRetryableCodes = map[int]struct{}{
	http.StatusRequestTimeout:      {}, // 408
	http.StatusConflict:            {}, // 409
	http.StatusLocked:              {}, // 423
	http.StatusFailedDependency:    {}, // 424
	http.StatusTooManyRequests:     {}, // 429
	http.StatusInternalServerError: {}, // 500
	http.StatusBadGateway:          {}, // 502
	http.StatusServiceUnavailable:  {}, // 503
	http.StatusGatewayTimeout:      {}, // 504
	httpStatusWebServerDown:        {}, // 521
	httpStatusConnectionTimedOut:   {}, // 522
	httpStatusTimeoutOccurred:      {}, // 524
}

func DefaultRetryableCodes() []int {
	keys := make([]int, 0, len(defaultRetryableCodes))
	for k := range defaultRetryableCodes {
		keys = append(keys, k)
	}

	return keys
}

type StandardRetry struct {
	backoff        backoff.Backoff
	retryableCodes map[int]struct{}
}

func NewStandardRetry(options ...StandardRetryOption) *StandardRetry {
	s := &StandardRetry{
		backoff:        backoff.NewExponentialJitter(backoff.MinDelay, backoff.MaxDelay),
		retryableCodes: defaultRetryableCodes,
	}
	for _, option := range options {
		option(s)
	}

	return s
}

func (s *StandardRetry) RetryDelay(attempt int, err error) (time.Duration, error) {
	if !s.IsRetryable(err) {
		return 0, clienterrors.ErrNoRetry
	}

	return s.backoff.RetryDelay(attempt), nil
}

func (s *StandardRetry) IsRetryable(err error) bool {
	switch {
	case mwserrors.IsAPIError(err):
		var apiError *mwserrors.APIError
		errors.As(err, &apiError)

		return s.codeIsRetryable(apiError.Code)
	case mwserrors.IsTransportError(err):
		return true
	default:
		return false
	}
}

func (s *StandardRetry) codeIsRetryable(code int) bool {
	_, ok := s.retryableCodes[code]
	return ok
}

type StandardRetryOption func(*StandardRetry)

func WithBackoff(backoff backoff.Backoff) StandardRetryOption {
	return func(r *StandardRetry) {
		r.backoff = backoff
	}
}

func WithRetryableCodes(retryableCodes ...int) StandardRetryOption {
	return func(r *StandardRetry) {
		r.retryableCodes = make(map[int]struct{}, len(retryableCodes))
		for _, code := range retryableCodes {
			r.retryableCodes[code] = struct{}{}
		}
	}
}
