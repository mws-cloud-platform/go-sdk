// Package duration provides types and utilities for working with durations.
package duration

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

// ErrInvalidDurationString is returned when the duration string is invalid.
const ErrInvalidDurationString = consterr.Error("invalid duration string")

const (
	maxStrLen   = 128
	maxDuration = uint64(1<<63 - 1)

	// equivalent to minimum time.Duration value
	// https://github.com/golang/go/issues/48629
	minDuration = uint64(1 << 63)
)

// Duration represents a time duration. It wraps [time.Duration] and adds
// additional functionality for parsing and formatting. It also stores the raw
// string from which it was parsed.
type Duration struct {
	value time.Duration
	raw   *string
}

// NewFromTimeDuration creates a new duration from the [time.Duration] value.
func NewFromTimeDuration(d time.Duration) Duration {
	return Duration{value: d}
}

// ParseString parses duration from a raw string.
//
// The following formats are accepted:
//   - ISO-8601 duration format (e.g., "P6DT5H6M7S").
//   - Simple duration format (e.g., "10s", "1h 30m", "-15m", "-(1h 30m)").
//   - Seconds number format (e.g., "10").
//
// Additional rules and restrictions:
//   - The only allowed non-time designator for ISO-8601 duration format is
//     days (D). Years (Y), weeks (W), and months (M) are not allowed.
//   - Only seconds (S) in ISO-8601 format can have fractional part.
//   - Valid time units for the simple duration format are "ns", "us", "ms",
//     "s", "m", "h", "d". Units must be ordered from largest to smallest.
//   - Valid time units for the simple duration format are "d", "h", "m", "s",
//     "ms", "us", "ns". Units must be ordered from largest to smallest.
//   - Parsing is case-insensitive (e.g., "1h 30m" == "1H 30M").
//   - Spaces are allowed and ignored (e.g., "1h30m" == "1h 30m").
//   - A negative duration is prefixed with "-" and can be surrounded with
//     parentheses in simple duration format (e.g., "-(1h 30m)").
func ParseString(s string) (Duration, error) {
	if len(s) > maxStrLen {
		return Duration{}, fmt.Errorf("%w: string is too long", ErrInvalidDurationString)
	}

	v, err := parse(s)
	if err != nil {
		return Duration{}, err
	}

	return Duration{value: v, raw: &s}, nil
}

// MustParseString is the same as [ParseString] but panics on errors.
func MustParseString(s string) Duration {
	d, err := ParseString(s)
	if err != nil {
		panic(err)
	}
	return d
}

// ToTimeDuration returns the duration as a [time.Duration].
func (d Duration) ToTimeDuration() time.Duration {
	return d.value
}

// RawValue returns the raw string value from which the duration was parsed.
func (d Duration) RawValue() *string {
	return d.raw
}

// String returns duration string in a ISO-8601 duration format.
//
// The string is presented in the format "PdDThHmMs.fS" (e.g., "P6DT5H6M7S"),
// where d is the number of days, h is the number of hours, m is the number of
// minutes, s is the number of seconds and f is a fractional part of second.
// Negative durations are prefixed with the "-" sign (e.g., "-PT1M30S").
func (d Duration) String() string {
	if d.value == 0 {
		return "PT0S"
	}

	v := d.value
	carry := time.Duration(0)
	if v == time.Duration(math.MinInt64) {
		// handle MinInt64 edge case
		carry = 1
		v += 1
	}

	sb := strings.Builder{}
	if v < 0 {
		sb.WriteByte('-')
		v = -v
	}
	sb.WriteByte('P')

	if days := v / timeDay; days > 0 {
		sb.WriteString(strconv.Itoa(int(days)))
		sb.WriteByte('D')
		v %= timeDay
	}

	if v > 0 {
		sb.WriteByte('T')
	}

	if hours := v / time.Hour; hours > 0 {
		sb.WriteString(strconv.Itoa(int(hours)))
		sb.WriteByte('H')
		v %= time.Hour
	}
	if minutes := v / time.Minute; minutes > 0 {
		sb.WriteString(strconv.Itoa(int(minutes)))
		sb.WriteByte('M')
		v %= time.Minute
	}
	if (v + carry) > 0 {
		sb.WriteString(strconv.FormatFloat(float64(v+carry)/float64(time.Second), 'f', -1, 64))
		sb.WriteByte('S')
	}

	return sb.String()
}

