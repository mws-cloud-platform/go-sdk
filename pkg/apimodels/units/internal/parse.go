package internal

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"

	uniterrors "go.mws.cloud/go-sdk/pkg/apimodels/units/errors"
)

const maxStrLen = 128

// Unit represents a unit of measurement.
type Unit interface {
	GetValue() *big.Int
}

// ParseString parses unit value from string in format "<quantity> <unit>". The
// resulting base quantity must be an integer. Spacing between quantity and unit
// are allowed and ignored. The same for leading and trailing whitespace.
func ParseString[U Unit](s string, re *regexp.Regexp, parseUnit func(string) (U, bool)) (*big.Int, U, error) {
	var unit U

	if len(s) > maxStrLen {
		return nil, unit, uniterrors.ErrStringTooLong
	}

	matches := re.FindStringSubmatch(strings.ToLower(strings.TrimSpace(s)))
	if len(matches) != 3 {
		return nil, unit, uniterrors.ErrInvalidString
	}

	quantity, err := decimal.NewFromString(matches[1])
	if err != nil {
		return nil, unit, fmt.Errorf("%w: %w", uniterrors.ErrInvalidQuantity, err)
	}
	u := matches[2]

	var ok bool
	unit, ok = parseUnit(strings.ToLower(u))
	if !ok {
		return nil, unit, fmt.Errorf("%w %q", uniterrors.ErrUnexpectedUnit, u)
	}

	if !quantity.IsZero() && quantity.IntPart() == 0 && !quantity.IsInteger() {
		return nil, unit, uniterrors.ErrZeroWithFractionalPart
	}

	quantity = quantity.Mul(decimal.NewFromBigInt(unit.GetValue(), 0))
	if !quantity.IsInteger() {
		return nil, unit, uniterrors.ErrBaseQuantityHaveFractionalPart
	}

	return quantity.BigInt(), unit, nil
}
