package largenumber

import (
	"math/big"
)

// Unit represents a unit of measurement for a large number.
type Unit uint8

const (
	// EmptyUnit represents no unit ("").
	EmptyUnit Unit = iota
	// K represents thousands ("K").
	K

	maxUnit
)

// IsValid reports whether unit is valid.
func (u Unit) IsValid() bool {
	return u < maxUnit
}

// GetValue returns the unit value. If unit is invalid, it returns zero value.
// Use [Unit.IsValid] to check if the unit is valid.
func (u Unit) GetValue() *big.Int {
	switch u {
	case EmptyUnit:
		return emptyUnit
	case K:
		return k
	default:
		return new(big.Int)
	}
}

// Decrease returns the next smaller unit and a boolean indicating if the new
// unit is valid.
func (u Unit) Decrease() (Unit, bool) {
	v := u - 1
	return v, v.IsValid()
}

// String returns the string representation of the unit. If unit is invalid, it
// returns an "invalid unit" string. Use [Unit.IsValid] to check if the unit is
// valid.
func (u Unit) String() string {
	switch u {
	case EmptyUnit:
		return ""
	case K:
		return "K"
	default:
		return "invalid unit"
	}
}

var (
	// EmptyUnit - minimum supported value
	emptyUnit = big.NewInt(1)
	// 1 K = 1_000 entities
	k = new(big.Int).Exp(big.NewInt(10), big.NewInt(3), nil)
)

func parseUnit(s string) (Unit, bool) {
	switch s {
	case "":
		return EmptyUnit, true
	case "k":
		return K, true
	default:
		return 0, false
	}
}
