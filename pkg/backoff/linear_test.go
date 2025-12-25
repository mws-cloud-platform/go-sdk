package backoff_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/backoff"
)

func TestLinearJitter(t *testing.T) {
	base, maximum := backoff.MinDelay, backoff.MaxDelay
	b := backoff.NewLinearJitter(base, maximum)
	for i := 0; i <= 1000; i++ {
		delay := b.RetryDelay(i)
		require.True(t, delay >= base && delay <= maximum)
	}
}

func TestLinearJitterZeroJitter(t *testing.T) {
	b := backoff.NewLinearJitter(5*time.Second, 30*time.Second, backoff.WithRand(func() float64 {
		return 0
	}))

	expected := []time.Duration{
		5 * time.Second,
		10 * time.Second,
		15 * time.Second,
		20 * time.Second,
		25 * time.Second,
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
