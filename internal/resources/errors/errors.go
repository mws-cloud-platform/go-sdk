package errors

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-faster/jx"
)

type ParseReferenceError struct {
	RawRef string
	Err    error
}

func (r *ParseReferenceError) Error() string {
	if r.Err == nil {
		return fmt.Sprintf("parse reference '%s'", r.RawRef)
	}
	return fmt.Sprintf("parse reference '%s': %s", r.RawRef, r.Err)
}

func (r *ParseReferenceError) Unwrap() error {
	return r.Err
}

func (r *ParseReferenceError) Is(err error) bool {
	return IsParseReferenceError(err)
}

func IsParseReferenceError(err error) bool {
	var target *ParseReferenceError
	return errors.As(err, &target)
}

func NewParseReferenceError(rawRef string, err error) *ParseReferenceError {
	return &ParseReferenceError{
		RawRef: rawRef,
		Err:    err,
	}
}

type ParseIDError struct {
	RawID string
	Err   error
}

func (r *ParseIDError) Error() string {
	if r.Err == nil {
		return fmt.Sprintf("parse id '%s'", r.RawID)
	}
	return fmt.Sprintf("parse id '%s': %s", r.RawID, r.Err)
}

func (r *ParseIDError) Unwrap() error {
	return r.Err
}

func (r *ParseIDError) Is(err error) bool {
	return IsParseIDError(err)
}

func IsParseIDError(err error) bool {
	var target *ParseIDError
	return errors.As(err, &target)
}

func NewParseIDError(rawID string, err error) *ParseIDError {
	return &ParseIDError{
		RawID: rawID,
		Err:   err,
	}
}

type InitOneOfReferenceError struct {
	BrokenReference any
}

func (o *InitOneOfReferenceError) Error() string {
	return fmt.Sprintf("type '%T' is not supported for init a one of reference", o.BrokenReference)
}

func (o *InitOneOfReferenceError) Is(err error) bool {
	return IsInitOneOfReferenceError(err)
}

func IsInitOneOfReferenceError(err error) bool {
	var target *InitOneOfReferenceError
	return errors.As(err, &target)
}

func NewInitOneOfReferenceError(brokenReference any) *InitOneOfReferenceError {
	return &InitOneOfReferenceError{
		BrokenReference: brokenReference,
	}
}

type InitOneOfIDError struct {
	BrokenID any
}

func (o *InitOneOfIDError) Error() string {
	return fmt.Sprintf("type '%T' is not supported for init a one of id", o.BrokenID)
}

func (o *InitOneOfIDError) Is(err error) bool {
	return IsInitOneOfIDError(err)
}

func IsInitOneOfIDError(err error) bool {
	var target *InitOneOfIDError
	return errors.As(err, &target)
}

func NewInitOneOfIDError(brokenID any) *InitOneOfIDError {
	return &InitOneOfIDError{
		BrokenID: brokenID,
	}
}

type PathAccumulatorError struct {
	Path string
	Err  error
}

func (r *PathAccumulatorError) Error() string {
	if r.Err == nil {
		return fmt.Sprintf("path '%s'", r.Path)
	}
	return fmt.Sprintf("path '%s': %s", r.Path, r.Err)
}

func (r *PathAccumulatorError) Unwrap() error {
	return r.Err
}

func (r *PathAccumulatorError) Is(err error) bool {
	return IsPathAccumulatorError(err)
}

func IsPathAccumulatorError(err error) bool {
	var target *PathAccumulatorError
	return errors.As(err, &target)
}

func NewPathAccumulatorError(path string, err error) *PathAccumulatorError {
	var target *PathAccumulatorError
	if errors.As(err, &target) {
		if len(target.Path) != 0 && target.Path[0] != '[' {
			path += "."
		}
		path += target.Path
		err = target.Err
	}
	return &PathAccumulatorError{
		Path: path,
		Err:  err,
	}
}

func NewPathAccumulatorErrorAsIndex[T int | string](idx T, err error) *PathAccumulatorError {
	return NewPathAccumulatorError(fmt.Sprintf(`[%v]`, idx), err)
}

func PathAccumulatorErrorObjBytesFuncWrap(f func(*jx.Decoder, []byte) error) func(*jx.Decoder, []byte) error {
	return func(d *jx.Decoder, key []byte) error {
		if err := f(d, key); err != nil {
			return NewPathAccumulatorError(string(key), err)
		}
		return nil
	}
}

func PathAccumulatorErrorAsIndexObjBytesFuncWrap(f func(*jx.Decoder, []byte) error) func(*jx.Decoder, []byte) error {
	return func(d *jx.Decoder, key []byte) error {
		if err := f(d, key); err != nil {
			return NewPathAccumulatorErrorAsIndex(strconv.Quote(string(key)), err)
		}
		return nil
	}
}

func PathAccumulatorErrorAsIndexArrFuncWrap(f func(*jx.Decoder) error) func(*jx.Decoder) error {
	idx := 0
	return func(d *jx.Decoder) error {
		if err := f(d); err != nil {
			return NewPathAccumulatorErrorAsIndex(idx, err)
		}
		idx++
		return nil
	}
}
