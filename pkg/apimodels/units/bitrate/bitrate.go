// Package bitrate provides functionality for working with bitrates.
package bitrate

import (
	"math/big"
	"regexp"

	"github.com/go-faster/jx"
	"github.com/shopspring/decimal"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.mws.cloud/util-toolset/pkg/utils/must"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	uniterrors "go.mws.cloud/go-sdk/pkg/apimodels/units/errors"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/internal"
)

var bitrateRegexp = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)\s*([kmgt]?bit/s)?$`)

// Bitrate represents a bitrate value.
type Bitrate struct {
	baseQuantity *big.Int
	unit         Unit
	raw          *string
}

// NewFromInt64 creates a new bitrate value from int64.
func NewFromInt64(quantity int64, unit Unit) (Bitrate, error) {
	if quantity < 0 {
		return Bitrate{}, uniterrors.ErrNegativeQuantity
	}
	if !unit.IsValid() {
		return Bitrate{}, uniterrors.ErrUnexpectedUnit
	}

	return Bitrate{
		baseQuantity: new(big.Int).Mul(big.NewInt(quantity), unit.GetValue()),
		unit:         unit,
	}, nil
}

// NewFromBigInt creates a new bitrate value from [big.Int]. If quantity is nil,
// it is set to zero.
func NewFromBigInt(quantity *big.Int, unit Unit) (Bitrate, error) {
	if quantity == nil {
		quantity = new(big.Int)
	}
	if quantity.Sign() == -1 {
		return Bitrate{}, uniterrors.ErrNegativeQuantity
	}
	if !unit.IsValid() {
		return Bitrate{}, uniterrors.ErrUnexpectedUnit
	}

	return Bitrate{
		baseQuantity: new(big.Int).Mul(quantity, unit.GetValue()),
		unit:         unit,
	}, nil
}

// ParseString parses bitrate from string in format "<quantity> <unit>".
// Quantity must be a non-negative number and it cannot be zero with a
// fractional part. Valid units are "bit/s", "kbit/s", "Mbit/s", "Gbit/s",
// "Tbit/s". If unit is not specified, it is assumed to be minimum unit of
// measurement ("bit/s"). The resulting base quantity (in bits per second) must
// be an integer. Spacing between quantity and unit are allowed and ignored. The
// same for leading and trailing whitespace.
func ParseString(s string) (Bitrate, error) {
	baseQuantity, unit, err := internal.ParseString(s, bitrateRegexp, parseUnit)
	if err != nil {
		return Bitrate{}, err
	}

	return Bitrate{
		baseQuantity: baseQuantity,
		unit:         unit,
		raw:          &s,
	}, nil
}

// MustNewFromInt64 is like [NewFromInt64] but panics on error.
func MustNewFromInt64(quantity int64, unit Unit) Bitrate {
	return must.Value(NewFromInt64(quantity, unit))
}

// MustNewFromBigInt is like [NewFromBigInt] but panics on error.
func MustNewFromBigInt(quantity *big.Int, unit Unit) Bitrate {
	return must.Value(NewFromBigInt(quantity, unit))
}

// MustParseString is like [ParseString] but panics on error.
func MustParseString(s string) Bitrate {
	return must.Value(ParseString(s))
}

// BigInt returns bitrate value in the minimum unit of measurement ("bit/s").
func (b Bitrate) BigInt() *big.Int {
	if b.baseQuantity == nil {
		return big.NewInt(0)
	}
	return b.baseQuantity
}

// RawValue returns the raw string value from which the bitrate was parsed.
func (b Bitrate) RawValue() *string {
	return b.raw
}

// String returns a bitrate string (e.g. "10 bit/s", "10 Gbit/s").
func (b Bitrate) String() string {
	return b.QuantityString() + " " + b.unit.String()
}

// QuantityString returns a bitrate quantity string (e.g. "10" for "10 bit/s").
func (b Bitrate) QuantityString() string {
	if b.baseQuantity == nil {
		return "0"
	}
	if b.unit == BITS {
		return b.baseQuantity.String()
	}
	return decimal.NewFromBigInt(b.baseQuantity, 0).
		Div(decimal.NewFromBigInt(b.unit.GetValue(), 0)).
		String()
}

// Unit returns the bitrate unit.
func (b Bitrate) Unit() Unit {
	return b.unit
}

// Clone returns a new bitrate with the same quantity and raw string
// representation.
func (b *Bitrate) Clone() *Bitrate {
	if b == nil {
		return nil
	}
	return &Bitrate{
		baseQuantity: new(big.Int).Set(b.baseQuantity),
		unit:         b.unit,
		raw:          ptr.Clone(b.raw),
	}
}

// Equal checks if the values of b and b2 are equal. Note that this method
// checks both the value and the raw string representation. If you want to make
// a logical comparison, use the [Bitrate.Cmp].
func (b Bitrate) Equal(b2 Bitrate) bool {
	return b.unit == b2.unit &&
		b.baseQuantity.Cmp(b2.baseQuantity) == 0 &&
		ptr.Equal(b.raw, b2.raw)
}

// Cmp compares b and b2 and returns:
//   - -1 if b < b2;
//   - 0 if b == b2;
//   - +1 if b > b2.
func (b Bitrate) Cmp(b2 Bitrate) int {
	return b.BigInt().Cmp(b2.BigInt())
}

func (b Bitrate) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	b.Encode(&e)
	return e.Bytes(), nil
}

func (b *Bitrate) UnmarshalJSON(bytes []byte) error {
	return b.Decode(jx.DecodeBytes(bytes))
}

func (b *Bitrate) Encode(e *jx.Encoder) {
	switch {
	case b == nil:
		e.Null()
	case b.raw != nil:
		e.Str(*b.raw)
	default:
		e.Str(b.String())
	}
}

func (b *Bitrate) Decode(d *jx.Decoder) error {
	if b == nil {
		return consterr.Error("decode to nil")
	}

	raw, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseString(raw)
	if err != nil {
		return err
	}

	b.baseQuantity = parsed.baseQuantity
	b.unit = parsed.unit
	b.raw = parsed.raw
	return nil
}
