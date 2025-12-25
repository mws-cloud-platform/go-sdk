package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var errInner = errors.New("inner")

func TestErrors_Comparison(t *testing.T) {
	for _, test := range []struct {
		name      string
		skipInner bool
		newError  func() error
		asError   func(error) bool
		errorsIs  func(error) bool
		customIs  func(error) bool
	}{
		{
			name:      "APIError",
			skipInner: true,
			newError: func() error {
				return NewAPIError(200, Unknown, "description")
			},
			asError: func(err error) bool {
				var target *APIError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &APIError{})
			},
			customIs: IsAPIError,
		},
		{
			name: "TransportError",
			newError: func() error {
				return NewTransportError(errInner)
			},
			asError: func(err error) bool {
				var target *TransportError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &TransportError{})
			},
			customIs: IsTransportError,
		},
		{
			name: "EncodeBodyError",
			newError: func() error {
				return NewEncodeBodyError(errInner)
			},
			asError: func(err error) bool {
				var target *EncodeBodyError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &EncodeBodyError{})
			},
			customIs: IsEncodeBodyError,
		},
		{
			name: "DecodeBodyError",
			newError: func() error {
				return NewDecodeBodyError("application/json", errInner)
			},
			asError: func(err error) bool {
				var target *DecodeBodyError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &DecodeBodyError{})
			},
			customIs: IsDecodeBodyError,
		},
		{
			name:      "InvalidContentTypeError",
			skipInner: true,
			newError: func() error {
				return NewInvalidContentTypeError("content-type")
			},
			asError: func(err error) bool {
				var target *InvalidContentTypeError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &InvalidContentTypeError{})
			},
			customIs: IsInvalidContentTypeError,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mwsError := test.newError()

			if !test.skipInner {
				require.ErrorIs(t, errors.Unwrap(mwsError), errInner)
				require.True(t, errors.Is(mwsError, errInner))
			}

			require.False(t, test.asError(errInner))
			require.False(t, test.errorsIs(errInner))
			require.False(t, test.customIs(errInner))

			require.True(t, test.asError(mwsError))
			require.True(t, test.errorsIs(mwsError))
			require.True(t, test.customIs(mwsError))

			wrapped := fmt.Errorf("%w", mwsError)
			require.True(t, test.errorsIs(wrapped))
			require.True(t, test.customIs(wrapped))
		})
	}
}

func TestErrors_Message(t *testing.T) {
	for _, test := range []struct {
		name          string
		newErrorInner func() error
		newErrorNil   func() error
		expectedInner string
		expectedNil   string
	}{
		{
			name: "APIError",
			newErrorNil: func() error {
				return NewAPIError(200, AlreadyExists, "description")
			},
			expectedNil: "api error. Code: 200. Status: ALREADY_EXISTS. Description: description",
		},
		{
			name: "TransportError",
			newErrorInner: func() error {
				return NewTransportError(errInner)
			},
			newErrorNil: func() error {
				return NewTransportError(nil)
			},
			expectedInner: "transport error: inner",
			expectedNil:   "transport error",
		},
		{
			name: "EncodeBodyError",
			newErrorInner: func() error {
				return NewEncodeBodyError(errInner)
			},
			newErrorNil: func() error {
				return NewEncodeBodyError(nil)
			},
			expectedInner: "encode error: inner",
			expectedNil:   "encode error",
		},
		{
			name: "DecodeBodyError",
			newErrorInner: func() error {
				return NewDecodeBodyError("content-type", errInner)
			},
			newErrorNil: func() error {
				return NewDecodeBodyError("content-type", nil)
			},
			expectedInner: "decode content type 'content-type': inner",
			expectedNil:   "decode content type 'content-type'",
		},
		{
			name: "InvalidContentTypeError",
			newErrorInner: func() error {
				return NewInvalidContentTypeError("content-type")
			},
			expectedInner: "invalid Content-Type: 'content-type'",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.newErrorInner != nil {
				mwsErr := test.newErrorInner()
				require.Equal(t, test.expectedInner, fmt.Sprintf("%s", mwsErr))
			}

			if test.newErrorNil != nil {
				mwsErr := test.newErrorNil()
				require.Equal(t, test.expectedNil, fmt.Sprintf("%s", mwsErr))
			}
		})
	}
}
