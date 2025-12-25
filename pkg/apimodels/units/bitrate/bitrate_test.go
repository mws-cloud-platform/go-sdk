package bitrate_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/apimodels/units/bitrate"
)

func ExampleNewFromInt64() {
	v, err := bitrate.NewFromInt64(256, bitrate.BITS)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 256 bit/s
}

func ExampleNewFromBigInt() {
	v, err := bitrate.NewFromBigInt(big.NewInt(100), bitrate.KBITS)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 100 kbit/s
}

func ExampleParseString() {
	v, err := bitrate.ParseString("20 Gbit/s")
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 20 Gbit/s
}

func ExampleBitrate_UnmarshalJSON() {
	var v bitrate.Bitrate
	if err := json.Unmarshal([]byte(`"256 Mbit/s"`), &v); err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 256 Mbit/s
}

func TestNewFromInt64(t *testing.T) {
	for _, tc := range []struct {
		quantity int64
		unit     bitrate.Unit
		expected string
	}{
		{0, bitrate.BITS, "0 bit/s"},
		{1, bitrate.BITS, "1 bit/s"},
		{2, bitrate.KBITS, "2 kbit/s"},
		{4, bitrate.MBITS, "4 Mbit/s"},
		{8, bitrate.GBITS, "8 Gbit/s"},
		{16, bitrate.TBITS, "16 Tbit/s"},
	} {
		t.Run(tc.expected, func(t *testing.T) {
			v, err := bitrate.NewFromInt64(tc.quantity, tc.unit)
			require.NoError(t, err)
			require.Equal(t, tc.expected, v.String())
		})
	}
}

func TestNewFromInt64_negative(t *testing.T) {
	_, err := bitrate.NewFromInt64(-1, bitrate.BITS)
	require.Error(t, err)
}

func TestNewFromBigInt(t *testing.T) {
	for _, tc := range []struct {
		quantity *big.Int
		unit     bitrate.Unit
		expected string
	}{
		{big.NewInt(0), bitrate.BITS, "0 bit/s"},
		{big.NewInt(1), bitrate.BITS, "1 bit/s"},
		{big.NewInt(2), bitrate.KBITS, "2 kbit/s"},
		{big.NewInt(4), bitrate.MBITS, "4 Mbit/s"},
		{big.NewInt(8), bitrate.GBITS, "8 Gbit/s"},
		{big.NewInt(16), bitrate.TBITS, "16 Tbit/s"},
	} {
		t.Run(tc.expected, func(t *testing.T) {
			v, err := bitrate.NewFromBigInt(tc.quantity, tc.unit)
			require.NoError(t, err)
			require.Equal(t, tc.expected, v.String())
		})
	}
}

func TestNewFromBigInt_negative(t *testing.T) {
	_, err := bitrate.NewFromBigInt(big.NewInt(-1), bitrate.BITS)
	require.Error(t, err)
}

func TestParseString(t *testing.T) {
	for _, tc := range []struct {
		input    string
		quantity *big.Int
		unit     bitrate.Unit
	}{
		{"0", big.NewInt(0), bitrate.BITS},
		{"1", big.NewInt(1), bitrate.BITS},
		{" 256  ", big.NewInt(256), bitrate.BITS},
		{"1 bit/s", big.NewInt(1), bitrate.BITS},
		{"   1   bit/s  ", big.NewInt(1), bitrate.BITS},
		{"1 BIT/S", big.NewInt(1), bitrate.BITS},
		{"100bit/s", big.NewInt(100), bitrate.BITS},
		{"2 kbit/s", big.NewInt(2 * 1 << 10), bitrate.KBITS},
		{"6.5 kbit/s", big.NewInt(6656), bitrate.KBITS},
		{"1.000000000 bit/s", big.NewInt(1), bitrate.BITS},
		{"5 mbit/s", big.NewInt(5 * 1 << 20), bitrate.MBITS},
		{"10 Gbit/s", big.NewInt(10 * 1 << 30), bitrate.GBITS},
		{"3 Tbit/s", big.NewInt(3 * 1 << 40), bitrate.TBITS},
	} {
		t.Run(tc.input, func(t *testing.T) {
			v, err := bitrate.ParseString(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.quantity, v.BigInt())
			require.Equal(t, tc.unit, v.Unit())
			require.Equal(t, tc.input, *v.RawValue())
		})
	}
}

