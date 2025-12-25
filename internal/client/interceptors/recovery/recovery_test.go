package recovery

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/internal/client"
)

func TestRecovery(t *testing.T) {
	var calls int
	err := Recovery(t.Context(), nil, nil, func(context.Context, any, client.APIResp) error {
		calls++
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 1, calls)
}

func TestRecovery_panic(t *testing.T) {
	err := Recovery(t.Context(), nil, nil, func(context.Context, any, client.APIResp) error {
		panic("oops!")
	})

	require.True(t, IsPanicError(err))
	require.EqualError(t, err, "panic: oops!")
}

func TestRecovery_panicError(t *testing.T) {
	err := Recovery(t.Context(), nil, nil, func(context.Context, any, client.APIResp) error {
		panic(errors.New("fail"))
	})

	require.True(t, IsPanicError(err))
	require.EqualError(t, err, "panic: fail")
}
