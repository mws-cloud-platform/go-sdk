package client

import (
	"go.mws.cloud/go-sdk/pkg/optional"
)

// Deprecated: use NewOptionalNil from go.mws.cloud/go-sdk/go/pkg/optional
func NewOptionalNil[T any](val T) OptionalNil[T] {
	return NewDirectOptionalNil(val, true, false)
}

// Deprecated: use NewOptionalNil or fill OptionalNil struct
// from go.mws.cloud/go-sdk/go/pkg/optional
func NewDirectOptionalNil[T any](val T, set, isNil bool) OptionalNil[T] {
	return OptionalNil[T]{
		Value: val,
		Set:   set,
		Null:  isNil,
	}
}

// Deprecated: use "go.mws.cloud/go-sdk/go/pkg/optional"
type OptionalNil[T any] = optional.OptionalNil[T]

// Deprecated: use NewOptional from go.mws.cloud/go-sdk/go/pkg/optional
func NewOptional[T any](val T) Optional[T] {
	return NewDirectOptional(val, true)
}

// Deprecated: use NewOptionalNil or fill OptionalNil struct
// from go.mws.cloud/go-sdk/go/pkg/optional
func NewDirectOptional[T any](val T, set bool) Optional[T] {
	return Optional[T]{
		Value: val,
		Set:   set,
	}
}

// Deprecated: use "go.mws.cloud/go-sdk/go/pkg/optional"
type Optional[T any] = optional.Optional[T]
