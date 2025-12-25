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

type fakeRetryer struct{}

func (fakeRetryer) RetryDelay(int, error) (time.Duration, error) {
	return 0, nil
}
