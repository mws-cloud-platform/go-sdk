package throughput

import (
	"math/big"
)

// Unit represents a unit of measurement for throughput.
type Unit uint8

const (
	// BPS represents bytes per second ("Bps").
	BPS Unit = iota
	// KBPS represents kilobytes per second ("KBps").
	KBPS
	// MBPS represents megabytes per second ("MBps").
	MBPS
	// GBPS represents gigabytes per second ("GBps").
	GBPS
	// TBPS represents terabytes per second ("TBps").
	TBPS
	// PBPS represents petabytes per second ("PBps").
	PBPS
	// EBPS represents exabytes per second ("EBps").
	EBPS
	// ZBPS represents zettabytes per second ("ZBps").
	ZBPS
	// YBPS represents yottabytes per second ("YBps").
	YBPS
	// RBPS represents robobytes per second ("RBps").
	RBPS
	// QBPS represents quettabytes per second ("QBps").
	QBPS

	maxUnit
)

// IsValid reports whether unit is valid.
func (u Unit) IsValid() bool {
	return u < maxUnit
}

// GetValue returns the unit value in bytes per second. If unit is invalid, it
// returns zero value. Use [Unit.IsValid] to check if the unit is valid.
func (u Unit) GetValue() *big.Int {
	switch u {
	case BPS:
		return bps
	case KBPS:
		return kbps
	case MBPS:
		return mbps
	case GBPS:
		return gbps
	case TBPS:
		return tbps
	case PBPS:
		return pbps
	case EBPS:
		return ebps
	case ZBPS:
		return zbps
	case YBPS:
		return ybps
	case RBPS:
		return rbps
	case QBPS:
		return qbps
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
	case BPS:
		return "Bps"
	case KBPS:
		return "KBps"
	case MBPS:
		return "MBps"
	case GBPS:
		return "GBps"
	case TBPS:
		return "TBps"
	case PBPS:
		return "PBps"
	case EBPS:
		return "EBps"
	case ZBPS:
		return "ZBps"
	case YBPS:
		return "YBps"
	case RBPS:
		return "RBps"
	case QBPS:
		return "QBps"
	default:
		return "invalid unit"
	}
}

const lsh = 10

var (
	// Bps - minimum supported value
	bps = big.NewInt(1)
	// 1 KBps = 1_024 Bps
	kbps = new(big.Int).Lsh(bps, lsh)
	// 1 MBps = 1_024 KBps = 1_048_576 Bps
	mbps = new(big.Int).Lsh(kbps, lsh)
	// 1 GBps = 1_024 MBps = 1_048_576 KBps
	gbps = new(big.Int).Lsh(mbps, lsh)
	// 1 TBps = 1_024 GBps = 1_048_576 MBps
	tbps = new(big.Int).Lsh(gbps, lsh)
	// 1 PBps = 1_024 TBps = 1_048_576 GBps
	pbps = new(big.Int).Lsh(tbps, lsh)
	// 1 EBps = 1_024 PBps = 1_048_576 TBps
	ebps = new(big.Int).Lsh(pbps, lsh)
	// 1 ZBps = 1_024 EBps = 1_048_576 PBps
	zbps = new(big.Int).Lsh(ebps, lsh)
	// 1 YBps = 1_024 ZBps = 1_048_576 EBps
	ybps = new(big.Int).Lsh(zbps, lsh)
	// 1 RBps = 1_024 YBps = 1_048_576 ZBps
	rbps = new(big.Int).Lsh(ybps, lsh)
	// 1 QBps = 1_024 RBps = 1_048_576 YBps
	qbps = new(big.Int).Lsh(rbps, lsh)
)

func parseUnit(s string) (Unit, bool) {
	switch s {
	case "", "bps":
		return BPS, true
	case "kbps":
		return KBPS, true
	case "mbps":
		return MBPS, true
	case "gbps":
		return GBPS, true
	case "tbps":
		return TBPS, true
	case "pbps":
		return PBPS, true
	case "ebps":
		return EBPS, true
	case "zbps":
		return ZBPS, true
	case "ybps":
		return YBPS, true
	case "rbps":
		return RBPS, true
	case "qbps":
		return QBPS, true
	default:
		return 0, false
	}
}
