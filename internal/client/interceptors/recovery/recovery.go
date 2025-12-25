package recovery

import (
	"context"
	"errors"
	"fmt"

	commonclient "go.mws.cloud/go-sdk/internal/client"
)

// Recovery is a client interceptor that recovers from panics and converts them
// into errors.
func Recovery(
	ctx context.Context,
	request any,
	response commonclient.APIResp,
	invoker commonclient.Invoker,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = PanicError{r}
		}
	}()
	return invoker(ctx, request, response)
}

// PanicError is an error that wraps panic.
type PanicError struct {
	reason any
}

func (e PanicError) Error() string {
	return fmt.Sprintf("panic: %v", e.reason)
}

func (e PanicError) Unwrap() error {
	if err, ok := e.reason.(error); ok {
		return err
	}
	return nil
}

func (e PanicError) Is(err error) bool {
	return IsPanicError(err)
}

// IsPanicError reports whether provided error matches [PanicError].
func IsPanicError(err error) bool {
	var target PanicError
	return errors.As(err, &target)
}
