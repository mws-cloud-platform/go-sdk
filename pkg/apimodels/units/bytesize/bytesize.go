// Package bytesize provides functionality for working with byte sizes.
package bytesize

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

var byteSizeRegexp = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)\s*([kmgtpezyrq]?b)?$`)

// ByteSize represents a byte size value.
type ByteSize struct {
	baseQuantity *big.Int
	unit         Unit
	raw          *string
}

// NewFromInt64 creates a new byte size value from int64.
func NewFromInt64(quantity int64, unit Unit) (ByteSize, error) {
	if quantity < 0 {
		return ByteSize{}, uniterrors.ErrNegativeQuantity
	}
	if !unit.IsValid() {
		return ByteSize{}, uniterrors.ErrUnexpectedUnit
	}

	return ByteSize{
		baseQuantity: new(big.Int).Mul(big.NewInt(quantity), unit.GetValue()),
		unit:         unit,
	}, nil
}

// NewFromBigInt creates a new byte size value from [big.Int]. If quantity is
// nil, it is set to zero.
func NewFromBigInt(quantity *big.Int, unit Unit) (ByteSize, error) {
	if quantity == nil {
		quantity = new(big.Int)
	}
	if quantity.Sign() == -1 {
		return ByteSize{}, uniterrors.ErrNegativeQuantity
	}
	if !unit.IsValid() {
		return ByteSize{}, uniterrors.ErrUnexpectedUnit
	}

	return ByteSize{
		baseQuantity: new(big.Int).Mul(quantity, unit.GetValue()),
		unit:         unit,
	}, nil
}

// ParseString parses byte size from string in format "<quantity> <unit>".
// Quantity must be a non-negative number and it cannot be zero with a
// fractional part. Valid units are "B", "KB", "MB", "GB", "TB", "PB", "EB",
// "ZB", "YB", "RB", "QB". If unit is not specified, it is assumed to be minimum
// unit of measurement ("B"). The resulting base quantity (in bytes) must be an
// integer. Spacing between quantity and unit are allowed and ignored. The same
// for leading and trailing whitespace.
func ParseString(s string) (ByteSize, error) {
	baseQuantity, unit, err := internal.ParseString(s, byteSizeRegexp, parseUnit)
	if err != nil {
		return ByteSize{}, err
	}

	return ByteSize{
		baseQuantity: baseQuantity,
		unit:         unit,
		raw:          &s,
	}, nil
}

// MustNewFromInt64 is like [NewFromInt64] but panics on error.
func MustNewFromInt64(quantity int64, unit Unit) ByteSize {
	return must.Value(NewFromInt64(quantity, unit))
}

// MustNewFromBigInt is like [NewFromBigInt] but panics on error.
func MustNewFromBigInt(quantity *big.Int, unit Unit) ByteSize {
	return must.Value(NewFromBigInt(quantity, unit))
}

// MustParseString is like [ParseString] but panics on error.
func MustParseString(s string) ByteSize {
	return must.Value(ParseString(s))
}

// BigInt returns byte size value in the minimum unit of measurement ("B").
func (b ByteSize) BigInt() *big.Int {
	if b.baseQuantity == nil {
		return big.NewInt(0)
	}
	return b.baseQuantity
}

// RawValue returns the raw string value from which the byte size was parsed.
func (b ByteSize) RawValue() *string {
	return b.raw
}

// String returns a byte size string (e.g. "10 B", "10 GB").
func (b ByteSize) String() string {
	return b.QuantityString() + " " + b.unit.String()
}

// QuantityString returns a byte size quantity string (e.g. "10" for "10 B").
func (b ByteSize) QuantityString() string {
	if b.baseQuantity == nil {
		return "0"
	}
	if b.unit == B {
		return b.baseQuantity.String()
	}
	return decimal.NewFromBigInt(b.baseQuantity, 0).
		Div(decimal.NewFromBigInt(b.unit.GetValue(), 0)).
		String()
}

// Unit returns the byte size unit.
func (b ByteSize) Unit() Unit {
	return b.unit
}

// Clone returns a new byte size with the same quantity and raw string
// representation.
func (b *ByteSize) Clone() *ByteSize {
	if b == nil {
		return nil
	}
	return &ByteSize{
		baseQuantity: new(big.Int).Set(b.baseQuantity),
		unit:         b.unit,
		raw:          ptr.Clone(b.raw),
	}
}

// Equal checks if the values of b and b2 are equal. Note that this method
// checks both the value and the raw string representation. If you want to make
// a logical comparison, use the [ByteSize.Cmp].
func (b ByteSize) Equal(b2 ByteSize) bool {
	return b.unit == b2.unit &&
		b.baseQuantity.Cmp(b2.baseQuantity) == 0 &&
		ptr.Equal(b.raw, b2.raw)
}

// Cmp compares b and b2 and returns:
//   - -1 if b < b2;
//   - 0 if b == b2;
//   - +1 if b > b2.
func (b ByteSize) Cmp(b2 ByteSize) int {
	return b.BigInt().Cmp(b2.BigInt())
}

func (b ByteSize) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	b.Encode(&e)
	return e.Bytes(), nil
}

func (b *ByteSize) UnmarshalJSON(bytes []byte) error {
	return b.Decode(jx.DecodeBytes(bytes))
}

func (b *ByteSize) Encode(e *jx.Encoder) {
	switch {
	case b == nil:
		e.Null()
	case b.raw != nil:
		e.Str(*b.raw)
	default:
		e.Str(b.String())
	}
}

func (b *ByteSize) Decode(d *jx.Decoder) error {
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
