// Package throughput provides functionality for working with throughput.
package throughput

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

var throughputRegexp = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)\s*([kmgtpezyrq]?bps)?$`)

// Throughput represents a throughput value.
type Throughput struct {
	baseQuantity *big.Int
	unit         Unit
	raw          *string
}

// NewFromInt64 creates a new throughput value from int64.
func NewFromInt64(quantity int64, unit Unit) (Throughput, error) {
	if quantity < 0 {
		return Throughput{}, uniterrors.ErrNegativeQuantity
	}
	if !unit.IsValid() {
		return Throughput{}, uniterrors.ErrUnexpectedUnit
	}

	return Throughput{
		baseQuantity: new(big.Int).Mul(big.NewInt(quantity), unit.GetValue()),
		unit:         unit,
	}, nil
}

// NewFromBigInt creates a new throughput value from [big.Int]. If quantity is
// nil, it is set to zero.
func NewFromBigInt(quantity *big.Int, unit Unit) (Throughput, error) {
	if quantity == nil {
		quantity = new(big.Int)
	}
	if quantity.Sign() == -1 {
		return Throughput{}, uniterrors.ErrNegativeQuantity
	}
	if !unit.IsValid() {
		return Throughput{}, uniterrors.ErrUnexpectedUnit
	}

	return Throughput{
		baseQuantity: new(big.Int).Mul(quantity, unit.GetValue()),
		unit:         unit,
	}, nil
}

// ParseString parses throughput from string in format "<quantity> <unit>".
// Quantity must be a non-negative number and it cannot be zero with a
// fractional part. Valid units are "Bps", "KBps", "MBps", "GBps", "TBps",
// "PBps", "EBps", "ZBps", "YBps", "RBps", "QBps". If unit is not specified, it
// is assumed to be minimum unit of measurement ("Bps"). The resulting base
// quantity (in bytes per second) must be an integer. Spacing between quantity
// and unit are allowed and ignored. The same for leading and trailing
// whitespace.
func ParseString(s string) (Throughput, error) {
	baseQuantity, unit, err := internal.ParseString(s, throughputRegexp, parseUnit)
	if err != nil {
		return Throughput{}, err
	}

	return Throughput{
		baseQuantity: baseQuantity,
		unit:         unit,
		raw:          &s,
	}, nil
}

// MustNewFromInt64 is like [NewFromInt64] but panics on error.
func MustNewFromInt64(quantity int64, unit Unit) Throughput {
	return must.Value(NewFromInt64(quantity, unit))
}

// MustNewFromBigInt is like [NewFromBigInt] but panics on error.
func MustNewFromBigInt(quantity *big.Int, unit Unit) Throughput {
	return must.Value(NewFromBigInt(quantity, unit))
}

// MustParseString is like [ParseString] but panics on error.
func MustParseString(s string) Throughput {
	return must.Value(ParseString(s))
}

// BigInt returns throughput value in the minimum unit of measurement ("Bps").
func (t Throughput) BigInt() *big.Int {
	if t.baseQuantity == nil {
		return big.NewInt(0)
	}
	return t.baseQuantity
}

// RawValue returns the raw string value from which the throughput was parsed.
func (t Throughput) RawValue() *string {
	return t.raw
}

// String returns a throughput string (e.g. "10 Bps", "10 GBps").
func (t Throughput) String() string {
	return t.QuantityString() + " " + t.unit.String()
}

// QuantityString returns a throughput quantity string (e.g. "10" for "10 Bps").
func (t Throughput) QuantityString() string {
	if t.baseQuantity == nil {
		return "0"
	}
	if t.unit == BPS {
		return t.baseQuantity.String()
	}
	return decimal.NewFromBigInt(t.baseQuantity, 0).
		Div(decimal.NewFromBigInt(t.unit.GetValue(), 0)).
		String()
}

// Unit returns the throughput unit.
func (t Throughput) Unit() Unit {
	return t.unit
}

// Clone returns a new throughput with the same quantity and raw string
// representation.
func (t *Throughput) Clone() *Throughput {
	if t == nil {
		return nil
	}
	return &Throughput{
		baseQuantity: new(big.Int).Set(t.baseQuantity),
		unit:         t.unit,
		raw:          ptr.Clone(t.raw),
	}
}

// Equal checks if the values of t and t2 are equal. Note that this method
// checks both the value and the raw string representation. If you want to make
// a logical comparison, use the [Throughput.Cmp].
func (t Throughput) Equal(t2 Throughput) bool {
	return t.unit == t2.unit &&
		t.baseQuantity.Cmp(t2.baseQuantity) == 0 &&
		ptr.Equal(t.raw, t2.raw)
}

// Cmp compares t and t2 and returns:
//   - -1 if t < t2;
//   - 0 if t == t2;
//   - +1 if t > t2.
func (t Throughput) Cmp(t2 Throughput) int {
	return t.BigInt().Cmp(t2.BigInt())
}

func (t Throughput) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	t.Encode(&e)
	return e.Bytes(), nil
}

func (t *Throughput) UnmarshalJSON(bytes []byte) error {
	return t.Decode(jx.DecodeBytes(bytes))
}

func (t *Throughput) Encode(e *jx.Encoder) {
	switch {
	case t == nil:
		e.Null()
	case t.raw != nil:
		e.Str(*t.raw)
	default:
		e.Str(t.String())
	}
}

func (t *Throughput) Decode(d *jx.Decoder) error {
	if t == nil {
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

	t.baseQuantity = parsed.baseQuantity
	t.unit = parsed.unit
	t.raw = parsed.raw
	return nil
}
