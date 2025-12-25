// Package frequency provides functionality for working with frequencies.
package frequency

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

var frequencyRegexp = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)\s*([kmgt]?hz)?$`)

// Frequency represents a frequency value.
type Frequency struct {
	baseQuantity *big.Int
	unit         Unit
	raw          *string
}

// NewFromInt64 creates a new frequency value from int64.
func NewFromInt64(quantity int64, unit Unit) (Frequency, error) {
	if quantity < 0 {
		return Frequency{}, uniterrors.ErrNegativeQuantity
	}
	if !unit.IsValid() {
		return Frequency{}, uniterrors.ErrUnexpectedUnit
	}

	return Frequency{
		baseQuantity: new(big.Int).Mul(big.NewInt(quantity), unit.GetValue()),
		unit:         unit,
	}, nil
}

// NewFromBigInt creates a new frequency value from [big.Int]. If quantity is
// nil, it is set to zero.
func NewFromBigInt(quantity *big.Int, unit Unit) (Frequency, error) {
	if quantity == nil {
		quantity = new(big.Int)
	}
	if quantity.Sign() == -1 {
		return Frequency{}, uniterrors.ErrNegativeQuantity
	}
	if !unit.IsValid() {
		return Frequency{}, uniterrors.ErrUnexpectedUnit
	}

	return Frequency{
		baseQuantity: new(big.Int).Mul(quantity, unit.GetValue()),
		unit:         unit,
	}, nil
}

// ParseString parses frequency from string in format "<quantity> <unit>".
// Quantity must be a non-negative number and it cannot be zero with a
// fractional part. Valid units are "Hz", "kHz", "MHz", "GHz", "THz". If unit is
// not specified, it is assumed to be minimum unit of measurement ("Hz"). The
// resulting base quantity (in hertz) must be an integer. Spacing between
// quantity and unit are allowed and ignored. The same for leading and trailing
// whitespace.
func ParseString(s string) (Frequency, error) {
	baseQuantity, unit, err := internal.ParseString(s, frequencyRegexp, parseUnit)
	if err != nil {
		return Frequency{}, err
	}

	return Frequency{
		baseQuantity: baseQuantity,
		unit:         unit,
		raw:          &s,
	}, nil
}

// MustNewFromInt64 is like [NewFromInt64] but panics on error.
func MustNewFromInt64(quantity int64, unit Unit) Frequency {
	return must.Value(NewFromInt64(quantity, unit))
}

// MustNewFromBigInt is like [NewFromBigInt] but panics on error.
func MustNewFromBigInt(quantity *big.Int, unit Unit) Frequency {
	return must.Value(NewFromBigInt(quantity, unit))
}

// MustParseString is like [ParseString] but panics on error.
func MustParseString(s string) Frequency {
	return must.Value(ParseString(s))
}

// BigInt returns frequency value in the minimum unit of measurement ("Hz").
func (f Frequency) BigInt() *big.Int {
	if f.baseQuantity == nil {
		return big.NewInt(0)
	}
	return f.baseQuantity
}

// RawValue returns the raw string value from which the frequency was parsed.
func (f Frequency) RawValue() *string {
	return f.raw
}

// String returns a frequency string (e.g. "10 Hz", "10 GHz").
func (f Frequency) String() string {
	return f.QuantityString() + " " + f.unit.String()
}

// QuantityString returns a frequency quantity string (e.g. "10" for "10 Hz").
func (f Frequency) QuantityString() string {
	if f.baseQuantity == nil {
		return "0"
	}
	if f.unit == HZ {
		return f.baseQuantity.String()
	}
	return decimal.NewFromBigInt(f.baseQuantity, 0).
		Div(decimal.NewFromBigInt(f.unit.GetValue(), 0)).
		String()
}

// Unit returns the frequency unit.
func (f Frequency) Unit() Unit {
	return f.unit
}

// Clone returns a new frequency with the same quantity and raw string
// representation.
func (f *Frequency) Clone() *Frequency {
	if f == nil {
		return nil
	}
	return &Frequency{
		baseQuantity: new(big.Int).Set(f.baseQuantity),
		unit:         f.unit,
		raw:          ptr.Clone(f.raw),
	}
}

// Equal checks if the values of f and f2 are equal. Note that this method
// checks both the value and the raw string representation. If you want to make
// a logical comparison, use the [Frequency.Cmp].
func (f Frequency) Equal(f2 Frequency) bool {
	return f.unit == f2.unit &&
		f.baseQuantity.Cmp(f2.baseQuantity) == 0 &&
		ptr.Equal(f.raw, f2.raw)
}

// Cmp compares f and f2 and returns:
//   - -1 if f < f2;
//   - 0 if f == f2;
//   - +1 if f > f2.
func (f Frequency) Cmp(f2 Frequency) int {
	return f.BigInt().Cmp(f2.BigInt())
}

func (f Frequency) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	f.Encode(&e)
	return e.Bytes(), nil
}

func (f *Frequency) UnmarshalJSON(bytes []byte) error {
	return f.Decode(jx.DecodeBytes(bytes))
}

func (f *Frequency) Encode(e *jx.Encoder) {
	switch {
	case f == nil:
		e.Null()
	case f.raw != nil:
		e.Str(*f.raw)
	default:
		e.Str(f.String())
	}
}

func (f *Frequency) Decode(d *jx.Decoder) error {
	if f == nil {
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

	f.baseQuantity = parsed.baseQuantity
	f.unit = parsed.unit
	f.raw = parsed.raw
	return nil
}
