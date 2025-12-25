package largenumber_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/apimodels/units/largenumber"
)

func ExampleNewFromInt64() {
	v, err := largenumber.NewFromInt64(256, largenumber.EmptyUnit)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 256
}

func ExampleNewFromBigInt() {
	v, err := largenumber.NewFromBigInt(big.NewInt(100), largenumber.K)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 100 K
}

func ExampleParseString() {
	v, err := largenumber.ParseString("20 K")
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 20 K
}

func ExampleLargeNumber_UnmarshalJSON() {
	var v largenumber.LargeNumber
	if err := json.Unmarshal([]byte(`"256 K"`), &v); err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 256 K
}

func TestNewFromInt64(t *testing.T) {
	for _, tc := range []struct {
		quantity int64
		unit     largenumber.Unit
		expected string
	}{
		{0, largenumber.EmptyUnit, "0"},
		{1, largenumber.EmptyUnit, "1"},
		{2, largenumber.K, "2 K"},
		{100, largenumber.K, "100 K"},
		{1000, largenumber.K, "1000 K"},
	} {
		t.Run(tc.expected, func(t *testing.T) {
			v, err := largenumber.NewFromInt64(tc.quantity, tc.unit)
			require.NoError(t, err)
			require.Equal(t, tc.expected, v.String())
		})
	}
}

func TestNewFromInt64_negative(t *testing.T) {
	_, err := largenumber.NewFromInt64(-1, largenumber.EmptyUnit)
	require.Error(t, err)
}

func TestNewFromBigInt(t *testing.T) {
	for _, tc := range []struct {
		quantity *big.Int
		unit     largenumber.Unit
		expected string
	}{
		{big.NewInt(0), largenumber.EmptyUnit, "0"},
		{big.NewInt(1), largenumber.EmptyUnit, "1"},
		{big.NewInt(2), largenumber.K, "2 K"},
		{big.NewInt(100), largenumber.K, "100 K"},
		{big.NewInt(1000), largenumber.K, "1000 K"},
	} {
		t.Run(tc.expected, func(t *testing.T) {
			v, err := largenumber.NewFromBigInt(tc.quantity, tc.unit)
			require.NoError(t, err)
			require.Equal(t, tc.expected, v.String())
		})
	}
}

func TestNewFromBigInt_negative(t *testing.T) {
	_, err := largenumber.NewFromBigInt(big.NewInt(-1), largenumber.EmptyUnit)
	require.Error(t, err)
}

func TestParseString(t *testing.T) {
	for _, tc := range []struct {
		input    string
		quantity *big.Int
		unit     largenumber.Unit
	}{
		{"0", big.NewInt(0), largenumber.EmptyUnit},
		{"1", big.NewInt(1), largenumber.EmptyUnit},
		{" 256  ", big.NewInt(256), largenumber.EmptyUnit},
		{"1k", big.NewInt(1000), largenumber.K},
		{"   1   k  ", big.NewInt(1000), largenumber.K},
		{"1 K", big.NewInt(1000), largenumber.K},
		{"100k", big.NewInt(100000), largenumber.K},
		{"2 k", big.NewInt(2000), largenumber.K},
		{"6.5 k", big.NewInt(6500), largenumber.K},
		{"1.000000000", big.NewInt(1), largenumber.EmptyUnit},
		{"-1", big.NewInt(-1), largenumber.EmptyUnit},
		{"-5 k", big.NewInt(-5000), largenumber.K},
	} {
		t.Run(tc.input, func(t *testing.T) {
			v, err := largenumber.ParseString(tc.input)
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
		"100 k k",
		"6000 000 k",
		"0.5",
		"5.0001 k",
	} {
		t.Run(tc, func(t *testing.T) {
			_, err := largenumber.ParseString(tc)
			require.Error(t, err)
		})
	}
}

func TestLargeNumber_MarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		value    largenumber.LargeNumber
		expected string
	}{
		{largenumber.MustNewFromInt64(1, largenumber.K), `"1 K"`},
		{largenumber.MustNewFromInt64(100, largenumber.EmptyUnit), `"100"`},
		{largenumber.MustParseString("1 k"), `"1 k"`},
		{largenumber.MustParseString("1K"), `"1K"`},
	} {
		t.Run(tc.value.String(), func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(data))
		})
	}
}

func TestLargeNumber_UnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		input string
		value *big.Int
		unit  largenumber.Unit
	}{
		{`"1 K"`, big.NewInt(1000), largenumber.K},
		{`"100"`, big.NewInt(100), largenumber.EmptyUnit},
		{`"2 k"`, big.NewInt(2000), largenumber.K},
		{`"1024K"`, big.NewInt(1024000), largenumber.K},
	} {
		t.Run(tc.input, func(t *testing.T) {
			var v largenumber.LargeNumber
			err := json.Unmarshal([]byte(tc.input), &v)
			require.NoError(t, err)
			require.Equal(t, tc.value, v.BigInt())
			require.Equal(t, strings.Trim(tc.input, `"`), *v.RawValue())
		})
	}
}

func TestLargeNumber_Clone(t *testing.T) {
	v := largenumber.MustParseString("10 K")
	clone := v.Clone()

	require.Zero(t, v.BigInt().Cmp(clone.BigInt()))
	require.Equal(t, v.Unit(), clone.Unit())
	require.Equal(t, *v.RawValue(), *clone.RawValue())

	// Ensure they are separate instances
	require.NotSame(t, v.BigInt(), clone.BigInt())
	require.NotSame(t, v.RawValue(), clone.RawValue())
}

func TestLargeNumber_Equal(t *testing.T) {
	v1 := largenumber.MustNewFromInt64(10, largenumber.K)
	v2 := largenumber.MustNewFromInt64(10, largenumber.K)
	v3 := largenumber.MustNewFromInt64(20, largenumber.K)
	v4 := largenumber.MustNewFromInt64(10, largenumber.EmptyUnit)
	v5 := largenumber.MustParseString("10 K")

	require.True(t, v1.Equal(v2))
	require.False(t, v1.Equal(v3))
	require.False(t, v1.Equal(v4))
	require.False(t, v1.Equal(v5)) // Different because v1 has no raw value
}

func TestLargeNumber_Cmp(t *testing.T) {
	v1 := largenumber.MustNewFromInt64(10, largenumber.K)
	v2 := largenumber.MustNewFromInt64(10, largenumber.K)
	v3 := largenumber.MustNewFromInt64(20, largenumber.K)
	v4 := largenumber.MustNewFromInt64(10, largenumber.EmptyUnit)
	v5 := largenumber.MustNewFromInt64(5, largenumber.EmptyUnit)

	require.Equal(t, 0, v1.Cmp(v2))
	require.Equal(t, -1, v1.Cmp(v3))
	require.Equal(t, 1, v1.Cmp(v4))
	require.Equal(t, 1, v1.Cmp(v5))
}
