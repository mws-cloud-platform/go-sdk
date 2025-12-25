package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	errInner                     = fmt.Errorf("err")
	errInnerPathAccumulatorError = NewPathAccumulatorError("inner", errInner)
	rawPath                      = "projects/1"
	path                         = "[0]"
)

func TestErrors_Comparison(t *testing.T) {
	testCases := []struct {
		name      string
		skipInner bool
		newError  func() error
		asError   func(error) bool
		errorsIs  func(error) bool
		customIs  func(error) bool
	}{
		{
			name: "ParseReferenceError",
			newError: func() error {
				return NewParseReferenceError(rawPath, errInner)
			},
			asError: func(err error) bool {
				var target *ParseReferenceError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &ParseReferenceError{})
			},
			customIs: IsParseReferenceError,
		},
		{
			name: "ParseIDError",
			newError: func() error {
				return NewParseIDError(rawPath, errInner)
			},
			asError: func(err error) bool {
				var target *ParseIDError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &ParseIDError{})
			},
			customIs: IsParseIDError,
		},
		{
			name: "InitOneOfReferenceError",
			newError: func() error {
				return NewInitOneOfReferenceError("")
			},
			asError: func(err error) bool {
				var target *InitOneOfReferenceError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &InitOneOfReferenceError{})
			},
			customIs: IsInitOneOfReferenceError,
		},
		{
			name: "InitOneOfIDError",
			newError: func() error {
				return NewInitOneOfIDError("")
			},
			asError: func(err error) bool {
				var target *InitOneOfIDError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &InitOneOfIDError{})
			},
			customIs: IsInitOneOfIDError,
		},
		{
			name: "PathAccumulatorError",
			newError: func() error {
				return NewPathAccumulatorError(path, errInner)
			},
			asError: func(err error) bool {
				var target *PathAccumulatorError
				return errors.As(err, &target)
			},
			errorsIs: func(err error) bool {
				return errors.Is(err, &PathAccumulatorError{})
			},
			customIs: IsPathAccumulatorError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testErr := testCase.newError()

			require.False(t, testCase.asError(errInner))
			require.False(t, testCase.errorsIs(errInner))
			require.False(t, testCase.customIs(errInner))

			require.True(t, testCase.asError(testErr))
			require.True(t, testCase.errorsIs(testErr))
			require.True(t, testCase.customIs(testErr))

			wrapped := fmt.Errorf("%w", testErr)
			require.True(t, testCase.errorsIs(wrapped))
			require.True(t, testCase.customIs(wrapped))
		})
	}
}

func TestErrors_Message(t *testing.T) {
	for _, test := range []struct {
		name                              string
		newErrorInnerPathAccumulatorError func() error
		newErrorInner                     func() error
		newErrorNil                       func() error
		expectedInnerPathAccumulatorError string
		expectedInner                     string
		expectedNil                       string
	}{
		{
			name: "ParseReferenceError",
			newErrorInner: func() error {
				return NewParseReferenceError(rawPath, errInner)
			},
			newErrorNil: func() error {
				return NewParseReferenceError(rawPath, nil)
			},
			expectedInner: "parse reference 'projects/1': err",
			expectedNil:   "parse reference 'projects/1'",
		},
		{
			name: "ParseIDError",
			newErrorInner: func() error {
				return NewParseIDError(rawPath, errInner)
			},
			newErrorNil: func() error {
				return NewParseIDError(rawPath, nil)
			},
			expectedInner: "parse id 'projects/1': err",
			expectedNil:   "parse id 'projects/1'",
		},
		{
			name: "InitOneOfReferenceError",
			newErrorInner: func() error {
				return NewInitOneOfReferenceError("")
			},
			expectedInner: "type 'string' is not supported for init a one of reference",
		},
		{
			name: "InitOneOfIDError",
			newErrorInner: func() error {
				return NewInitOneOfIDError("")
			},
			expectedInner: "type 'string' is not supported for init a one of id",
		},
		{
			name: "PathAccumulatorError",
			newErrorInnerPathAccumulatorError: func() error {
				return NewPathAccumulatorError(path, errInnerPathAccumulatorError)
			},
			newErrorInner: func() error {
				return NewPathAccumulatorError(path, errInner)
			},
			newErrorNil: func() error {
				return NewPathAccumulatorError(path, nil)
			},

			// errInnerPathAccumulatorError
			expectedInnerPathAccumulatorError: "path '[0].inner': err",
			expectedInner:                     "path '[0]': err",
			expectedNil:                       "path '[0]'",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mwsErr := test.newErrorInner()
			require.Equal(t, test.expectedInner, fmt.Sprintf("%s", mwsErr))

			if test.newErrorInnerPathAccumulatorError != nil {
				mwsErr = test.newErrorInnerPathAccumulatorError()
				require.Equal(t, test.expectedInnerPathAccumulatorError, fmt.Sprintf("%s", mwsErr))
			}
			if test.newErrorNil != nil {
				mwsErr = test.newErrorNil()
				require.Equal(t, test.expectedNil, fmt.Sprintf("%s", mwsErr))
			}
		})
	}
}
