package backoff_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/backoff"
)

func TestFixed(t *testing.T) {
	b := backoff.NewFixed(backoff.MinDelay)
	for i := 0; i <= 1000; i++ {
		require.Equal(t, backoff.MinDelay, b.RetryDelay(i))
	}
}
