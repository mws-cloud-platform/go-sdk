package examples_test

import (
	"context"
	"errors"
	"time"

	mwserrors "go.mws.cloud/go-sdk/mws/errors"
)

func wait[T any](ctx context.Context, get func(context.Context) (T, error), check func(T, error) (stop bool, err error)) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			response, err := get(ctx)
			stop, err := check(response, err)
			if err != nil || stop {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func waitForDeletion(ctx context.Context, get func(context.Context) (any, error)) error {
	return wait(ctx, get, func(_ any, err error) (bool, error) {
		if err == nil {
			return false, nil
		}
		var target *mwserrors.APIError
		if errors.As(err, &target) && target.Status == mwserrors.NotFound {
			return true, nil
		}
		return false, err
	})
}