// Equal returns true if the duration is equal to the other duration. Note that
// this method checks both the value and the raw string representation.
func (d Duration) Equal(other Duration) bool {
	return d.value == other.value && ptr.Equal(d.raw, other.raw)
}

// Clone returns a new duration with the same value and raw string
// representation.
func (d *Duration) Clone() *Duration {
	if d == nil {
		return nil
	}
	return &Duration{
		value: d.value,
		raw:   ptr.Clone(d.raw),
	}
}

func (d Duration) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	d.Encode(&e)
	return e.Bytes(), nil
}

func (d *Duration) UnmarshalJSON(bytes []byte) error {
	return d.Decode(jx.DecodeBytes(bytes))
}

func (d *Duration) Encode(e *jx.Encoder) {
	switch {
	case d == nil:
		e.Null()
	case d.raw != nil:
		e.Str(*d.raw)
	default:
		e.Str(d.String())
	}
}

func (d *Duration) Decode(dec *jx.Decoder) error {
	if d == nil {
		return consterr.Error("decode to nil")
	}

	raw, err := dec.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseString(raw)
	if err != nil {
		return err
	}

	d.value = parsed.value
	d.raw = parsed.raw
	return nil
}

func parse(s string) (_ time.Duration, err error) {
	s = strings.ToLower(strings.ReplaceAll(s, " ", ""))
	if len(s) == 0 {
		return 0, fmt.Errorf("%w: empty string", ErrInvalidDurationString)
	}

	if _, parseErr := strconv.ParseFloat(s, 64); parseErr == nil {
		s += "s"
	}

	var d uint64

	hasSign := false
	isNegative := false
	switch s[0] {
	case '-':
		hasSign = true
		isNegative = true
	case '+':
		hasSign = true
	}
	if hasSign {
		s = s[1:]
	}
	if len(s) == 0 {
		return 0, ErrInvalidDurationString
	}

	if s[0] == 'p' {
		d, err = parseISO8601(s[1:])
	} else {
		if hasSign && s[0] == '(' && s[len(s)-1] == ')' {
			s = s[1 : len(s)-1]
		}
		d, err = parseSimple(s)
	}
	if err != nil {
		return 0, ErrInvalidDurationString
	}

	return uint64ToTimeDuration(d, isNegative)
}

//nolint:cyclop // complex parsing algorithm
func parseISO8601(s string) (d uint64, err error) {
	if len(s) == 0 {
		return 0, ErrInvalidDurationString
	}

	var (
		prevUnit        uint64
		isTimeComponent bool
	)

	for len(s) > 0 {
		var (
			n   number
			err error
		)

		if s[0] == 't' {
			isTimeComponent = true
			s = s[1:]
			if len(s) == 0 {
				return 0, ErrInvalidDurationString
			}
		}

		n, s, err = parseNumber(s)
		if err != nil {
			return 0, err
		}

		// Consume unit.
		if len(s) == 0 {
			return 0, ErrInvalidDurationString
		}
		u := s[:1]
		s = s[1:]
		unit, ok := isoUnitMap[u]
		if !ok {
			return 0, fmt.Errorf("%w: unknown unit %q", ErrInvalidDurationString, u)
		}
		if (prevUnit != 0 && unit >= prevUnit) || (u == "d" && isTimeComponent) || (u != "d" && !isTimeComponent) {
			return 0, fmt.Errorf("%w: unexpected unit position", ErrInvalidDurationString)
		}
		if n.fraction > 0 && u != "s" {
			return 0, fmt.Errorf("%w: fraction not allowed for unit %q", ErrInvalidDurationString, u)
		}
		prevUnit = unit

		d, err = n.addTo(d, unit)
		if err != nil {
			return 0, err
		}
	}

	return d, nil
}

