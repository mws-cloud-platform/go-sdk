// Package optional provides generic container types for optional values.
package optional

import (
	"encoding/json"

	"github.com/go-faster/jx"
)

// OptionalNil represents an optional value of type T that can be in one of three valid states:
//   - not set (Set = false, Null = false)
//   - set to a value (Set = true, Null = false)
//   - explicitly set to null (Set = true, Null = true)
//
// Fourth possible state is considered invalid (Set = false, Null = true)
//
// This is used when a field may be omitted, present as a value,
// or present as null. The zero value is "not set".
type OptionalNil[T any] struct {
	Value T
	Set   bool
	Null  bool
}

// NewOptionalNil returns an OptionalNil[T] in the "set to value" state,
// with the provided val and Set = true, Null = false.
func NewOptionalNil[T any](val T) OptionalNil[T] {
	return OptionalNil[T]{
		Value: val,
		Set:   true,
		Null:  false,
	}
}

// IsSet returns true if optional value is set.
func (o *OptionalNil[T]) IsSet() bool {
	return o.Set
}

// IsNull returns true if the optional is set to null.
// It returns false when the optional is not set or when it is set to a value.
func (o *OptionalNil[T]) IsNull() bool {
	return o.Null
}

// Reset unsets value.
func (o *OptionalNil[T]) Reset() {
	var v T
	o.Value = v
	o.Set = false
	o.Null = false
}

// SetTo sets value to v.
func (o *OptionalNil[T]) SetTo(v T) {
	o.Set = true
	o.Null = false
	o.Value = v
}

// SetToNull sets value to null
func (o *OptionalNil[T]) SetToNull() {
	o.Set = true
	o.Null = true
	var v T
	o.Value = v
}

// Get returns value and a boolean indicating whether the value is set.
// It returns (zero value, false) if the optional is not set or null.
func (o OptionalNil[T]) Get() (v T, ok bool) {
	if o.Null {
		return v, false
	}
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set to a non-null value, otherwise returns the provided default d.
func (o OptionalNil[T]) Or(d T) T {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

func (o *OptionalNil[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.SetToNull()
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.SetTo(value)
	return nil
}

func (o *OptionalNil[T]) Decode(d *jx.Decoder) error {
	data, err := d.Raw()
	if err != nil {
		return err
	}

	return o.UnmarshalJSON(data)
}

// Optional represents an optional value of type T that can be in one of two states:
//   - not set (Set = false)
//   - set to a value (Set = true)
//
// The zero value is "not set". To distinguish between null and a missing field,
// use OptionalNil[T] instead.
type Optional[T any] struct {
	Value T
	Set   bool
}

// NewOptional returns an Optional[T] in the "set" state with the provided value.
func NewOptional[T any](val T) Optional[T] {
	return Optional[T]{
		Value: val,
		Set:   true,
	}
}

// IsSet returns true if optional value is set.
func (o Optional[T]) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *Optional[T]) Reset() {
	var v T
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *Optional[T]) SetTo(v T) {
	o.Set = true
	o.Value = v
}

// Get returns value and a boolean indicating whether the value is set.
// It returns (zero value, false) if the optional is not set.
func (o Optional[T]) Get() (v T, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, otherwise returns the provided default d.
func (o Optional[T]) Or(d T) T {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.SetTo(value)
	return nil
}

func (o *Optional[T]) Decode(d *jx.Decoder) error {
	data, err := d.Raw()
	if err != nil {
		return err
	}

	return o.UnmarshalJSON(data)
}
