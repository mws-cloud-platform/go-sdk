package backoff_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/backoff"
)

func TestBackoff(t *testing.T) {
	for _, test := range []struct {
		name  string
		scale backoff.Scale
		typ   any
		err   error
	}{
		{
			name:  "fixed",
			scale: backoff.ScaleFixed,
			typ:   zero[backoff.Func](),
			err:   nil,
		},
		{
			name:  "linear",
			scale: backoff.ScaleLinear,
			typ:   zero[backoff.LinearJitter](),
			err:   nil,
		},
		{
			name:  "progressive",
			scale: backoff.ScaleProgressive,
			typ:   zero[backoff.ProgressiveJitter](),
			err:   nil,
		},
		{
			name:  "exponential",
			scale: backoff.ScaleExponential,
			typ:   zero[backoff.ExponentialJitter](),
			err:   nil,
		},
		{
			name:  "invalid",
			scale: backoff.ScaleInvalid,
			typ:   nil,
			err:   backoff.ErrInvalidScale,
		},
		{
			name:  "unknown",
			scale: 42,
			typ:   nil,
			err:   backoff.ErrUnknownScale,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			b, err := backoff.New(test.scale, backoff.MinDelay, backoff.MinDelay)
			if test.err == nil {
				require.NoError(t, err)
				return
			}

			require.ErrorIs(t, err, test.err)
			require.IsType(t, test.typ, b)
		})
	}
}

func zero[T any]() T {
	var z T
	return z
}
