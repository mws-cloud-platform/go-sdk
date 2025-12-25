package bytesize_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
)

func ExampleNewFromInt64() {
	v, err := bytesize.NewFromInt64(256, bytesize.B)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 256 B
}

func ExampleNewFromBigInt() {
	v, err := bytesize.NewFromBigInt(big.NewInt(100), bytesize.KB)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 100 KB
}

func ExampleParseString() {
	v, err := bytesize.ParseString("20 GB")
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 20 GB
}

func ExampleByteSize_UnmarshalJSON() {
	var v bytesize.ByteSize
	if err := json.Unmarshal([]byte(`"256 MB"`), &v); err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output: 256 MB
}

func TestNewFromInt64(t *testing.T) {
	for _, tc := range []struct {
		quantity int64
		unit     bytesize.Unit
		expected string
	}{
		{0, bytesize.B, "0 B"},
		{1, bytesize.B, "1 B"},
		{2, bytesize.KB, "2 KB"},
		{4, bytesize.MB, "4 MB"},
		{8, bytesize.GB, "8 GB"},
		{16, bytesize.TB, "16 TB"},
		{32, bytesize.PB, "32 PB"},
		{64, bytesize.EB, "64 EB"},
		{128, bytesize.ZB, "128 ZB"},
		{256, bytesize.YB, "256 YB"},
		{512, bytesize.RB, "512 RB"},
		{1024, bytesize.QB, "1024 QB"},
	} {
		t.Run(tc.expected, func(t *testing.T) {
			v, err := bytesize.NewFromInt64(tc.quantity, tc.unit)
			require.NoError(t, err)
			require.Equal(t, tc.expected, v.String())
		})
	}
}

func TestNewFromInt64_negative(t *testing.T) {
	_, err := bytesize.NewFromInt64(-1, bytesize.B)
	require.Error(t, err)
}

func TestNewFromBigInt(t *testing.T) {
	for _, tc := range []struct {
		quantity *big.Int
		unit     bytesize.Unit
		expected string
	}{
		{big.NewInt(0), bytesize.B, "0 B"},
		{big.NewInt(1), bytesize.B, "1 B"},
		{big.NewInt(2), bytesize.KB, "2 KB"},
		{big.NewInt(4), bytesize.MB, "4 MB"},
		{big.NewInt(8), bytesize.GB, "8 GB"},
		{big.NewInt(16), bytesize.TB, "16 TB"},
		{big.NewInt(32), bytesize.PB, "32 PB"},
		{big.NewInt(64), bytesize.EB, "64 EB"},
		{big.NewInt(128), bytesize.ZB, "128 ZB"},
		{big.NewInt(256), bytesize.YB, "256 YB"},
		{big.NewInt(512), bytesize.RB, "512 RB"},
		{big.NewInt(1024), bytesize.QB, "1024 QB"},
	} {
		t.Run(tc.expected, func(t *testing.T) {
			v, err := bytesize.NewFromBigInt(tc.quantity, tc.unit)
			require.NoError(t, err)
			require.Equal(t, tc.expected, v.String())
		})
	}
}

func TestNewFromBigInt_negative(t *testing.T) {
	_, err := bytesize.NewFromBigInt(big.NewInt(-1), bytesize.B)
	require.Error(t, err)
}

