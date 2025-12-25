package client

import (
	"encoding/json"

	"github.com/go-faster/jx"
)

func NewOptionalNil[T any](val T) OptionalNil[T] {
	return NewDirectOptionalNil(val, true, false)
}

func NewDirectOptionalNil[T any](val T, set, isNil bool) OptionalNil[T] {
	return OptionalNil[T]{
		Value: val,
		Set:   set,
		Null:  isNil,
	}
}

type OptionalNil[T any] struct {
	Value T
	Set   bool
	Null  bool
}

// IsSet returns true if optional value was set.
func (o *OptionalNil[T]) IsSet() bool {
	return o.Set
}

// IsNull returns true if value is Null
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

// Get returns value and boolean that denotes whether value was set.
func (o OptionalNil[T]) Get() (v T, ok bool) { //nolint:nonamedreturns // generic
	if o.Null {
		return v, false
	}
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
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

func NewOptional[T any](val T) Optional[T] {
	return NewDirectOptional(val, true)
}

func NewDirectOptional[T any](val T, set bool) Optional[T] {
	return Optional[T]{
		Value: val,
		Set:   set,
	}
}

// Optional is optional value.
type Optional[T any] struct {
	Value T
	Set   bool
}

// IsSet returns true if optional value was set.
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

// Get returns value and boolean that denotes whether value was set.
func (o Optional[T]) Get() (v T, ok bool) { //nolint:nonamedreturns // generic
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
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
