// Package largenumber provides functionality for working with large numbers.
package largenumber

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

var largenumberRegexp = regexp.MustCompile(`^(-?[0-9]+(?:\.[0-9]+)?)\s*(k)?$`)

// LargeNumber represents a large number.
type LargeNumber struct {
	baseQuantity *big.Int
	unit         Unit
	raw          *string
}

// NewFromInt64 creates a new large number from int64.
func NewFromInt64(quantity int64, unit Unit) (LargeNumber, error) {
	if quantity < 0 {
		return LargeNumber{}, uniterrors.ErrNegativeQuantity
	}
	if !unit.IsValid() {
		return LargeNumber{}, uniterrors.ErrUnexpectedUnit
	}

	return LargeNumber{
		baseQuantity: new(big.Int).Mul(big.NewInt(quantity), unit.GetValue()),
		unit:         unit,
	}, nil
}

// NewFromBigInt creates a new large number from [big.Int]. If quantity is nil,
// it is set to zero.
func NewFromBigInt(quantity *big.Int, unit Unit) (LargeNumber, error) {
	if quantity == nil {
		quantity = new(big.Int)
	}
	if quantity.Sign() == -1 {
		return LargeNumber{}, uniterrors.ErrNegativeQuantity
	}
	if !unit.IsValid() {
		return LargeNumber{}, uniterrors.ErrUnexpectedUnit
	}

	return LargeNumber{
		baseQuantity: new(big.Int).Mul(quantity, unit.GetValue()),
		unit:         unit,
	}, nil
}

// ParseString parses value from string in format "<quantity> <unit>". Quantity
// must be a non-negative number and it cannot be zero with a fractional part.
// Valid unit is "K". If unit is not specified, it is assumed to be minimum unit
// of measurement (empty unit). The resulting base quantity must be an integer.
// Spacing between quantity and unit are allowed and ignored. The same for
// leading and trailing whitespace.
func ParseString(s string) (LargeNumber, error) {
	baseQuantity, unit, err := internal.ParseString(s, largenumberRegexp, parseUnit)
	if err != nil {
		return LargeNumber{}, err
	}

	return LargeNumber{
		baseQuantity: baseQuantity,
		unit:         unit,
		raw:          &s,
	}, nil
}

// MustNewFromInt64 is like [NewFromInt64] but panics on error.
func MustNewFromInt64(quantity int64, unit Unit) LargeNumber {
	return must.Value(NewFromInt64(quantity, unit))
}

// MustNewFromBigInt is like [NewFromBigInt] but panics on error.
func MustNewFromBigInt(quantity *big.Int, unit Unit) LargeNumber {
	return must.Value(NewFromBigInt(quantity, unit))
}

// MustParseString is like [ParseString] but panics on error.
func MustParseString(s string) LargeNumber {
	return must.Value(ParseString(s))
}

// BigInt returns large number in the minimum unit of measurement (empty unit).
func (m LargeNumber) BigInt() *big.Int {
	if m.baseQuantity == nil {
		return big.NewInt(0)
	}
	return m.baseQuantity
}

// RawValue returns the raw string value from which the large number was parsed.
func (m LargeNumber) RawValue() *string {
	return m.raw
}

// String returns a large number string (e.g. "10", "10 K").
func (m LargeNumber) String() string {
	unitStr := m.unit.String()
	if unitStr == "" {
		return m.QuantityString()
	}
	return m.QuantityString() + " " + unitStr
}

// QuantityString returns a large number quantity string (e.g. "10" for "10 K").
func (m LargeNumber) QuantityString() string {
	if m.baseQuantity == nil {
		return "0"
	}
	if m.unit == EmptyUnit {
		return m.baseQuantity.String()
	}
	return decimal.NewFromBigInt(m.baseQuantity, 0).
		Div(decimal.NewFromBigInt(m.unit.GetValue(), 0)).
		String()
}

// Unit returns the large number unit.
func (m LargeNumber) Unit() Unit {
	return m.unit
}

// Clone returns a new large number with the same quantity and raw string
// representation.
func (m *LargeNumber) Clone() *LargeNumber {
	if m == nil {
		return nil
	}
	return &LargeNumber{
		baseQuantity: new(big.Int).Set(m.baseQuantity),
		unit:         m.unit,
		raw:          ptr.Clone(m.raw),
	}
}

// Equal checks if the values of m and m2 are equal. Note that this method
// checks both the value and the raw string representation. If you want to make
// a logical comparison, use the [LargeNumber.Cmp].
func (m LargeNumber) Equal(m2 LargeNumber) bool {
	return m.unit == m2.unit &&
		m.baseQuantity.Cmp(m2.baseQuantity) == 0 &&
		ptr.Equal(m.raw, m2.raw)
}

// Cmp compares m and m2 and returns:
//   - -1 if m < m2;
//   - 0 if m == m2;
//   - +1 if m > m2.
func (m LargeNumber) Cmp(m2 LargeNumber) int {
	return m.BigInt().Cmp(m2.BigInt())
}

func (m LargeNumber) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	m.Encode(&e)
	return e.Bytes(), nil
}

func (m *LargeNumber) UnmarshalJSON(bytes []byte) error {
	return m.Decode(jx.DecodeBytes(bytes))
}

func (m *LargeNumber) Encode(e *jx.Encoder) {
	switch {
	case m == nil:
		e.Null()
	case m.raw != nil:
		e.Str(*m.raw)
	default:
		e.Str(m.String())
	}
}

func (m *LargeNumber) Decode(d *jx.Decoder) error {
	if m == nil {
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

	m.baseQuantity = parsed.baseQuantity
	m.unit = parsed.unit
	m.raw = parsed.raw
	return nil
}