func TestParseString(t *testing.T) {
	for _, tc := range []struct {
		input    string
		quantity *big.Int
		unit     bytesize.Unit
	}{
		{"0", big.NewInt(0), bytesize.B},
		{"1", big.NewInt(1), bytesize.B},
		{" 256  ", big.NewInt(256), bytesize.B},
		{"1 b", big.NewInt(1), bytesize.B},
		{"   1   b  ", big.NewInt(1), bytesize.B},
		{"1 B", big.NewInt(1), bytesize.B},
		{"100b", big.NewInt(100), bytesize.B},
		{"2 kb", big.NewInt(2 * 1 << 10), bytesize.KB},
		{"6.5 kb", big.NewInt(6656), bytesize.KB},
		{"1.000000000 b", big.NewInt(1), bytesize.B},
		{"5 mb", big.NewInt(5 * 1 << 20), bytesize.MB},
		{"10 Gb", big.NewInt(10 * 1 << 30), bytesize.GB},
		{"3 Tb", big.NewInt(3 * 1 << 40), bytesize.TB},
		{"1 QB", new(big.Int).Lsh(big.NewInt(1), 100), bytesize.QB},
	} {
		t.Run(tc.input, func(t *testing.T) {
			v, err := bytesize.ParseString(tc.input)
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
		"-5 kb",
		"100 k b",
		"6000 000 b",
		"0.5 b",
		"5.0001 kb",
	} {
		t.Run(tc, func(t *testing.T) {
			_, err := bytesize.ParseString(tc)
			require.Error(t, err)
		})
	}
}

func TestByteSize_MarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		value    bytesize.ByteSize
		expected string
	}{
		{bytesize.MustNewFromInt64(1, bytesize.GB), `"1 GB"`},
		{bytesize.MustNewFromInt64(100, bytesize.KB), `"100 KB"`},
		{bytesize.MustParseString("1 KB"), `"1 KB"`},
		{bytesize.MustParseString("1KB"), `"1KB"`},
	} {
		t.Run(tc.value.String(), func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(data))
		})
	}
}

func TestByteSize_UnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		input string
		value *big.Int
		unit  bytesize.Unit
	}{
		{`"1 Gb"`, big.NewInt(1 << 30), bytesize.GB},
		{`"100 b"`, big.NewInt(100), bytesize.KB},
		{`"2 kb"`, big.NewInt(2 * 1 << 10), bytesize.KB},
		{`"1024KB"`, big.NewInt(1024 * 1 << 10), bytesize.KB},
	} {
		t.Run(tc.input, func(t *testing.T) {
			var v bytesize.ByteSize
			err := json.Unmarshal([]byte(tc.input), &v)
			require.NoError(t, err)
			require.Equal(t, tc.value, v.BigInt())
			require.Equal(t, strings.Trim(tc.input, `"`), *v.RawValue())
		})
	}
}

func TestByteSize_Clone(t *testing.T) {
	v := bytesize.MustParseString("10 GB")
	clone := v.Clone()

	require.Zero(t, v.BigInt().Cmp(clone.BigInt()))
	require.Equal(t, v.Unit(), clone.Unit())
	require.Equal(t, *v.RawValue(), *clone.RawValue())

	// Ensure they are separate instances
	require.NotSame(t, v.BigInt(), clone.BigInt())
	require.NotSame(t, v.RawValue(), clone.RawValue())
}

func TestByteSize_Equal(t *testing.T) {
	v1 := bytesize.MustNewFromInt64(10, bytesize.GB)
	v2 := bytesize.MustNewFromInt64(10, bytesize.GB)
	v3 := bytesize.MustNewFromInt64(20, bytesize.GB)
	v4 := bytesize.MustNewFromInt64(10, bytesize.KB)
	v5 := bytesize.MustParseString("10 GB")

	require.True(t, v1.Equal(v2))
	require.False(t, v1.Equal(v3))
	require.False(t, v1.Equal(v4))
	require.False(t, v1.Equal(v5)) // Different because v1 has no raw value
}

func TestByteSize_Cmp(t *testing.T) {
	v1 := bytesize.MustNewFromInt64(10, bytesize.GB)
	v2 := bytesize.MustNewFromInt64(10, bytesize.GB)
	v3 := bytesize.MustNewFromInt64(20, bytesize.GB)
	v4 := bytesize.MustNewFromInt64(10, bytesize.MB)
	v5 := bytesize.MustNewFromInt64(5, bytesize.KB)

	require.Equal(t, 0, v1.Cmp(v2))
	require.Equal(t, -1, v1.Cmp(v3))
	require.Equal(t, 1, v1.Cmp(v4))
	require.Equal(t, 1, v1.Cmp(v5))
}
