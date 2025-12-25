package duration_test

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/apimodels/units/duration"
)

func ExampleNewFromTimeDuration() {
	d := duration.NewFromTimeDuration(161*time.Hour + 6*time.Minute + 7*time.Second)
	fmt.Println(d)
	// Output: P6DT17H6M7S
}

func ExampleParseString_iso8601() {
	d, err := duration.ParseString("P6DT17H6M7S")
	if err != nil {
		panic(err)
	}

	fmt.Println(d)
	// Output: P6DT17H6M7S
}

func ExampleParseString_simple() {
	d, err := duration.ParseString("1h 30m")
	if err != nil {
		panic(err)
	}

	fmt.Println(d)
	// Output: PT1H30M
}

func ExampleParseString_seconds() {
	d, err := duration.ParseString("30")
	if err != nil {
		panic(err)
	}

	fmt.Println(d)
	// Output: PT30S
}

func ExampleDuration_UnmarshalJSON() {
	var d duration.Duration
	if err := json.Unmarshal([]byte(`"P6DT17H6M7S"`), &d); err != nil {
		panic(err)
	}

	fmt.Println(d)
	// Output: P6DT17H6M7S
}

func TestParseString(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected time.Duration
	}{
		// seconds
		{"0", 0},
		{"30", 30 * time.Second},
		{"3600", time.Hour},
		{"0.5", 500 * time.Millisecond},
		{"0.000000001", time.Nanosecond},
		{"+30", 30 * time.Second},
		{"-30", -30 * time.Second},
		{"-2562047h47m16.854775808s", time.Duration(math.MinInt64)},
		{"2562047h47m16.854775807s", time.Duration(math.MaxInt64)},

		// simple
		{"0s", 0},
		{"0m", 0},
		{"0h", 0},
		{"0d", 0},
		{"1ns", time.Nanosecond},
		{"1us", time.Microsecond},
		{"1ms", time.Millisecond},
		{"1s", time.Second},
		{"1m", time.Minute},
		{"1h", time.Hour},
		{"1d", 24 * time.Hour},
		{"1h30m", 1*time.Hour + 30*time.Minute},
		{"1h 30m", 1*time.Hour + 30*time.Minute},
		{"101d", 101 * 24 * time.Hour},
		{"39.5h", 39*time.Hour + 30*time.Minute},
		{"1h0m0.5s", time.Hour + 500*time.Millisecond},
		{"1H0M0.5S", time.Hour + 500*time.Millisecond},
		{"+(1h 30m)", 1*time.Hour + 30*time.Minute},
		{"-(1h 30m)", -1*time.Hour - 30*time.Minute},
		{"0.001s", time.Millisecond},
		{"0.000001s", time.Microsecond},
		{"0.000000001s", time.Nanosecond},
		{"0.0000000001s", 0},
		{"10.333333333333333s", 10*time.Second + 333*time.Millisecond + 333*time.Microsecond + 333*time.Nanosecond},
		{"-9223372036854775808ns", time.Duration(math.MinInt64)},
		{"-2562047h47m16.854775808s", time.Duration(math.MinInt64)},
		{"9223372036854775807ns", time.Duration(math.MaxInt64)},
		{"2562047h47m16.854775807s", time.Duration(math.MaxInt64)},

		// ISO-8601
		{"PT0S", 0}, {"P0D", 0}, {"PT0H", 0}, {"PT0M", 0}, {"PT0H0M0S", 0},
		{"PT1S", time.Second},
		{"PT2.5S", 2*time.Second + 500*time.Millisecond},
		{"PT1M", time.Minute},
		{"PT1H", time.Hour},
		{"P1D", 24 * time.Hour},
		{"PT1H1M", time.Hour + time.Minute},
		{"PT1H1M1S", time.Hour + time.Minute + time.Second},
		{"P6DT17H6M7S", 6*24*time.Hour + 17*time.Hour + 6*time.Minute + 7*time.Second},
		{"P6D T5H6M7S  ", 6*24*time.Hour + 5*time.Hour + 6*time.Minute + 7*time.Second},
		{"p6dt17h6m7s", 6*24*time.Hour + 17*time.Hour + 6*time.Minute + 7*time.Second},
		{"PT1H0M0.500S", time.Hour + 500*time.Millisecond},
		{"+PT1H30M", 1*time.Hour + 30*time.Minute},
		{"-PT5M", -5 * time.Minute},
		{"-PT1H30M", -1*time.Hour - 30*time.Minute},
		{"  -P6D T5H6M 7S ", -6*24*time.Hour - 5*time.Hour - 6*time.Minute - 7*time.Second},
	} {
		t.Run(tc.input, func(t *testing.T) {
			d, err := duration.ParseString(tc.input)
			require.NoError(t, err, tc.input)
			require.Equal(t, tc.expected, d.ToTimeDuration())
			require.Equal(t, tc.input, *d.RawValue())
		})
	}
}

func TestParseString_invalid(t *testing.T) {
	for _, tc := range []string{
		"",
		"    ",
		"-",
		"+",
		"()",
		"+()",
		"-()",
		"30m 1h",
		"PT",
		"P1DT",
		"P1DT10",
		"PT0SP0D",

		// overflow
		"-P106751DT23H47M16.854775809S",
		"P106751DT23H47M16.854775808S",
		"-9223372036.86",
		"9223372036.86",
		"-9223372037",
		"9223372037",
		"1" + strings.Repeat("0", 100) + "s",
	} {
		t.Run(tc, func(t *testing.T) {
			_, err := duration.ParseString(tc)
			require.Error(t, err, tc)
		})
	}
}

func TestDuration_String(t *testing.T) {
	for _, tc := range []struct {
		d        time.Duration
		expected string
	}{
		{0, "PT0S"},
		{time.Nanosecond, "PT0.000000001S"},
		{time.Microsecond, "PT0.000001S"},
		{time.Millisecond, "PT0.001S"},
		{500 * time.Millisecond, "PT0.5S"},
		{time.Second, "PT1S"},
		{time.Second + 500*time.Millisecond, "PT1.5S"},
		{time.Minute, "PT1M"},
		{time.Hour, "PT1H"},
		{24 * time.Hour, "P1D"},
		{24*time.Hour + time.Minute, "P1DT1M"},
		{48*time.Hour + 30*time.Minute, "P2DT30M"},
		{161*time.Hour + 6*time.Minute + 7*time.Second, "P6DT17H6M7S"},
		{-24 * time.Hour, "-P1D"},
		{-(time.Hour + 30*time.Minute), "-PT1H30M"},
		{time.Duration(math.MinInt64), "-P106751DT23H47M16.854775808S"},
		{time.Duration(math.MinInt64 + 1), "-P106751DT23H47M16.854775807S"},
		{time.Duration(math.MaxInt64), "P106751DT23H47M16.854775807S"},
	} {
		t.Run(tc.d.String(), func(t *testing.T) {
			actual := duration.NewFromTimeDuration(tc.d).String()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	t.Run("duration", func(t *testing.T) {
		expected := []byte(`"P6DT17H6M7S"`)
		d := duration.NewFromTimeDuration(161*time.Hour + 6*time.Minute + 7*time.Second)
		actual, err := json.Marshal(d)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("parsed", func(t *testing.T) {
		expected := []byte(`"1d 2h 3m 7s"`)
		d := duration.MustParseString("1d 2h 3m 7s")
		actual, err := json.Marshal(d)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}