func TestParseString_invalid(t *testing.T) {
	for _, tc := range []string{
		"",
		"invalid",
		"-1",
		"-5 kbit/s",
		"100 k bit/s",
		"6000 000 bit/s",
		"0.5 bit/s",
		"5.0001 kbit/s",
	} {
		t.Run(tc, func(t *testing.T) {
			_, err := bitrate.ParseString(tc)
			require.Error(t, err)
		})
	}
}

func TestBitrate_MarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		value    bitrate.Bitrate
		expected string
	}{
		{bitrate.MustNewFromInt64(1, bitrate.GBITS), `"1 Gbit/s"`},
		{bitrate.MustNewFromInt64(100, bitrate.KBITS), `"100 kbit/s"`},
		{bitrate.MustParseString("1 kbit/s"), `"1 kbit/s"`},
		{bitrate.MustParseString("1KBIT/S"), `"1KBIT/S"`},
	} {
		t.Run(tc.value.String(), func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(data))
		})
	}
}

func TestBitrate_UnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		input string
		value *big.Int
		unit  bitrate.Unit
	}{
		{`"1 Gbit/s"`, big.NewInt(1 << 30), bitrate.GBITS},
		{`"100 bit/s"`, big.NewInt(100), bitrate.KBITS},
		{`"2 kbit/s"`, big.NewInt(2 * 1 << 10), bitrate.KBITS},
		{`"1024KBIT/S"`, big.NewInt(1024 * 1 << 10), bitrate.KBITS},
	} {
		t.Run(tc.input, func(t *testing.T) {
			var v bitrate.Bitrate
			err := json.Unmarshal([]byte(tc.input), &v)
			require.NoError(t, err)
			require.Equal(t, tc.value, v.BigInt())
			require.Equal(t, strings.Trim(tc.input, `"`), *v.RawValue())
		})
	}
}

func TestBitrate_Clone(t *testing.T) {
	v := bitrate.MustParseString("10 Gbit/s")
	clone := v.Clone()

	require.Zero(t, v.BigInt().Cmp(clone.BigInt()))
	require.Equal(t, v.Unit(), clone.Unit())
	require.Equal(t, *v.RawValue(), *clone.RawValue())

	// Ensure they are separate instances
	require.NotSame(t, v.BigInt(), clone.BigInt())
	require.NotSame(t, v.RawValue(), clone.RawValue())
}

func TestBitrate_Equal(t *testing.T) {
	v1 := bitrate.MustNewFromInt64(10, bitrate.GBITS)
	v2 := bitrate.MustNewFromInt64(10, bitrate.GBITS)
	v3 := bitrate.MustNewFromInt64(20, bitrate.GBITS)
	v4 := bitrate.MustNewFromInt64(10, bitrate.KBITS)
	v5 := bitrate.MustParseString("10 Gbit/s")

	require.True(t, v1.Equal(v2))
	require.False(t, v1.Equal(v3))
	require.False(t, v1.Equal(v4))
	require.False(t, v1.Equal(v5)) // Different because v1 has no raw value
}

func TestBitrate_Cmp(t *testing.T) {
	v1 := bitrate.MustNewFromInt64(10, bitrate.GBITS)
	v2 := bitrate.MustNewFromInt64(10, bitrate.GBITS)
	v3 := bitrate.MustNewFromInt64(20, bitrate.GBITS)
	v4 := bitrate.MustNewFromInt64(10, bitrate.MBITS)
	v5 := bitrate.MustNewFromInt64(5, bitrate.KBITS)

	require.Equal(t, 0, v1.Cmp(v2))
	require.Equal(t, -1, v1.Cmp(v3))
	require.Equal(t, 1, v1.Cmp(v4))
	require.Equal(t, 1, v1.Cmp(v5))
}
