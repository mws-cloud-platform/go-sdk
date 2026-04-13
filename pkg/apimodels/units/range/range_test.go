package rangeunit

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/pkg/apimodels/units/bitrate"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
)

const (
	rawValue = "5-10b"
)

var (
	bitrate5bps  = bitrate.MustParseString("5bit/s")
	bitrate10bps = bitrate.MustParseString("10bit/s")
)

func TestRange_New(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		minValue    bitrate.Bitrate
		maxValue    bitrate.Bitrate
		expected    Range[bitrate.Bitrate]
		errExpected bool
	}{
		{
			name:     "valid_range",
			minValue: bitrate5bps,
			maxValue: bitrate10bps,
			expected: Range[bitrate.Bitrate]{
				minValue: bitrate5bps,
				maxValue: bitrate10bps,
			},
		},
		{
			name:     "equal_values",
			minValue: bitrate5bps,
			maxValue: bitrate5bps,
			expected: Range[bitrate.Bitrate]{
				minValue: bitrate5bps,
				maxValue: bitrate5bps,
			},
		},
		{
			name:        "min_value_greater_than_max_value",
			minValue:    bitrate10bps,
			maxValue:    bitrate5bps,
			errExpected: true,
		},
		{
			name:        "different_units",
			minValue:    bitrate5bps,
			maxValue:    bitrate.MustParseString("10Gbit/s"),
			errExpected: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := NewRange[bitrate.Bitrate](testCase.minValue, testCase.maxValue)
			if testCase.errExpected {
				require.Error(t, err)
				return
			}

			require.Equal(t, testCase.expected, result)
		})
	}
}

func TestRange_MaxValue(t *testing.T) {
	result, err := NewRange[bitrate.Bitrate](bitrate5bps, bitrate10bps)
	require.NoError(t, err)

	require.Equal(t, bitrate10bps, result.MaxValue())
}

func TestRange_MinValue(t *testing.T) {
	result, err := NewRange[bitrate.Bitrate](bitrate5bps, bitrate10bps)
	require.NoError(t, err)

	require.Equal(t, bitrate5bps, result.MinValue())
}

func TestRange_RawValue(t *testing.T) {
	result, err := ParseString[bytesize.ByteSize](rawValue)
	require.NoError(t, err)

	require.Equal(t, rawValue, ptr.Value(result.RawValue()))
}

func TestRange_Clone(t *testing.T) {
	r, err := ParseString[bytesize.ByteSize](rawValue)
	require.NoError(t, err)

	clone := r.Clone()
	*r.rawValue = "rawValue"
	r.minValue = bytesize.MustParseString("1b")
	r.maxValue = bytesize.MustParseString("2b")

	require.NotEqual(t, r.rawValue, clone.rawValue)
	require.NotEqual(t, r.minValue, clone.minValue)
	require.NotEqual(t, r.maxValue, clone.maxValue)
}

func TestRange_Equal(t *testing.T) {
	for _, testCase := range []struct {
		name   string
		input1 Range[bitrate.Bitrate]
		input2 Range[bitrate.Bitrate]
		equal  bool
	}{
		{
			name: "equal_ranges",
			input1: Range[bitrate.Bitrate]{
				minValue: bitrate5bps,
				maxValue: bitrate10bps,
				rawValue: ptr.Get("5-10bit/s"),
			},
			input2: Range[bitrate.Bitrate]{
				minValue: bitrate5bps,
				maxValue: bitrate10bps,
				rawValue: ptr.Get("5-10bit/s"),
			},
			equal: true,
		},
		{
			name: "not_equal_min_value",
			input1: Range[bitrate.Bitrate]{
				minValue: bitrate5bps,
				maxValue: bitrate10bps,
				rawValue: ptr.Get("5-10bit/s"),
			},
			input2: Range[bitrate.Bitrate]{
				minValue: bitrate10bps,
				maxValue: bitrate10bps,
				rawValue: ptr.Get("5-10bit/s"),
			},
			equal: false,
		},
		{
			name: "not_equal_max_value",
			input1: Range[bitrate.Bitrate]{
				minValue: bitrate5bps,
				maxValue: bitrate10bps,
				rawValue: ptr.Get("5-10bit/s"),
			},
			input2: Range[bitrate.Bitrate]{
				minValue: bitrate5bps,
				maxValue: bitrate5bps,
				rawValue: ptr.Get("5-10bit/s"),
			},
			equal: false,
		},
		{
			name: "not_equal_raw_value",
			input1: Range[bitrate.Bitrate]{
				minValue: bitrate5bps,
				maxValue: bitrate10bps,
				rawValue: ptr.Get("5-10bit/s"),
			},
			input2: Range[bitrate.Bitrate]{
				minValue: bitrate5bps,
				maxValue: bitrate10bps,
				rawValue: ptr.Get("5-10 bit/s"),
			},
			equal: false,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.equal, testCase.input1.Equal(testCase.input2))
		})
	}
}

