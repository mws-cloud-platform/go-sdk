package errors

import (
	"errors"
	"fmt"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"

	mwserrors "go.mws.cloud/go-sdk/mws/errors"
)

const (
	// ErrNoRetry is the error returned when no retry is allowed.
	ErrNoRetry = consterr.Error("no retry")
)

type DecodeBodyError = mwserrors.DecodeBodyError

func NewDecodeBodyError(contentType string, err error) *DecodeBodyError {
	return mwserrors.NewDecodeBodyError(contentType, err)
}

func IsDecodeBodyError(err error) bool {
	return mwserrors.IsDecodeBodyError(err)
}

type InvalidContentTypeError = mwserrors.InvalidContentTypeError

func InvalidContentType(contentType string) *InvalidContentTypeError {
	return mwserrors.NewInvalidContentTypeError(contentType)
}

func IsInvalidContentTypeError(err error) bool {
	return mwserrors.IsInvalidContentTypeError(err)
}

type UnexpectedStatusCodeError = mwserrors.UnexpectedStatusCodeError

func UnexpectedStatusCodeWithData(statusCode int, data []byte) *UnexpectedStatusCodeError {
	return mwserrors.NewUnexpectedStatusCodeErrorWithData(statusCode, data)
}

func IsUnexpectedStatusCodeError(err error) bool {
	return mwserrors.IsUnexpectedStatusCodeError(err)
}

// RetryAttemptsExhaustedError is the error returned when no retry attempts left.
type RetryAttemptsExhaustedError struct {
	Attempts int
	Err      error
}

func NewRetryAttemptsExhaustedError(attempts int, err error) RetryAttemptsExhaustedError {
	return RetryAttemptsExhaustedError{
		Attempts: attempts,
		Err:      err,
	}
}

func (e RetryAttemptsExhaustedError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("failed to execute request after %d attempts", e.Attempts)
	}
	return fmt.Sprintf("failed to execute request after %d attempts: %v", e.Attempts, e.Err)
}

func (e RetryAttemptsExhaustedError) Unwrap() error {
	return e.Err
}

func (e RetryAttemptsExhaustedError) Is(err error) bool {
	return IsRetryAttemptsExhaustedError(err)
}

func IsRetryAttemptsExhaustedError(err error) bool {
	return errors.As(err, &RetryAttemptsExhaustedError{})
}
