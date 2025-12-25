// Package errors provides common errors related to units.
package errors

import (
	"fmt"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

const (
	ErrInvalidString   = consterr.Error("invalid string")
	ErrInvalidQuantity = consterr.Error("invalid quantity")
	ErrInvalidUnit     = consterr.Error("invalid unit")
)

var (
	ErrUnexpectedUnit                 = fmt.Errorf("%w: unexpected unit", ErrInvalidUnit)
	ErrStringTooLong                  = fmt.Errorf("%w: string too long", ErrInvalidString)
	ErrStringDoesNotMatchRegexp       = fmt.Errorf("%w: string does not match regexp", ErrInvalidString)
	ErrNegativeQuantity               = fmt.Errorf("%w: negative quantity", ErrInvalidQuantity)
	ErrZeroWithFractionalPart         = fmt.Errorf("%w: quantity is a zero with fractional part", ErrInvalidQuantity)
	ErrBaseQuantityHaveFractionalPart = fmt.Errorf("%w: quantity after conversion to the base quantity have a fractional part", ErrInvalidQuantity)
)
