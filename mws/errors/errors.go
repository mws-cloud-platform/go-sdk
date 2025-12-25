// Package errors provides types and utilities for API errors.
package errors

import (
	"errors"
	"fmt"
)

// DecodeBodyError reports that client decoder was unable to decode the body.
type DecodeBodyError struct {
	// Content-Type from the response header.
	ContentType string
	// Underlying decoding error.
	Err error
}

func (d *DecodeBodyError) Error() string {
	if d.Err == nil {
		return fmt.Sprintf("decode content type '%s'", d.ContentType)
	}
	return fmt.Sprintf("decode content type '%s': %s", d.ContentType, d.Err)
}

// NewDecodeBodyError creates a new decode error.
func NewDecodeBodyError(contentType string, err error) *DecodeBodyError {
	return &DecodeBodyError{
		ContentType: contentType,
		Err:         err,
	}
}

// Unwrap returns the underlying error.
func (d *DecodeBodyError) Unwrap() error {
	return d.Err
}

// Is checks if the error is a [DecodeBodyError].
func (d *DecodeBodyError) Is(err error) bool {
	return IsDecodeBodyError(err)
}

// IsDecodeBodyError checks if an error is a [DecodeBodyError].
func IsDecodeBodyError(err error) bool {
	var target *DecodeBodyError
	return errors.As(err, &target)
}

// InvalidContentTypeError reports that client decoder got invalid content-type.
type InvalidContentTypeError struct {
	ContentType string
}

func (e *InvalidContentTypeError) Error() string {
	return fmt.Sprintf("invalid Content-Type: '%s'", e.ContentType)
}

// NewInvalidContentTypeError creates a new invalid content type error.
func NewInvalidContentTypeError(contentType string) *InvalidContentTypeError {
	return &InvalidContentTypeError{
		ContentType: contentType,
	}
}

// Is checks if the error is a [InvalidContentTypeError].
func (e *InvalidContentTypeError) Is(err error) bool {
	return IsInvalidContentTypeError(err)
}

// IsInvalidContentTypeError checks if an error is a [InvalidContentTypeError].
func IsInvalidContentTypeError(err error) bool {
	var target *InvalidContentTypeError
	return errors.As(err, &target)
}

// UnexpectedStatusCodeError reports that client got unexpected status code.
type UnexpectedStatusCodeError struct {
	StatusCode int
	Data       []byte
}

// NewUnexpectedStatusCodeError creates a new unexpected status code error.
func NewUnexpectedStatusCodeError(statusCode int) *UnexpectedStatusCodeError {
	return NewUnexpectedStatusCodeErrorWithData(statusCode, nil)
}

// NewUnexpectedStatusCodeErrorWithData creates a new unexpected status code error with data.
func NewUnexpectedStatusCodeErrorWithData(statusCode int, data []byte) *UnexpectedStatusCodeError {
	return &UnexpectedStatusCodeError{
		StatusCode: statusCode,
		Data:       data,
	}
}

func (e *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("unexpected status code: '%d'", e.StatusCode)
}

// Is checks if the error is a [UnexpectedStatusCodeError].
func (e *UnexpectedStatusCodeError) Is(err error) bool {
	return IsUnexpectedStatusCodeError(err)
}

// IsUnexpectedStatusCodeError checks if an error is a [UnexpectedStatusCodeError].
func IsUnexpectedStatusCodeError(err error) bool {
	var target *UnexpectedStatusCodeError
	return errors.As(err, &target)
}

// TransportError reports that client got transport error.
type TransportError struct {
	Err error
}

// NewTransportError creates a new transport error.
func NewTransportError(err error) *TransportError {
	return &TransportError{
		Err: err,
	}
}

func (t *TransportError) Error() string {
	if t.Err == nil {
		return "transport error"
	}

	return fmt.Sprintf("transport error: %s", t.Err)
}

// Unwrap returns the underlying error.
func (t *TransportError) Unwrap() error {
	return t.Err
}

// Is checks if the error is a [TransportError].
func (t *TransportError) Is(err error) bool {
	return IsTransportError(err)
}

// IsTransportError checks if an error is a [TransportError].
func IsTransportError(err error) bool {
	var target *TransportError
	return errors.As(err, &target)
}

// EncodeBodyError reports that client got encode body error.
type EncodeBodyError struct {
	Err error
}

// NewEncodeBodyError creates a new encode body error.
func NewEncodeBodyError(err error) *EncodeBodyError {
	return &EncodeBodyError{
		Err: err,
	}
}

func (e *EncodeBodyError) Error() string {
	if e.Err == nil {
		return "encode error"
	}

	return fmt.Sprintf("encode error: %s", e.Err)
}

// Unwrap returns the underlying error.
func (e *EncodeBodyError) Unwrap() error {
	return e.Err
}

// Is checks if the error is a [EncodeBodyError].
func (e *EncodeBodyError) Is(err error) bool {
	return IsEncodeBodyError(err)
}

// IsEncodeBodyError checks if an error is a [EncodeBodyError].
func IsEncodeBodyError(err error) bool {
	var target *EncodeBodyError
	return errors.As(err, &target)
}
