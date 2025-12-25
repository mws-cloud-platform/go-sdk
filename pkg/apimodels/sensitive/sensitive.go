// Package sensitive provides types and utilities for working with sensitive
// data such as passwords, tokens, or other confidential information.
package sensitive

import (
	"encoding/json"
	"log/slog"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

const defaultFormat = "****"

// Option is a sensitive value option.
type Option[T any] func(*Sensitive[T])

// WithFormat is an option that sets the custom format function for the
// sensitive value.
func WithFormat[T any](format func(T) string) Option[T] {
	return func(s *Sensitive[T]) {
		s.format = format
	}
}

// Sensitive is a sensitive value container that can be used to store sensitive
// data such as passwords, tokens, or other confidential information. The value
// will be hidden when formatting (string, text, json or yaml) or logging (slog
// and zap). To get the raw value, use the [Sensitive.Value] method, but it is
// recommended to call this method only where the value is really used - in
// other cases, pass the data in this container.
type Sensitive[T any] struct {
	value  T
	format func(T) string
}

// New creates a new sensitive value container with the given value and options.
func New[T any](value T, opts ...Option[T]) Sensitive[T] {
	v := Sensitive[T]{value: value}
	for _, opt := range opts {
		opt(&v)
	}
	return v
}

// Value returns the raw value of the sensitive container.
func (s Sensitive[T]) Value() T {
	return s.value
}

// String returns the redacted value of the sensitive container.
func (s Sensitive[T]) String() string {
	return s.redacted()
}

func (s Sensitive[T]) MarshalText() ([]byte, error) {
	return []byte(s.redacted()), nil
}

func (s Sensitive[T]) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.redacted() + `"`), nil
}

func (s *Sensitive[T]) UnmarshalJSON(data []byte) error {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	s.value = value
	return nil
}

func (s Sensitive[T]) MarshalYAML() (any, error) {
	return s.redacted(), nil
}

func (s *Sensitive[T]) UnmarshalYAML(node *yaml.Node) error {
	var value T
	if err := node.Decode(&value); err != nil {
		return err
	}
	s.value = value
	return nil
}

func (s Sensitive[T]) LogValue() slog.Value {
	return slog.StringValue(s.redacted())
}

func (s Sensitive[T]) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("value", s.redacted())
	return nil
}

func (s Sensitive[T]) redacted() string {
	if s.format != nil {
		return s.format(s.value)
	}
	return defaultFormat
}
