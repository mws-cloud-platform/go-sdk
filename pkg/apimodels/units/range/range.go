// Package rangeunit provides types and utilities for working with ranges.
package rangeunit

import (
	"fmt"
	"regexp"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/must"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	uniterrors "go.mws.cloud/go-sdk/pkg/apimodels/units/errors"
)

type ValueUnit[T any] interface {
	Cmp(T) int
	CloneByValue() T
	Equal(other T) bool
	ParseString(s string) (T, error)
	QuantityString() string
	UnitString() string
}

// Range wrapper that supports parsing of Range values.
type Range[T ValueUnit[T]] struct {
	minValue T
	maxValue T
	rawValue *string
}

func NewRange[T ValueUnit[T]](minValue, maxValue T) (Range[T], error) {
	if minValue.Cmp(maxValue) > 0 {
		return Range[T]{}, ErrMinValueGreaterThanMaxValue
	}
	if minValue.UnitString() != maxValue.UnitString() {
		return Range[T]{}, ErrDifferentUnitValues
	}
	return Range[T]{
		minValue: minValue,
		maxValue: maxValue,
	}, nil
}

// MaxValue returns the maximum value in the range.
func (r Range[T]) MaxValue() T {
	return r.maxValue
}

// MinValue returns the minimum value in the range.
func (r Range[T]) MinValue() T {
	return r.minValue
}

// RawValue returns a raw value if it was created from a string.
func (r Range[T]) RawValue() *string {
	return r.rawValue
}

// Clone returns a clone Range pointer with new pointer values.
func (r *Range[T]) Clone() *Range[T] {
	if r == nil {
		return nil
	}

	clone := *r
	clone.minValue = r.minValue.CloneByValue()
	clone.maxValue = r.maxValue.CloneByValue()
	clone.rawValue = ptr.Clone(r.rawValue)

	return &clone
}

// Equal checks if the values of r and r2 are equal.
// This method is a deep comparison of models.
func (r Range[T]) Equal(r2 Range[T]) bool {
	return r.minValue.Equal(r2.minValue) && r.maxValue.Equal(r2.maxValue) && ptr.Equal(r.rawValue, r2.rawValue)
}

// String returns a string representing the range in the format of "10 - 20<unit>".
func (r Range[T]) String() string {
	return fmt.Sprintf("%s - %s%s", r.minValue.QuantityString(), r.maxValue.QuantityString(), r.minValue.UnitString())
}

func (r Range[T]) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	r.Encode(&e)
	return e.Bytes(), nil
}

func (r *Range[T]) Encode(e *jx.Encoder) {
	switch {
	case r == nil:
		e.Null()
	case r.rawValue != nil:
		e.Str(*r.rawValue)
	default:
		e.Str(r.String())
	}
}

func (r *Range[T]) UnmarshalJSON(bytes []byte) error {
	return r.Decode(jx.DecodeBytes(bytes))
}

func (r *Range[T]) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseString[T](rawValue)
	if err != nil {
		return err
	}

	r.minValue = parsed.minValue
	r.maxValue = parsed.maxValue
	r.rawValue = parsed.rawValue
	return nil
}

const (
	numSubMatchesInRegexp = 4
	maxStrLen             = 128
)

var validationRegexp = regexp.MustCompile(`(?i)^(\s*-?[0-9]+(?:\.[0-9]+)?)\s*-\s*(-?[0-9]+(?:\.[0-9]+)?)\s*(\S*)\s*$`)

// ParseString parses a range string.
// A range string is a string consisting of a
// decimal number and possibly a unit of measurement.
// The range string is a string consisting two decimal numbers
// through the dash and their unit of measurement.
// Examples of string range are "10-20 bit/s", "10 - 20 GB", "1.5-2 MB".
func ParseString[T ValueUnit[T]](s string) (Range[T], error) {
	if len(s) > maxStrLen {
		return Range[T]{}, uniterrors.ErrStringTooLong
	}

	matches := validationRegexp.FindStringSubmatch(s)
	if len(matches) != numSubMatchesInRegexp {
		return Range[T]{}, uniterrors.ErrStringDoesNotMatchRegexp
	}

	var v T

	minValue, err := v.ParseString(matches[1] + matches[3])
	if err != nil {
		return Range[T]{}, fmt.Errorf("%w: %w", ErrIncorrectMinValue, err)
	}

	maxValue, err := v.ParseString(matches[2] + matches[3])
	if err != nil {
		return Range[T]{}, fmt.Errorf("%w: %w", ErrIncorrectMaxValue, err)
	}

	if minValue.Cmp(maxValue) > 0 {
		return Range[T]{}, fmt.Errorf("%w: %w", uniterrors.ErrInvalidString, ErrMinValueGreaterThanMaxValue)
	}

	return Range[T]{
		minValue: minValue,
		maxValue: maxValue,
		rawValue: &s,
	}, nil
}

// MustParseString is like [ParseString] but panics on error.
func MustParseString[T ValueUnit[T]](s string) Range[T] {
	return must.Value(ParseString[T](s))
}