// parseSimple is a modified version of [time.ParseDuration] that can handle
// days (e.g., "1d12h30m") and checks units order (e.g., "30m1h" is invalid).
func parseSimple(s string) (uint64, error) {
	if len(s) == 0 {
		return 0, ErrInvalidDurationString
	}

	var (
		d        uint64
		prevUnit uint64
	)

	for len(s) > 0 {
		var (
			n   number
			err error
		)

		n, s, err = parseNumber(s)
		if err != nil {
			return 0, err
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, fmt.Errorf("%w: missing unit", ErrInvalidDurationString)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, fmt.Errorf("%w: unknown unit %q", ErrInvalidDurationString, u)
		}
		if prevUnit != 0 && unit >= prevUnit {
			return 0, fmt.Errorf("%w: unexpected unit position", ErrInvalidDurationString)
		}
		prevUnit = unit

		d, err = n.addTo(d, unit)
		if err != nil {
			return 0, err
		}
	}

	return d, nil
}

func uint64ToTimeDuration(d uint64, isNegative bool) (time.Duration, error) {
	if isNegative {
		return -time.Duration(d), nil
	}
	if d > maxDuration {
		return 0, fmt.Errorf("%w: overflow", ErrInvalidDurationString)
	}
	return time.Duration(d), nil
}

type number struct {
	value    uint64  // integer part of the number
	fraction uint64  // fractional part of the number
	scale    float64 // value = v + f/scale
}

// https://github.com/golang/go/blob/go1.25.5/src/time/format.go#L1692-L1709
func (n number) addTo(d, unit uint64) (uint64, error) {
	v := n.value
	if v > minDuration/unit {
		return 0, ErrInvalidDurationString
	}
	v *= unit
	if n.fraction > 0 {
		// float64 is needed to be nanosecond accurate for fractions of hours.
		// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
		v += uint64(float64(n.fraction) * (float64(unit) / n.scale))
		if v > minDuration {
			return 0, ErrInvalidDurationString
		}
	}
	d += v
	if d > minDuration {
		return 0, ErrInvalidDurationString
	}
	return d, nil
}

// parseNumber consumes the number in format [0-9]+(\.[0-9]*)? from s.
func parseNumber(s string) (n number, rem string, err error) {
	// The next character must be [0-9.]
	if s[0] != '.' && (s[0] < '0' || s[0] > '9') {
		return number{}, "", ErrInvalidDurationString
	}

	pl := len(s)
	n.value, s, err = parseInt(s)
	if err != nil {
		return number{}, "", err
	}

	pre := pl != len(s) // whether we consumed anything before a period

	// Consume (\.[0-9]*)?
	post := false
	if s != "" && s[0] == '.' {
		s = s[1:]
		pl := len(s)
		n.fraction, n.scale, s = parseFraction(s)
		post = pl != len(s)
	}
	if !pre && !post {
		// no digits (e.g., ".s" or "-.s")
		return number{}, "", ErrInvalidDurationString
	}

	return n, s, nil
}

// parseInt consumes the leading [0-9]* from s. Cloned from [time.leadingInt].
func parseInt(s string) (x uint64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > minDuration/10 {
			return 0, rem, fmt.Errorf("%w: overflow", ErrInvalidDurationString)
		}
		x = x*10 + uint64(c) - '0'
		if x > minDuration {
			return 0, rem, fmt.Errorf("%w: overflow", ErrInvalidDurationString)
		}
	}
	return x, s[i:], nil
}

// parseFraction consumes the leading [0-9]* from s. It is used only for
// fractions, so does not return an error on overflow, it just stops
// accumulating precision. Cloned from [time.leadingFraction].
func parseFraction(s string) (x uint64, scale float64, rem string) {
	i := 0
	scale = 1
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if overflow {
			continue
		}
		if x > maxDuration/10 {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		y := x*10 + uint64(c) - '0'
		if y > minDuration {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return x, scale, s[i:]
}

const timeDay = 24 * time.Hour

var unitMap = map[string]uint64{
	"ns": uint64(time.Nanosecond),
	"us": uint64(time.Microsecond),
	"ms": uint64(time.Millisecond),
	"s":  uint64(time.Second),
	"m":  uint64(time.Minute),
	"h":  uint64(time.Hour),
	"d":  uint64(timeDay),
}

var isoUnitMap = map[string]uint64{
	"s": uint64(time.Second),
	"m": uint64(time.Minute),
	"h": uint64(time.Hour),
	"d": uint64(timeDay),
}
