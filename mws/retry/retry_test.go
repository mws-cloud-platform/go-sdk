package retry_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	clienterrors "go.mws.cloud/go-sdk/internal/client/errors"
	mwserrors "go.mws.cloud/go-sdk/mws/errors"
	"go.mws.cloud/go-sdk/mws/retry"
	"go.mws.cloud/go-sdk/pkg/backoff"
)

func TestNoRetry(t *testing.T) {
	d, err := retry.NoRetry{}.RetryDelay(1, net.ErrClosed)
	require.Zero(t, d)
	require.ErrorIs(t, err, clienterrors.ErrNoRetry)
}

func TestStandardRetry(t *testing.T) {
	r := retry.NewStandardRetry(
		retry.WithBackoff(backoff.NewFixed(backoff.MinDelay)),
		retry.WithRetryableCodes(append(retry.DefaultRetryableCodes(), 418)...),
	)

	for _, test := range []struct {
		err    error
		outErr error
	}{
		{
			err:    &mwserrors.APIError{},
			outErr: clienterrors.ErrNoRetry,
		},
		{
			err: &mwserrors.APIError{
				Code: 400,
			},
			outErr: clienterrors.ErrNoRetry,
		},
		{
			err: &mwserrors.APIError{
				Code: 400,
			},
			outErr: nil,
		},
		{
			err: &mwserrors.APIError{
				Code: 418,
			},
			outErr: nil,
		},
		{
			err: &mwserrors.APIError{
				Code: 524,
			},
			outErr: nil,
		},
		{
			err:    &mwserrors.TransportError{},
			outErr: nil,
		},
		{
			err:    errors.New("some error"),
			outErr: clienterrors.ErrNoRetry,
		},
		{
			err:    fmt.Errorf("some error"),
			outErr: clienterrors.ErrNoRetry,
		},
	} {
		delay, err := r.RetryDelay(0, test.err)
		if test.outErr != nil {
			require.ErrorIs(t, err, test.outErr)
			return
		}

		require.NoError(t, err)
		require.Equal(t, backoff.MinDelay, delay)
	}
}
