package backoff_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/backoff"
)

func TestExponentialJitter(t *testing.T) {
	base, maximum := backoff.MinDelay, backoff.MaxDelay
	b := backoff.NewExponentialJitter(base, maximum)
	for i := 0; i <= 1000; i++ {
		delay := b.RetryDelay(i)
		require.True(t, delay >= base && delay <= maximum)
	}
}

func TestExponentialJitterZeroJitter(t *testing.T) {
	b := backoff.NewExponentialJitter(500*time.Millisecond, 30*time.Second, backoff.WithRand(func() float64 {
		return 0
	}))

	expected := []time.Duration{
		500 * time.Millisecond,
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
		16 * time.Second,
		30 * time.Second,
		30 * time.Second,
		30 * time.Second,
		30 * time.Second,
		30 * time.Second,
	}
	actual := make([]time.Duration, 11)
	for i := 0; i <= 10; i++ {
		actual[i] = b.RetryDelay(i)
	}

	require.Equal(t, expected, actual)
}
