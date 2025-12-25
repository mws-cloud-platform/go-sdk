package frequency

import (
	"math/big"
)

// Unit represents a unit of measurement for a frequency.
type Unit uint8

const (
	// HZ represents hertz ("Hz").
	HZ Unit = iota
	// KHZ represents kilohertz ("kHz").
	KHZ
	// MHZ represents megahertz ("MHz").
	MHZ
	// GHZ represents gigahertz ("GHz").
	GHZ
	// THZ represents terahertz ("THz").
	THZ

	maxUnit
)

// IsValid reports whether unit is valid.
func (u Unit) IsValid() bool {
	return u < maxUnit
}

// GetValue returns the unit value in hertz. If unit is invalid, it returns zero
// value. Use [Unit.IsValid] to check if the unit is valid.
func (u Unit) GetValue() *big.Int {
	switch u {
	case HZ:
		return hz
	case KHZ:
		return khz
	case MHZ:
		return mhz
	case GHZ:
		return ghz
	case THZ:
		return thz
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
	case HZ:
		return "Hz"
	case KHZ:
		return "kHz"
	case MHZ:
		return "MHz"
	case GHZ:
		return "GHz"
	case THZ:
		return "THz"
	default:
		return "invalid unit"
	}
}

var (
	ten = big.NewInt(10)

	// Hz - minimum supported value
	hz = big.NewInt(1)
	// 1 kHz = 1_000 Hz
	khz = new(big.Int).Exp(ten, big.NewInt(3), nil)
	// 1 MHz = 1_000 kHz = 1_000_000 Hz
	mhz = new(big.Int).Exp(ten, big.NewInt(6), nil)
	// 1 GHz = 1_000 MHz = 1_000_000 kHz = 1_000_000_000 Hz
	ghz = new(big.Int).Exp(ten, big.NewInt(9), nil)
	// 1 THz = 1_000 GHz = 1_000_000 MHz = 1_000_000_000 kHz = 1_000_000_000_000 Hz
	thz = new(big.Int).Exp(ten, big.NewInt(12), nil)
)

func parseUnit(s string) (Unit, bool) {
	switch s {
	case "", "hz":
		return HZ, true
	case "khz":
		return KHZ, true
	case "mhz":
		return MHZ, true
	case "ghz":
		return GHZ, true
	case "thz":
		return THZ, true
	default:
		return 0, false
	}
}
