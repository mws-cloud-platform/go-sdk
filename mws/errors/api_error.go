package errors

import (
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.mws.cloud/go-sdk/pkg/backoff"
)

// NewAPIError creates a new API error.
func NewAPIError(code int, status Status, description string) *APIError {
	return &APIError{
		Code:        code,
		Status:      status,
		Description: description,
	}
}

// APIError represents an error response from the API. Returned by the client
// when the API returns an error response.
type APIError struct {
	// HTTP status code (e.g., 400, 404, 500).
	Code int
	// Typed error status for programmatic handling.
	Status Status
	// Human-readable error description.
	Description string
	// Retry policy suggested by the API.
	RetryPolicy *RetryPolicy
	// Additional human-readable error context. Not suitable for automated
	// processing.
	Details Details
}

// Error implements the error interface for [APIError]. Constructs a formatted
// error string based on available fields.
func (a *APIError) Error() string {
	if a == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("api error. Code: " + strconv.Itoa(a.Code))
	status := a.StatusString()
	if status != "" {
		sb.WriteString(". Status: " + status)
	}
	if a.Description != "" {
		sb.WriteString(". Description: " + a.Description)
	}
	return sb.String()
}

// StatusString returns the string representation of the error status.
func (a *APIError) StatusString() string {
	if a == nil {
		return ""
	}
	strStatus, ok := statusToStringMap[a.Status]
	if !ok {
		return ""
	}
	return strStatus
}

// Is checks if the given error is an [APIError]. Implements the interface for
// error comparison using [errors.Is].
func (a *APIError) Is(err error) bool {
	return IsAPIError(err)
}

// IsAPIError checks if the error is or wraps an [APIError].
func IsAPIError(err error) bool {
	var target *APIError
	return errors.As(err, &target)
}

// Details stores additional error context as key-value pairs. Intended for
// human-readable information, not machine processing.
type Details map[string]any

// String returns a flattened JSON representation of the details.
func (d Details) String() string {
	return flattenJSON(d, "")
}

// RetryPolicy defines retry behavior for API errors. Used by clients to
// determine when and how to retry failed requests.
type RetryPolicy struct {
	// Maximum total time allowed for all retry attempts.
	MaxRetryTimeout time.Duration
	// Maximum number of retry attempts.
	RetryCount int
	// Base timeout between retry attempts.
	RetryTimeout time.Duration
	// Multiplier for scaling timeout between consecutive retries.
	RetryTimeoutScale backoff.Scale
}

// IsAPIErrorAlreadyExistsStatus checks if the error is an [APIError] with
// [AlreadyExists] status. Returns false if the error is not an [APIError] or
// has a different status.
func IsAPIErrorAlreadyExistsStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == AlreadyExists
}

// IsAPIErrorCancelledStatus checks if the error is an [APIError] with
// [Cancelled] status. Returns false if the error is not an [APIError] or has a
// different status.
func IsAPIErrorCancelledStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == Cancelled
}

// IsAPIErrorDeadlineExceededStatus checks if the error is an [APIError] with
// [DeadlineExceeded] status. Returns false if the error is not an [APIError] or
// has a different status.
func IsAPIErrorDeadlineExceededStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == DeadlineExceeded
}

// IsAPIErrorFailedPreconditionStatus checks if the error is an [APIError] with
// [FailedPrecondition] status. Returns false if the error is not an [APIError]
// or has a different status.
func IsAPIErrorFailedPreconditionStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == FailedPrecondition
}

// IsAPIErrorInternalStatus checks if the error is an [APIError] with [Internal]
// status. Returns false if the error is not an [APIError] or has a different
// status.
func IsAPIErrorInternalStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == Internal
}

// IsAPIErrorInvalidArgumentStatus checks if the error is an [APIError] with
// [InvalidArgument] status. Returns false if the error is not an [APIError] or
// has a different status.
func IsAPIErrorInvalidArgumentStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == InvalidArgument
}

// IsAPIErrorNotFoundStatus checks if the error is an [APIError] with [NotFound]
// status. Returns false if the error is not an [APIError] or has a different
// status.
func IsAPIErrorNotFoundStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == NotFound
}

// IsAPIErrorPermissionDeniedStatus checks if the error is an [APIError] with
// [PermissionDenied] status. Returns false if the error is not an [APIError] or
// has a different status.
func IsAPIErrorPermissionDeniedStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == PermissionDenied
}

// IsAPIErrorQuotaExceededStatus checks if the error is an [APIError] with
// [QuotaExceeded] status. Returns false if the error is not an [APIError] or
// has a different status.
func IsAPIErrorQuotaExceededStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == QuotaExceeded
}

// IsAPIErrorIdempotencyKeyAlreadyUsedStatus checks if the error is an
// [APIError] with [IdempotencyKeyAlreadyUsed] status. Returns false if the
// error is not an [APIError] or has a different status.
func IsAPIErrorIdempotencyKeyAlreadyUsedStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == IdempotencyKeyAlreadyUsed
}

