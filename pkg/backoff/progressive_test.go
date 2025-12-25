package backoff_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/backoff"
)

func TestProgressiveJitter(t *testing.T) {
	base, maximum := backoff.MinDelay, backoff.MaxDelay
	b := backoff.NewProgressiveJitter(base, maximum)
	for i := 0; i <= 1000; i++ {
		delay := b.RetryDelay(i)
		require.True(t, delay >= base && delay <= maximum)
	}
}

func TestProgressiveJitterZeroJitter(t *testing.T) {
	b := backoff.NewProgressiveJitter(20*time.Millisecond, 30*time.Second, backoff.WithRand(func() float64 {
		return 0
	}))

	expected := []time.Duration{
		20 * time.Millisecond,
		400 * time.Millisecond,
		8 * time.Second,
		30 * time.Second,
		30 * time.Second,
		30 * time.Second,
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
