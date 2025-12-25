package conv

import (
	"fmt"
)

type DecodeToNilError struct {
	name string
}

func NewDecodeToNilError(name string) error {
	return DecodeToNilError{name: name}
}

func (d DecodeToNilError) Error() string {
	return fmt.Sprintf("unable to decode %q to nil", d.name)
}

type DecodeToObjectError struct {
	name   string
	reason string
}

func NewDecodeToObjectError(name, reason string) error {
	return DecodeToObjectError{name: name, reason: reason}
}

func (d DecodeToObjectError) Error() string {
	return fmt.Sprintf("unable to decode %q to object: %s", d.name, d.reason)
}

type StringToTypeError struct {
	name string
	err  error
}

func NewStringToTypeError(name string, err error) error {
	return StringToTypeError{name: name, err: err}
}

func (d StringToTypeError) Error() string {
	return fmt.Sprintf("unable to convert %q string: %s", d.name, d.err)
}

func (d StringToTypeError) Unwrap() error {
	return d.err
}