// IsAPIErrorInvalidEtagKeyStatus checks if the error is an APIError with
// InvalidEtagKey status. Returns false if the error is not an [APIError] or has
// a different status.
func IsAPIErrorInvalidEtagKeyStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == InvalidEtagKey
}

// IsAPIErrorUnauthenticatedStatus checks if the error is an [APIError] with
// [Unauthenticated] status. Returns false if the error is not an [APIError] or
// has a different status.
func IsAPIErrorUnauthenticatedStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == Unauthenticated
}

// IsAPIErrorUnavailableStatus checks if the error is an [APIError] with
// [Unavailable] status. Returns false if the error is not an [APIError] or has
// a different status.
func IsAPIErrorUnavailableStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == Unavailable
}

// IsAPIErrorMethodNotAllowedStatus checks if the error is an [APIError] with
// [MethodNotAllowed] status. Returns false if the error is not an [APIError] or
// has a different status.
func IsAPIErrorMethodNotAllowedStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == MethodNotAllowed
}

// IsAPIErrorTooManyRequestsStatus checks if the error is an [APIError] with
// [TooManyRequests] status. Returns false if the error is not an [APIError] or
// has a different status.
func IsAPIErrorTooManyRequestsStatus(err error) bool {
	var target *APIError
	if !errors.As(err, &target) {
		return false
	}

	return target.Status == TooManyRequests
}

const (
	httpStatusClientClosedRequest   = 499
	httpStatusNetworkConnectTimeout = 599
)

var (
	statusToStringMap = map[Status]string{
		Unknown:                   "UNKNOWN",
		AlreadyExists:             "ALREADY_EXISTS",
		Cancelled:                 "CANCELLED",
		DeadlineExceeded:          "DEADLINE_EXCEEDED",
		FailedPrecondition:        "FAILED_PRECONDITION",
		Internal:                  "INTERNAL",
		InvalidArgument:           "INVALID_ARGUMENT",
		NotFound:                  "NOT_FOUND",
		PermissionDenied:          "PERMISSION_DENIED",
		QuotaExceeded:             "QUOTA_EXCEEDED",
		IdempotencyKeyAlreadyUsed: "IDEMPOTENCY_KEY_ALREADY_USED",
		InvalidEtagKey:            "INVALID_ETAG_KEY",
		Unauthenticated:           "UNAUTHENTICATED",
		Unavailable:               "UNAVAILABLE",
		MethodNotAllowed:          "METHOD_NOT_ALLOWED",
		TooManyRequests:           "TOO_MANY_REQUESTS",
	}

	statusToHTTPCodesMap = map[Status][]int{
		Unknown:                   {http.StatusInternalServerError},
		AlreadyExists:             {http.StatusConflict},
		Cancelled:                 {httpStatusClientClosedRequest, httpStatusNetworkConnectTimeout},
		DeadlineExceeded:          {http.StatusRequestTimeout},
		FailedPrecondition:        {http.StatusPreconditionFailed},
		Internal:                  {http.StatusInternalServerError},
		InvalidArgument:           {http.StatusBadRequest},
		NotFound:                  {http.StatusNotFound},
		PermissionDenied:          {http.StatusForbidden},
		QuotaExceeded:             {http.StatusForbidden},
		IdempotencyKeyAlreadyUsed: {http.StatusUnprocessableEntity},
		InvalidEtagKey:            {http.StatusConflict},
		Unauthenticated:           {http.StatusUnauthorized},
		Unavailable:               {http.StatusServiceUnavailable},
		MethodNotAllowed:          {http.StatusMethodNotAllowed},
		TooManyRequests:           {http.StatusTooManyRequests},
	}
)

// Status represents an [APIError] status.
type Status uint8

// String returns the string representation of the status.
func (s Status) String() string {
	return statusToStringMap[s]
}

// HTTPCodes returns HTTP status codes associated with this error status.
// Some statuses may map to multiple HTTP codes depending on context.
func (s Status) HTTPCodes() []int {
	return statusToHTTPCodesMap[s]
}

const (
	Unknown Status = iota
	AlreadyExists
	Cancelled
	DeadlineExceeded
	FailedPrecondition
	Internal
	InvalidArgument
	NotFound
	PermissionDenied
	QuotaExceeded
	IdempotencyKeyAlreadyUsed
	InvalidEtagKey
	Unauthenticated
	Unavailable
	MethodNotAllowed
	TooManyRequests
)

func flattenJSON(data map[string]any, prefix string) string {
	var result strings.Builder

	keys := slices.Sorted(maps.Keys(data))
	for i, k := range keys {
		v := data[k]
		fullKey := prefix + k
		switch v := v.(type) {
		case map[string]any:
			result.WriteString(flattenJSON(v, fullKey+"."))
		default:
			fmt.Fprintf(&result, "%s - %v", fullKey, v)
		}
		if i < len(keys)-1 {
			result.WriteByte('\n')
		}
	}
	return result.String()
}
