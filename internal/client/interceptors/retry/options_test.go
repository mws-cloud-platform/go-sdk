package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithRetryer(t *testing.T) {
	cfg := defaultConfig()
	WithRetryer(fakeRetryer{}).apply(cfg)
	assert.Equal(t, fakeRetryer{}, cfg.retryer)
}

func TestWithMaxAttempts(t *testing.T) {
	cfg := defaultConfig()
	WithMaxAttempts(5).apply(cfg)
	assert.Equal(t, 5, cfg.maxAttempts)

	WithMaxAttempts(0).apply(cfg)
	assert.Equal(t, 5, cfg.maxAttempts)

	WithMaxAttempts(-1).apply(cfg)
	assert.Equal(t, 5, cfg.maxAttempts)
}

type fakeRetryer struct{}

func (fakeRetryer) RetryDelay(int, error) (time.Duration, error) {
	return 0, nil
}
