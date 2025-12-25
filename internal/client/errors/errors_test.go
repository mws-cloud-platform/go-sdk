package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var errInner = errors.New("inner")

func TestClientErrors_Comparison(t *testing.T) {
	for _, test := range []struct {
		name      string
		skipInner bool
		newError  func() error
		asError   func(error) bool
		errorsIs  func(error) bool
		customIs  func(error) bool
	}{
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
				return InvalidContentType("content-type")
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
		{
			name:      "UnexpectedStatusCode",
			skipInner: true,
			newError: func() error {
				return UnexpectedStatusCodeWithData(418, nil)
			},
			asError: func(err error) bool {
				var target *UnexpectedStatusCodeError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &UnexpectedStatusCodeError{})
			},
			customIs: IsUnexpectedStatusCodeError,
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

func TestClientErrors_Message(t *testing.T) {
	for _, test := range []struct {
		name          string
		newErrorInner func() error
		newErrorNil   func() error
		expectedInner string
		expectedNil   string
	}{
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
				return InvalidContentType("content-type")
			},
			expectedInner: "invalid Content-Type: 'content-type'",
		},
		{
			name: "UnexpectedStatusCodeError",
			newErrorInner: func() error {
				return UnexpectedStatusCodeWithData(418, nil)
			},
			expectedInner: "unexpected status code: '418'",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mwsErr := test.newErrorInner()
			require.Equal(t, test.expectedInner, fmt.Sprintf("%s", mwsErr))

			if test.newErrorNil != nil {
				mwsErr = test.newErrorNil()
				require.Equal(t, test.expectedNil, fmt.Sprintf("%s", mwsErr))
			}
		})
	}
}
