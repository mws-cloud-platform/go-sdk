package wait_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	commonerrors "go.mws.cloud/go-sdk/mws/errors"
	"go.mws.cloud/go-sdk/mws/wait"
	"go.mws.cloud/go-sdk/pkg/clock/fakeclock"
)

var ErrNotFound = &commonerrors.APIError{
	Status: commonerrors.NotFound,
}

func TestWaiterOK(t *testing.T) {
	var (
		actual string
		err    error
	)

	clock := fakeclock.NewFake()

	callback := func(_ context.Context) (string, bool, error) {
		return "ok", true, nil
	}
	ctx := t.Context()

	w := wait.NewWaiter(callback, wait.WithClock(clock))

	done := make(chan struct{})
	go func() {
		actual, err = w.Wait(ctx)
		close(done)
	}()

	clock.BlockUntilContext(ctx, 2)

	clock.Advance(wait.DefaultRetryInterval)

	<-done
	require.NoError(t, err)
	require.Equal(t, "ok", actual)
}

func TestWaiterOKAfterRetry(t *testing.T) {
	var (
		actual string
		err    error
	)

	clock := fakeclock.NewFake()
	times := 3
	count := 0
	called := make(chan struct{})
	callback := func(_ context.Context) (string, bool, error) {
		count++
		called <- struct{}{}
		if count < times {
			return "processing", false, nil
		}
		return "ok", true, nil
	}
	ctx := t.Context()

	w := wait.NewWaiter(callback, wait.WithClock(clock))

	done := make(chan struct{})
	go func() {
		actual, err = w.Wait(ctx)
		close(done)
	}()

	clock.BlockUntilContext(ctx, 2)

	for range times {
		clock.Advance(wait.DefaultRetryInterval)
		<-called
	}

	<-done
	require.NoError(t, err)
	require.Equal(t, "ok", actual)
}

func TestWaiterError(t *testing.T) {
	var err error

	clock := fakeclock.NewFake()
	callback := func(_ context.Context) (string, bool, error) {
		return "", true, errors.New("fail")
	}
	ctx := t.Context()

	w := wait.NewWaiter(callback, wait.WithClock(clock))

	done := make(chan struct{})
	go func() {
		_, err = w.Wait(ctx)
		close(done)
	}()

	clock.BlockUntilContext(ctx, 2)
	clock.Advance(wait.DefaultRetryInterval)

	<-done
	require.EqualError(t, err, "fail")
}

func TestWaiterTimeout(t *testing.T) {
	var err error

	clock := fakeclock.NewFake()
	called := make(chan struct{})
	callback := func(_ context.Context) (string, bool, error) {
		called <- struct{}{}
		return "processing", false, nil
	}
	ctx := t.Context()

	w := wait.NewWaiter(callback, wait.WithClock(clock))

	done := make(chan struct{})
	go func() {
		_, err = w.Wait(ctx)
		close(done)
	}()

	clock.BlockUntilContext(ctx, 2)
	for {
		clock.Advance(wait.DefaultRetryInterval)
		select {
		case <-called:
		case <-done:
			require.ErrorIs(t, err, context.DeadlineExceeded)
			return
		}
	}
}

func TestWaiterCancel(t *testing.T) {
	var err error

	clock := fakeclock.NewFake()
	called := make(chan struct{})
	callback := func(ctx context.Context) (string, bool, error) {
		select {
		case called <- struct{}{}:
			return "processing", false, nil
		case <-ctx.Done():
			return "", false, ctx.Err()
		}
	}
	ctx, cancel := context.WithCancel(t.Context())

	w := wait.NewWaiter(callback, wait.WithClock(clock))

	done := make(chan struct{})
	go func() {
		_, err = w.Wait(ctx)
		close(done)
	}()

	clock.BlockUntilContext(ctx, 2)
	for range 3 {
		clock.Advance(wait.DefaultRetryInterval)
		<-called
	}

	cancel()

	<-done
	require.ErrorIs(t, err, context.Canceled)
}

func TestWaiterNotFound(t *testing.T) {
	var (
		actual string
		err    error
	)

	clock := fakeclock.NewFake()
	called := make(chan struct{})
	callback := func(_ context.Context) (string, bool, error) {
		called <- struct{}{}
		return "not_found", false, ErrNotFound
	}
	ctx := t.Context()

	w := wait.NewWaiter(callback, wait.WithClock(clock))

	done := make(chan struct{})
	go func() {
		actual, err = w.Wait(ctx)
		close(done)
	}()

	clock.BlockUntilContext(ctx, 2)

	for range wait.DefaultNotFoundAllowed + 1 {
		clock.Advance(wait.DefaultRetryInterval)
		<-called
	}

	<-done
	require.EqualError(t, err, ErrNotFound.Error())
	require.Equal(t, "not_found", actual)
}
