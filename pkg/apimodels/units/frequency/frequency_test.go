package frequency_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/apimodels/units/frequency"
)

func ExampleNewFromInt64() {
	v, err := frequency.NewFromInt64(256, frequency.HZ)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 256 Hz
}

func ExampleNewFromBigInt() {
	v, err := frequency.NewFromBigInt(big.NewInt(100), frequency.KHZ)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 100 kHz
}

func ExampleParseString() {
	v, err := frequency.ParseString("20 GHz")
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 20 GHz
}

func ExampleFrequency_UnmarshalJSON() {
	var v frequency.Frequency
	if err := json.Unmarshal([]byte(`"256 MHz"`), &v); err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 256 MHz
}

func TestNewFromInt64(t *testing.T) {
	for _, tc := range []struct {
		quantity int64
		unit     frequency.Unit
		expected string
	}{
		{0, frequency.HZ, "0 Hz"},
		{1, frequency.HZ, "1 Hz"},
		{2, frequency.KHZ, "2 kHz"},
		{4, frequency.MHZ, "4 MHz"},
		{8, frequency.GHZ, "8 GHz"},
		{16, frequency.THZ, "16 THz"},
	} {
		t.Run(tc.expected, func(t *testing.T) {
			v, err := frequency.NewFromInt64(tc.quantity, tc.unit)
			require.NoError(t, err)
			require.Equal(t, tc.expected, v.String())
		})
	}
}

func TestNewFromInt64_negative(t *testing.T) {
	_, err := frequency.NewFromInt64(-1, frequency.HZ)
	require.Error(t, err)
}

func TestNewFromBigInt(t *testing.T) {
	for _, tc := range []struct {
		quantity *big.Int
		unit     frequency.Unit
		expected string
	}{
		{big.NewInt(0), frequency.HZ, "0 Hz"},
		{big.NewInt(1), frequency.HZ, "1 Hz"},
		{big.NewInt(2), frequency.KHZ, "2 kHz"},
		{big.NewInt(4), frequency.MHZ, "4 MHz"},
		{big.NewInt(8), frequency.GHZ, "8 GHz"},
		{big.NewInt(16), frequency.THZ, "16 THz"},
	} {
		t.Run(tc.expected, func(t *testing.T) {
			v, err := frequency.NewFromBigInt(tc.quantity, tc.unit)
			require.NoError(t, err)
			require.Equal(t, tc.expected, v.String())
		})
	}
}

func TestNewFromBigInt_negative(t *testing.T) {
	_, err := frequency.NewFromBigInt(big.NewInt(-1), frequency.HZ)
	require.Error(t, err)
}

func TestParseString(t *testing.T) {
	for _, tc := range []struct {
		input    string
		quantity *big.Int
		unit     frequency.Unit
	}{
		{"0", big.NewInt(0), frequency.HZ},
		{"1", big.NewInt(1), frequency.HZ},
		{" 256  ", big.NewInt(256), frequency.HZ},
		{"1 hz", big.NewInt(1), frequency.HZ},
		{"   1   hz  ", big.NewInt(1), frequency.HZ},
		{"1 HZ", big.NewInt(1), frequency.HZ},
		{"100hz", big.NewInt(100), frequency.HZ},
		{"2 khz", big.NewInt(2000), frequency.KHZ},
		{"6.5 khz", big.NewInt(6500), frequency.KHZ},
		{"1.000000000 hz", big.NewInt(1), frequency.HZ},
		{"5 mhz", big.NewInt(5000000), frequency.MHZ},
		{"10 Ghz", big.NewInt(10000000000), frequency.GHZ},
		{"3 Thz", big.NewInt(3000000000000), frequency.THZ},
	} {
		t.Run(tc.input, func(t *testing.T) {
			v, err := frequency.ParseString(tc.input)
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
		"-5 khz",
		"100 k hz",
		"6000 000 hz",
		"0.5 hz",
		"5.0001 khz",
	} {
		t.Run(tc, func(t *testing.T) {
			_, err := frequency.ParseString(tc)
			require.Error(t, err)
		})
	}
}

func TestFrequency_MarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		value    frequency.Frequency
		expected string
	}{
		{frequency.MustNewFromInt64(1, frequency.GHZ), `"1 GHz"`},
		{frequency.MustNewFromInt64(100, frequency.KHZ), `"100 kHz"`},
		{frequency.MustParseString("1 khz"), `"1 khz"`},
		{frequency.MustParseString("1KHZ"), `"1KHZ"`},
	} {
		t.Run(tc.value.String(), func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(data))
		})
	}
}

func TestFrequency_UnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		input string
		value *big.Int
		unit  frequency.Unit
	}{
		{`"1 Ghz"`, big.NewInt(1000000000), frequency.GHZ},
		{`"100 hz"`, big.NewInt(100), frequency.HZ},
		{`"2 khz"`, big.NewInt(2000), frequency.KHZ},
		{`"1024KHZ"`, big.NewInt(1024000), frequency.KHZ},
	} {
		t.Run(tc.input, func(t *testing.T) {
			var v frequency.Frequency
			err := json.Unmarshal([]byte(tc.input), &v)
			require.NoError(t, err)
			require.Equal(t, tc.value, v.BigInt())
			require.Equal(t, strings.Trim(tc.input, `"`), *v.RawValue())
		})
	}
}

func TestFrequency_Clone(t *testing.T) {
	v := frequency.MustParseString("10 GHz")
	clone := v.Clone()

	require.Zero(t, v.BigInt().Cmp(clone.BigInt()))
	require.Equal(t, v.Unit(), clone.Unit())
	require.Equal(t, *v.RawValue(), *clone.RawValue())

	// Ensure they are separate instances
	require.NotSame(t, v.BigInt(), clone.BigInt())
	require.NotSame(t, v.RawValue(), clone.RawValue())
}

func TestFrequency_Equal(t *testing.T) {
	v1 := frequency.MustNewFromInt64(10, frequency.GHZ)
	v2 := frequency.MustNewFromInt64(10, frequency.GHZ)
	v3 := frequency.MustNewFromInt64(20, frequency.GHZ)
	v4 := frequency.MustNewFromInt64(10, frequency.KHZ)
	v5 := frequency.MustParseString("10 GHz")

	require.True(t, v1.Equal(v2))
	require.False(t, v1.Equal(v3))
	require.False(t, v1.Equal(v4))
	require.False(t, v1.Equal(v5)) // Different because v1 has no raw value
}

func TestFrequency_Cmp(t *testing.T) {
	v1 := frequency.MustNewFromInt64(10, frequency.GHZ)
	v2 := frequency.MustNewFromInt64(10, frequency.GHZ)
	v3 := frequency.MustNewFromInt64(20, frequency.GHZ)
	v4 := frequency.MustNewFromInt64(10, frequency.MHZ)
	v5 := frequency.MustNewFromInt64(5, frequency.KHZ)

	require.Equal(t, 0, v1.Cmp(v2))
	require.Equal(t, -1, v1.Cmp(v3))
	require.Equal(t, 1, v1.Cmp(v4))
	require.Equal(t, 1, v1.Cmp(v5))
}
