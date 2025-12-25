package errors

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mws.cloud/go-sdk/mws/errors"
	"go.mws.cloud/go-sdk/pkg/backoff"
)

type APIGenError interface {
	GetCodeOr(string) string
	GetDescriptionOr(string) string
}

type ExtAPIGenError interface {
	APIGenError
	GetCommonRetryPolicy() *errors.RetryPolicy
	GetType() *string
	GetDetails() map[string]json.RawMessage
}

func WrapAPIGenError(code int, err APIGenError) *errors.APIError {
	result := &errors.APIError{
		Code:        code,
		Description: err.GetDescriptionOr(""),
	}

	rawStrStatus := err.GetCodeOr("")

	status, ok := stringToStatusMap[rawStrStatus]
	if !ok {
		status = errors.Unknown
		if result.Description != "" {
			result.Description += "\n"
		}
		result.Description += fmt.Sprintf("Invalid status code: %q", rawStrStatus)
	}
	result.Status = status

	if extErr, ok := err.(ExtAPIGenError); ok {
		result.RetryPolicy = extErr.GetCommonRetryPolicy()
		rawDetails := extErr.GetDetails()
		if len(rawDetails) != 0 {
			data := make(map[string]any, len(rawDetails))
			for k, v := range rawDetails {
				var fieldValue any
				if errUnmarshal := json.Unmarshal(v, &fieldValue); errUnmarshal != nil {
					data[k] = map[string]json.RawMessage{"raw_json": v}
					continue
				}
				data[k] = fieldValue
			}
			result.Details = data
		}
	}

	return result
}

func WrapAPIGenDefaultError(code int, apiGenDefaultErr any) *errors.APIError {
	var description string

	switch e := apiGenDefaultErr.(type) {
	case APIGenDefaultError:
		description = e.GetErrorOr("")
	case error:
		description = e.Error()
	case fmt.Stringer:
		description = e.String()
	default:
		description = "unknown error"
	}
	return errors.NewAPIError(code, errors.Unknown, description)
}

type APIGenDefaultError interface {
	GetErrorOr(string) string
}

var (
	stringToStatusMap = map[string]errors.Status{
		"ALREADY_EXISTS":               errors.AlreadyExists,
		"CANCELLED":                    errors.Cancelled,
		"DEADLINE_EXCEEDED":            errors.DeadlineExceeded,
		"FAILED_PRECONDITION":          errors.FailedPrecondition,
		"INTERNAL":                     errors.Internal,
		"INVALID_ARGUMENT":             errors.InvalidArgument,
		"NOT_FOUND":                    errors.NotFound,
		"PERMISSION_DENIED":            errors.PermissionDenied,
		"QUOTA_EXCEEDED":               errors.QuotaExceeded,
		"IDEMPOTENCY_KEY_ALREADY_USED": errors.IdempotencyKeyAlreadyUsed,
		"INVALID_ETAG_KEY":             errors.InvalidEtagKey,
		"UNAUTHENTICATED":              errors.Unauthenticated,
		"UNAVAILABLE":                  errors.Unavailable,
		"METHOD_NOT_ALLOWED":           errors.MethodNotAllowed,
		"TOO_MANY_REQUESTS":            errors.TooManyRequests,
	}
)

const (
	defaultRetryTimeout    = 20 * time.Millisecond
	defaultMaxRetryTimeout = 10 * time.Second
	defaultRetryCount      = 3
)

func NewDefaultRetryPolicy() *errors.RetryPolicy {
	return &errors.RetryPolicy{
		MaxRetryTimeout:   defaultMaxRetryTimeout,
		RetryCount:        defaultRetryCount,
		RetryTimeout:      defaultRetryTimeout,
		RetryTimeoutScale: backoff.ScaleExponential,
	}
}
