package bytesize

import (
	"math/big"
)

// Unit represents a unit of measurement for a byte size.
type Unit uint8

const (
	// B represents bytes ("B").
	B Unit = iota
	// KB represents kilobytes ("KB").
	KB
	// MB represents megabytes ("MB").
	MB
	// GB represents gigabytes ("GB").
	GB
	// TB represents terabytes ("TB").
	TB
	// PB represents petabytes ("PB").
	PB
	// EB represents exabytes ("EB").
	EB
	// ZB represents zettabytes ("ZB").
	ZB
	// YB represents yottabytes ("YB").
	YB
	// RB represents robobytes ("RB").
	RB
	// QB represents quettabytes ("QB").
	QB

	maxUnit
)

// IsValid reports whether unit is valid.
func (u Unit) IsValid() bool {
	return u < maxUnit
}

// GetValue returns the unit value in bytes. If unit is invalid, it returns zero
// value. Use [Unit.IsValid] to check if the unit is valid.
func (u Unit) GetValue() *big.Int {
	switch u {
	case B:
		return b
	case KB:
		return kb
	case MB:
		return mb
	case GB:
		return gb
	case TB:
		return tb
	case PB:
		return pb
	case EB:
		return eb
	case ZB:
		return zb
	case YB:
		return yb
	case RB:
		return rb
	case QB:
		return qb
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
	case B:
		return "B"
	case KB:
		return "KB"
	case MB:
		return "MB"
	case GB:
		return "GB"
	case TB:
		return "TB"
	case PB:
		return "PB"
	case EB:
		return "EB"
	case ZB:
		return "ZB"
	case YB:
		return "YB"
	case RB:
		return "RB"
	case QB:
		return "QB"
	default:
		return "invalid unit"
	}
}

const lsh = 10

var (
	// B - minimum supported value
	b = big.NewInt(1)
	// 1 KB = 1_024 B
	kb = new(big.Int).Lsh(b, lsh)
	// 1 MB = 1_024 KB = 1_048_576 B
	mb = new(big.Int).Lsh(kb, lsh)
	// 1 GB = 1_024 MB = 1_048_576 KB
	gb = new(big.Int).Lsh(mb, lsh)
	// 1 TB = 1_024 GB = 1_048_576 MB
	tb = new(big.Int).Lsh(gb, lsh)
	// 1 PB = 1_024 TB = 1_048_576 GB
	pb = new(big.Int).Lsh(tb, lsh)
	// 1 EB = 1_024 PB = 1_048_576 TB
	eb = new(big.Int).Lsh(pb, lsh)
	// 1 ZB = 1_024 EB = 1_048_576 PB
	zb = new(big.Int).Lsh(eb, lsh)
	// 1 YB = 1_024 ZB = 1_048_576 EB
	yb = new(big.Int).Lsh(zb, lsh)
	// 1 RB = 1_024 YB = 1_048_576 ZB
	rb = new(big.Int).Lsh(yb, lsh)
	// 1 QB = 1_024 RB = 1_048_576 YB
	qb = new(big.Int).Lsh(rb, lsh)
)

func parseUnit(s string) (Unit, bool) {
	switch s {
	case "", "b":
		return B, true
	case "kb":
		return KB, true
	case "mb":
		return MB, true
	case "gb":
		return GB, true
	case "tb":
		return TB, true
	case "pb":
		return PB, true
	case "eb":
		return EB, true
	case "zb":
		return ZB, true
	case "yb":
		return YB, true
	case "rb":
		return RB, true
	case "qb":
		return QB, true
	default:
		return 0, false
	}
}
