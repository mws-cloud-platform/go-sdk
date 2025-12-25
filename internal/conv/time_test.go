package conv_test

import (
	"testing"
	"time"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/testing/golden"

	"go.mws.cloud/go-sdk/internal/conv"
	"go.mws.cloud/go-sdk/internal/decode"
)

func TestEncodeDateTime(t *testing.T) {
	for _, f := range []struct {
		name string
		enc  func(*jx.Encoder, time.Time)
	}{
		{
			name: "EncodeDateTime",
			enc:  conv.EncodeDateTime,
		},
		{
			name: "EncodeDateTimeUTC",
			enc:  conv.EncodeDateTimeUTC,
		},
	} {
		for _, tc := range []struct {
			name  string
			input time.Time
		}{
			{
				name:  "UTC timezone",
				input: time.Date(2023, 10, 15, 14, 30, 45, 123456789, time.UTC),
			},
			{
				name:  "negative timezone",
				input: time.Date(2023, 10, 15, 14, 30, 45, 123456789, time.FixedZone("EST", -5*3600)),
			},
			{
				name:  "positive timezone",
				input: time.Date(2023, 10, 15, 14, 30, 45, 123456789, time.FixedZone("CST", 8*3600)),
			},
			{
				name:  "leap second",
				input: time.Date(2016, 12, 31, 23, 59, 60, 0, time.UTC),
			},
			{
				name:  "zero time",
				input: time.Time{},
			},
			{
				name:  "max time",
				input: time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
			},
		} {
			t.Run(tc.name, func(tt *testing.T) {
				tt.Parallel()

				dir := golden.NewDir(tt,
					golden.WithPath("testdata/time/"+f.name+"/"+tc.name),
					golden.WithRecreateOnUpdate(),
				)

				encoder := jx.GetEncoder()
				defer jx.PutEncoder(encoder)

				f.enc(encoder, tc.input)
				result := encoder.String()
				dir.String(tt, "output.txt", result)

				actualTime, err := decode.DateTime(jx.DecodeBytes([]byte(result)))
				if err != nil {
					tt.Fatalf("Failed to parse actual result '%s': %v", result, err)
				}

				expectedTime := conv.DateTime(tc.input) // trim milliseconds
				if !actualTime.Equal(expectedTime) {
					tt.Errorf("Expected time %v, got %v", expectedTime, actualTime)
				}
			})
		}
	}
}