func TestRange_String(t *testing.T) {
	r, err := ParseString[bytesize.ByteSize](rawValue)
	require.NoError(t, err)
	require.Equal(t, "5 - 10B", r.String())

	r = Range[bytesize.ByteSize]{}
	require.Equal(t, "0 - 0B", r.String())
}

func TestRange_MarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/marshal_json.golden"), golden.WithRecreateOnUpdate())

	r, err := NewRange[bitrate.Bitrate](bitrate.MustParseString("10bit/s"), bitrate.MustParseString("20bit/s"))
	require.NoError(t, err)

	result, err := json.Marshal(r)
	require.NoError(t, err)

	dir.String(t, "10-20.txt", string(result))
}

func TestRange_MarshalJSONWithRawValue(t *testing.T) {
	r, err := NewRange(bytesize.MustParseString("5b"), bytesize.MustParseString("10b"))
	require.NoError(t, err)

	r.rawValue = ptr.Get(rawValue)

	result, err := json.Marshal(r)
	require.NoError(t, err)
	require.Equal(t, string(result), strconv.Quote(rawValue))
}

func TestRange_UnmarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/unmarshal_json.golden"), golden.WithRecreateOnUpdate())

	for _, testCase := range []struct {
		name        string
		rawValue    string
		errExpected bool
	}{
		{
			name:        "invalid unit",
			rawValue:    "1-2b",
			errExpected: true,
		},
		{
			name:        "revert range",
			rawValue:    "2-1bit/s",
			errExpected: true,
		},
		{
			name:        "negative unit",
			rawValue:    "-1-2bit/s",
			errExpected: true,
		},
		{
			name:        "empty string",
			rawValue:    "",
			errExpected: true,
		},
		{
			name:        "invalid string 1",
			rawValue:    "1",
			errExpected: true,
		},
		{
			name:        "invalid string 2",
			rawValue:    "1----------2bit/s",
			errExpected: true,
		},
		{
			name:        "invalid string 3",
			rawValue:    "1bit/s-2bit/s",
			errExpected: true,
		},
		{
			name:        "invalid string 4",
			rawValue:    "1bit/s2bit/s",
			errExpected: true,
		},
		{
			name:     "normal range",
			rawValue: "5-10bit/s",
		},
		{
			name:     "range with spaces",
			rawValue: "   5   -     10     bit/s   ",
		},
		{
			name:     "range without unit",
			rawValue: "5 - 10   ",
		},
		{
			name:     "equal units",
			rawValue: "1-1bit/s",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var r Range[bitrate.Bitrate]

			err := json.Unmarshal([]byte(strconv.Quote(testCase.rawValue)), &r)
			if testCase.errExpected {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			dir.String(t, strings.ReplaceAll(testCase.name, " ", "_")+".txt",
				"minValue: "+strconv.Quote(r.minValue.String())+"\nmaxValue: "+strconv.Quote(r.maxValue.String()))
		})
	}
}

func ExampleNewRange() {
	parsedRange, err := NewRange[bytesize.ByteSize](bytesize.MustParseString("10GB"), bytesize.MustParseString("20GB"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(parsedRange.String())
	fmt.Println(ptr.Value(parsedRange.RawValue()))
	fmt.Println(parsedRange.MinValue().String())
	fmt.Println(parsedRange.MaxValue().String())
	// Output:
	// 10 - 20GB
	//
	// 10 GB
	// 20 GB
}

func ExampleParseString() {
	parsedRange, err := ParseString[bytesize.ByteSize]("10-20   GB")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(parsedRange.String())
	fmt.Println(ptr.Value(parsedRange.RawValue()))
	fmt.Println(parsedRange.MinValue().String())
	fmt.Println(parsedRange.MaxValue().String())
	// Output:
	// 10 - 20GB
	// 10-20   GB
	// 10 GB
	// 20 GB
}

func ExampleRange_UnmarshalJSON() {
	var parsedRange Range[bytesize.ByteSize]
	if err := json.Unmarshal([]byte(strconv.Quote("10-20   GB")), &parsedRange); err != nil {
		log.Fatal(err)
	}

	fmt.Println(parsedRange.String())
	fmt.Println(ptr.Value(parsedRange.RawValue()))
	fmt.Println(parsedRange.MinValue().String())
	fmt.Println(parsedRange.MaxValue().String())
	// Output:
	// 10 - 20GB
	// 10-20   GB
	// 10 GB
	// 20 GB
}
