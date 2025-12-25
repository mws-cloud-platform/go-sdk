package conv_test

import (
	"strings"
	"testing"

	"go.mws.cloud/util-toolset/pkg/testing/golden"

	"go.mws.cloud/go-sdk/internal/conv"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bitrate"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/frequency"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/largenumber"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/throughput"
)

func TestValueToString(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name  string
		kind  string
		value string
		enc   func() string
	}{
		// Bitrate
		{
			name:  "bits_bits",
			kind:  "bitrate",
			value: "12 123 bit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(12123, bitrate.BITS)
				return conv.BitrateToString(b, bitrate.BITS)
			},
		},
		{
			name:  "bits_gbits",
			kind:  "bitrate",
			value: "12 123 123 123 bit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(12123123123, bitrate.BITS)
				return conv.BitrateToString(b, bitrate.GBITS)
			},
		},
		{
			name:  "bits_gbits_out_kbits",
			kind:  "bitrate",
			value: "150 000 bit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(150000, bitrate.BITS)
				return conv.BitrateToString(b, bitrate.GBITS)
			},
		},
		{
			name:  "bits_gbits_out_mbits",
			kind:  "bitrate",
			value: "150 000 000 bit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(150000000, bitrate.BITS)
				return conv.BitrateToString(b, bitrate.GBITS)
			},
		},
		{
			name:  "bits_kbits",
			kind:  "bitrate",
			value: "12 123 123 123 bit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(12123123123, bitrate.BITS)
				return conv.BitrateToString(b, bitrate.KBITS)
			},
		},
		{
			name:  "bits_mbits",
			kind:  "bitrate",
			value: "12 123 123 123 bit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(12123123123, bitrate.BITS)
				return conv.BitrateToString(b, bitrate.MBITS)
			},
		},
		{
			name:  "bits_tbits",
			kind:  "bitrate",
			value: "12 123 123 123 123 123 bit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(12123123123123123, bitrate.BITS)
				return conv.BitrateToString(b, bitrate.TBITS)
			},
		},

		{
			name:  "kbits_gbits",
			kind:  "bitrate",
			value: "10 kbit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(10, bitrate.KBITS)
				return conv.BitrateToString(b, bitrate.GBITS)
			},
		},
		{
			name:  "mbits_gbits",
			kind:  "bitrate",
			value: "123 Mbit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(123, bitrate.MBITS)
				return conv.BitrateToString(b, bitrate.GBITS)
			},
		},
		{
			name:  "gbits_tbits",
			kind:  "bitrate",
			value: "10 Gbit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(10, bitrate.GBITS)
				return conv.BitrateToString(b, bitrate.TBITS)
			},
		},
		{
			name:  "gbits_kbits",
			kind:  "bitrate",
			value: "10 Gbit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(10, bitrate.GBITS)
				return conv.BitrateToString(b, bitrate.KBITS)
			},
		},
		{
			name:  "gbits_bits",
			kind:  "bitrate",
			value: "10 Gbit/s",
			enc: func() string {
				b := bitrate.MustNewFromInt64(10, bitrate.GBITS)
				return conv.BitrateToString(b, bitrate.BITS)
			},
		},
		{
			name:  "tbits_gbits",
			kind:  "bitrate",
			value: "123 Tbits",
			enc: func() string {
				b := bitrate.MustNewFromInt64(123, bitrate.TBITS)
				return conv.BitrateToString(b, bitrate.GBITS)
			},
		},
		{
			name:  "zero_gbits",
			kind:  "bitrate",
			value: "0 Gbit/s",
			enc: func() string {
				var b bitrate.Bitrate
				return conv.BitrateToString(b, bitrate.GBITS)
			},
		},
		{
			name:  "zero_bits",
			kind:  "bitrate",
			value: "0 bit/s",
			enc: func() string {
				var b bitrate.Bitrate
				return conv.BitrateToString(b, bitrate.BITS)
			},
		},

		// bytesize
		{
			name:  "gb_zb",
			kind:  "bytesize",
			value: "10 GB",
			enc: func() string {
				b := bytesize.MustNewFromInt64(10, bytesize.GB)
				return conv.ByteSizeToString(b, bytesize.ZB)
			},
		},
		{
			name:  "gb_gb",
			kind:  "bytesize",
			value: "10 GB",
			enc: func() string {
				b := bytesize.MustNewFromInt64(10, bytesize.GB)
				return conv.ByteSizeToString(b, bytesize.GB)
			},
		},
		{
			name:  "gb_b",
			kind:  "bytesize",
			value: "10 GB",
			enc: func() string {
				b := bytesize.MustNewFromInt64(10, bytesize.GB)
				return conv.ByteSizeToString(b, bytesize.B)
			},
		},
		{
			name:  "zero_b",
			kind:  "bytesize",
			value: "0 B",
			enc: func() string {
				var b bytesize.ByteSize
				return conv.ByteSizeToString(b, bytesize.B)
			},
		},

		// Frequency
		{
			name:  "ghz_thz",
			kind:  "frequency",
			value: "10 GHz",
			enc: func() string {
				f := frequency.MustNewFromInt64(10, frequency.GHZ)
				return conv.FrequencyToString(f, frequency.THZ)
			},
		},
		{
			name:  "ghz_mhz",
			kind:  "frequency",
			value: "10 GHz",
			enc: func() string {
				f := frequency.MustNewFromInt64(10, frequency.GHZ)
				return conv.FrequencyToString(f, frequency.MHZ)
			},
		},
		{
			name:  "ghz_hz",
			kind:  "frequency",
			value: "10 GHz",
			enc: func() string {
				f := frequency.MustNewFromInt64(10, frequency.GHZ)
				return conv.FrequencyToString(f, frequency.HZ)
			},
		},
		{
			name:  "zero_hz",
			kind:  "frequency",
			value: "0 Hz",
			enc: func() string {
				var f frequency.Frequency
				return conv.FrequencyToString(f, frequency.HZ)
			},
		},

		// LargeNumber (минимальная единица — пустая строка)
		{
			name:  "k_k",
			kind:  "largenumber",
			value: "10 K",
			enc: func() string {
				l := largenumber.MustNewFromInt64(10, largenumber.K)
				return conv.LargeNumberToString(l, largenumber.K)
			},
		},
		{
			name:  "k_empty",
			kind:  "largenumber",
			value: "10 K",
			enc: func() string {
				l := largenumber.MustNewFromInt64(10, largenumber.K)
				return conv.LargeNumberToString(l, largenumber.EmptyUnit)
			},
		},
		{
			name:  "zero_empty",
			kind:  "largenumber",
			value: "0",
			enc: func() string {
				var l largenumber.LargeNumber
				return conv.LargeNumberToString(l, largenumber.EmptyUnit)
			},
		},

		// Throughput
		{
			name:  "gbps_qbps",
			kind:  "throughput",
			value: "10 GBps",
			enc: func() string {
				th := throughput.MustNewFromInt64(10, throughput.GBPS)
				return conv.ThroughputToString(th, throughput.QBPS)
			},
		},
		{
			name:  "gbps_gbps",
			kind:  "throughput",
			value: "10 GBps",
			enc: func() string {
				th := throughput.MustNewFromInt64(10, throughput.GBPS)
				return conv.ThroughputToString(th, throughput.GBPS)
			},
		},
		{
			name:  "gbps_bps",
			kind:  "throughput",
			value: "10 GBps",
			enc: func() string {
				th := throughput.MustNewFromInt64(10, throughput.GBPS)
				return conv.ThroughputToString(th, throughput.BPS)
			},
		},
		{
			name:  "zero_bps",
			kind:  "throughput",
			value: "0 Bps",
			enc: func() string {
				var th throughput.Throughput
				return conv.ThroughputToString(th, throughput.BPS)
			},
		},
	} {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()

			dir := golden.NewDir(tt,
				golden.WithPath("testdata/encode/"+tc.kind+"/"+tc.name),
				golden.WithRecreateOnUpdate(),
			)

			var sb strings.Builder
			sb.WriteString(tc.value)
			sb.WriteString(" => ")
			sb.WriteString(tc.enc())
			dir.String(tt, "output.txt", sb.String())
		})
	}
}
