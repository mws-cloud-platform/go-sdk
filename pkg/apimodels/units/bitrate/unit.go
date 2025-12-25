package bitrate

import (
	"math/big"
)

// Unit represents a unit of measurement for a bitrate.
type Unit uint8

const (
	// BITS represents bits per second ("bit/s").
	BITS Unit = iota
	// KBITS represents kilobits per second ("kbit/s").
	KBITS
	// MBITS represents megabits per second ("Mbit/s").
	MBITS
	// GBITS represents gigabits per second ("Gbit/s").
	GBITS
	// TBITS represents terabits per second ("Tbit/s").
	TBITS

	maxUnit
)

// IsValid reports whether unit is valid.
func (u Unit) IsValid() bool {
	return u < maxUnit
}

// GetValue returns unit value in bits per second. If unit is invalid, it
// returns zero value. Use [Unit.IsValid] to check if the unit is valid.
func (u Unit) GetValue() *big.Int {
	switch u {
	case BITS:
		return bits
	case KBITS:
		return kbits
	case MBITS:
		return mbits
	case GBITS:
		return gbits
	case TBITS:
		return tbits
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
	case BITS:
		return "bit/s"
	case KBITS:
		return "kbit/s"
	case MBITS:
		return "Mbit/s"
	case GBITS:
		return "Gbit/s"
	case TBITS:
		return "Tbit/s"
	default:
		return "invalid unit"
	}
}

func parseUnit(s string) (Unit, bool) {
	switch s {
	case "", "bit/s":
		return BITS, true
	case "kbit/s":
		return KBITS, true
	case "mbit/s":
		return MBITS, true
	case "gbit/s":
		return GBITS, true
	case "tbit/s":
		return TBITS, true
	default:
		return 0, false
	}
}

const lsh = 10

var (
	// bit/s - minimum supported value
	bits = big.NewInt(1)
	// 1 Kbit/s = 1_024 bit/s
	kbits = new(big.Int).Lsh(bits, lsh)
	// 1 Mbit/s = 1_024 Kbit/s = 1_048_576 bit/s
	mbits = new(big.Int).Lsh(kbits, lsh)
	// 1 Gbit/s = 1_024 Mbit/s = 1_048_576 Kbit/s
	gbits = new(big.Int).Lsh(mbits, lsh)
	// 1 Tbit/s = 1_024 Gbit/s = 1_048_576 Mbit/s
	tbits = new(big.Int).Lsh(gbits, lsh)
)
