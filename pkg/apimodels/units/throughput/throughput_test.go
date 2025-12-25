package throughput_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/apimodels/units/throughput"
)

func ExampleNewFromInt64() {
	v, err := throughput.NewFromInt64(256, throughput.BPS)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 256 Bps
}

func ExampleNewFromBigInt() {
	v, err := throughput.NewFromBigInt(big.NewInt(100), throughput.KBPS)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 100 KBps
}

func ExampleParseString() {
	v, err := throughput.ParseString("20 GBps")
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 20 GBps
}

func ExampleThroughput_UnmarshalJSON() {
	var v throughput.Throughput
	if err := json.Unmarshal([]byte(`"256 MBps"`), &v); err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 256 MBps
}

func TestNewFromInt64(t *testing.T) {
	for _, tc := range []struct {
		quantity int64
		unit     throughput.Unit
		expected string
	}{
		{0, throughput.BPS, "0 Bps"},
		{1, throughput.BPS, "1 Bps"},
		{2, throughput.KBPS, "2 KBps"},
		{4, throughput.MBPS, "4 MBps"},
		{8, throughput.GBPS, "8 GBps"},
		{16, throughput.TBPS, "16 TBps"},
	} {
		t.Run(tc.expected, func(t *testing.T) {
			v, err := throughput.NewFromInt64(tc.quantity, tc.unit)
			require.NoError(t, err)
			require.Equal(t, tc.expected, v.String())
		})
	}
}

func TestNewFromInt64_negative(t *testing.T) {
	_, err := throughput.NewFromInt64(-1, throughput.BPS)
	require.Error(t, err)
}

func TestNewFromBigInt(t *testing.T) {
	for _, tc := range []struct {
		quantity *big.Int
		unit     throughput.Unit
		expected string
	}{
		{big.NewInt(0), throughput.BPS, "0 Bps"},
		{big.NewInt(1), throughput.BPS, "1 Bps"},
		{big.NewInt(2), throughput.KBPS, "2 KBps"},
		{big.NewInt(4), throughput.MBPS, "4 MBps"},
		{big.NewInt(8), throughput.GBPS, "8 GBps"},
		{big.NewInt(16), throughput.TBPS, "16 TBps"},
	} {
		t.Run(tc.expected, func(t *testing.T) {
			v, err := throughput.NewFromBigInt(tc.quantity, tc.unit)
			require.NoError(t, err)
			require.Equal(t, tc.expected, v.String())
		})
	}
}

func TestNewFromBigInt_negative(t *testing.T) {
	_, err := throughput.NewFromBigInt(big.NewInt(-1), throughput.BPS)
	require.Error(t, err)
}

func TestParseString(t *testing.T) {
	for _, tc := range []struct {
		input    string
		quantity *big.Int
		unit     throughput.Unit
	}{
		{"0", big.NewInt(0), throughput.BPS},
		{"1", big.NewInt(1), throughput.BPS},
		{" 256  ", big.NewInt(256), throughput.BPS},
		{"1 bps", big.NewInt(1), throughput.BPS},
		{"   1   bps  ", big.NewInt(1), throughput.BPS},
		{"1 BPS", big.NewInt(1), throughput.BPS},
		{"100bps", big.NewInt(100), throughput.BPS},
		{"2 kbps", big.NewInt(2048), throughput.KBPS},
		{"6.5 kbps", big.NewInt(6656), throughput.KBPS},
		{"1.000000000", big.NewInt(1), throughput.BPS},
		{"5 mbps", big.NewInt(5242880), throughput.MBPS},
		{"10 gbps", big.NewInt(10737418240), throughput.GBPS},
		{"3 tbps", big.NewInt(3298534883328), throughput.TBPS},
	} {
		t.Run(tc.input, func(t *testing.T) {
			v, err := throughput.ParseString(tc.input)
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
		"-5 kbps",
		"100 kb ps",
		"6000 000 bps",
		"0.5 bps",
		"5.0001 kbps",
	} {
		t.Run(tc, func(t *testing.T) {
			_, err := throughput.ParseString(tc)
			require.Error(t, err)
		})
	}
}

func TestThroughput_MarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		value    throughput.Throughput
		expected string
	}{
		{throughput.MustNewFromInt64(1, throughput.GBPS), `"1 GBps"`},
		{throughput.MustNewFromInt64(100, throughput.KBPS), `"100 KBps"`},
		{throughput.MustParseString("1 kbps"), `"1 kbps"`},
		{throughput.MustParseString("1KBPS"), `"1KBPS"`},
	} {
		t.Run(tc.value.String(), func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(data))
		})
	}
}

func TestThroughput_UnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		input string
		value *big.Int
		unit  throughput.Unit
	}{
		{`"1 gbps"`, big.NewInt(1073741824), throughput.GBPS},
		{`"100 bps"`, big.NewInt(100), throughput.BPS},
		{`"2 kbps"`, big.NewInt(2048), throughput.KBPS},
		{`"1024KBPS"`, big.NewInt(1048576), throughput.KBPS},
	} {
		t.Run(tc.input, func(t *testing.T) {
			var v throughput.Throughput
			err := json.Unmarshal([]byte(tc.input), &v)
			require.NoError(t, err)
			require.Equal(t, tc.value, v.BigInt())
			require.Equal(t, strings.Trim(tc.input, `"`), *v.RawValue())
		})
	}
}

func TestThroughput_Clone(t *testing.T) {
	v := throughput.MustParseString("10 GBps")
	clone := v.Clone()

	require.Zero(t, v.BigInt().Cmp(clone.BigInt()))
	require.Equal(t, v.Unit(), clone.Unit())
	require.Equal(t, *v.RawValue(), *clone.RawValue())

	// Ensure they are separate instances
	require.NotSame(t, v.BigInt(), clone.BigInt())
	require.NotSame(t, v.RawValue(), clone.RawValue())
}

func TestThroughput_Equal(t *testing.T) {
	v1 := throughput.MustNewFromInt64(10, throughput.GBPS)
	v2 := throughput.MustNewFromInt64(10, throughput.GBPS)
	v3 := throughput.MustNewFromInt64(20, throughput.GBPS)
	v4 := throughput.MustNewFromInt64(10, throughput.KBPS)
	v5 := throughput.MustParseString("10 GBps")

	require.True(t, v1.Equal(v2))
	require.False(t, v1.Equal(v3))
	require.False(t, v1.Equal(v4))
	require.False(t, v1.Equal(v5)) // Different because v1 has no raw value
}

func TestThroughput_Cmp(t *testing.T) {
	v1 := throughput.MustNewFromInt64(10, throughput.GBPS)
	v2 := throughput.MustNewFromInt64(10, throughput.GBPS)
	v3 := throughput.MustNewFromInt64(20, throughput.GBPS)
	v4 := throughput.MustNewFromInt64(10, throughput.MBPS)
	v5 := throughput.MustNewFromInt64(5, throughput.KBPS)

	require.Equal(t, 0, v1.Cmp(v2))
	require.Equal(t, -1, v1.Cmp(v3))
	require.Equal(t, 1, v1.Cmp(v4))
	require.Equal(t, 1, v1.Cmp(v5))
}
