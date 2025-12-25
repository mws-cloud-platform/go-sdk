package errors

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/mws/errors"
)

func TestWrapAPIGenError(t *testing.T) {
	tests := []struct {
		name        string
		code        int
		apiGenError APIGenError
		expected    errors.APIError
	}{
		{
			name: "wrap_base_error_with_existing_status",
			code: http.StatusInternalServerError,
			apiGenError: &testGenBaseError{
				code:        ptr.Get("INTERNAL"),
				description: ptr.Get("base error existing status"),
			},
			expected: errors.APIError{
				Code:        http.StatusInternalServerError,
				Status:      errors.Internal,
				Description: "base error existing status",
			},
		},
		{
			name: "wrap_base_error_with_unexpected_status",
			code: http.StatusConflict,
			apiGenError: &testGenBaseError{
				code:        ptr.Get("foobar"),
				description: ptr.Get("base error unexpected status"),
			},
			expected: errors.APIError{
				Code:        http.StatusConflict,
				Status:      errors.Unknown,
				Description: "base error unexpected status\nInvalid status code: \"foobar\"",
			},
		},
		{
			name: "wrap_base_error_with_http_code_into_status",
			code: http.StatusConflict,
			apiGenError: &testGenBaseError{
				code:        ptr.Get("400"),
				description: ptr.Get("base error with http code into status"),
			},
			expected: errors.APIError{
				Code:        http.StatusConflict,
				Status:      errors.Unknown,
				Description: "base error with http code into status\nInvalid status code: \"400\"",
			},
		},
		{
			name: "wrap_base_error_with_unexpected_http_code",
			code: 505,
			apiGenError: &testGenBaseError{
				code:        nil,
				description: ptr.Get("base error with unexpected http code"),
			},
			expected: errors.APIError{
				Code:        505,
				Status:      errors.Unknown,
				Description: "base error with unexpected http code\nInvalid status code: \"\"",
			},
		},
		{
			name: "wrap_api_error_with_existing_status",
			code: http.StatusBadRequest,
			apiGenError: &testGenAPIError{
				code:        ptr.Get("INVALID_ARGUMENT"),
				description: ptr.Get("api error with existing status"),
			},
			expected: errors.APIError{
				Code:        http.StatusBadRequest,
				Status:      errors.InvalidArgument,
				Description: "api error with existing status",
			},
		},
		{
			name: "wrap_api_error_with_details_and_retry_policy",
			code: http.StatusBadRequest,
			apiGenError: &testGenAPIError{
				code:        ptr.Get("IDEMPOTENCY_KEY_ALREADY_USED"),
				description: ptr.Get("api error with details and retry policy"),
				rawDetails: map[string]json.RawMessage{
					"hello": json.RawMessage(`"details"`),
				},
				retryPolicy: NewDefaultRetryPolicy(),
			},
			expected: errors.APIError{
				Code:        http.StatusBadRequest,
				Status:      errors.IdempotencyKeyAlreadyUsed,
				Description: "api error with details and retry policy",
				Details:     map[string]any{"hello": "details"},
				RetryPolicy: NewDefaultRetryPolicy(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapAPIGenError(tt.code, tt.apiGenError)
			require.Equal(t, tt.expected, *result)
		})
	}
}

type testGenBaseError struct {
	code        *string
	description *string
}

func (t testGenBaseError) GetCodeOr(s string) string {
	if t.code != nil {
		return *t.code
	}
	return s
}

func (t testGenBaseError) GetDescriptionOr(s string) string {
	if t.description != nil {
		return *t.description
	}
	return s
}

type testGenAPIError struct {
	code        *string
	description *string
	retryPolicy *errors.RetryPolicy
	rawDetails  map[string]json.RawMessage
}

func (t testGenAPIError) GetCodeOr(s string) string {
	if t.code != nil {
		return *t.code
	}
	return s
}

func (t testGenAPIError) GetDescriptionOr(s string) string {
	if t.description != nil {
		return *t.description
	}
	return s
}

func (t testGenAPIError) GetCommonRetryPolicy() *errors.RetryPolicy {
	return t.retryPolicy
}

func (t testGenAPIError) GetType() *string {
	return nil
}

func (t testGenAPIError) GetDetails() map[string]json.RawMessage {
	return t.rawDetails
}
