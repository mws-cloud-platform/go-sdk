package macaddress

import (
	"encoding/json"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

const (
	macEui48Raw = "08-00-2b-01-02-03"
	macEui64Raw = "08-00-2b-01-02-03-04-05"
)

var (
	testMacEui48Address, _ = net.ParseMAC("08:00:2b:01:02:03")
	testMacEui64Address, _ = net.ParseMAC("08:00:2b:01:02:03:04:05")
)

func TestNewMACAddress(t *testing.T) {
	_, err := NewMACAddress(nil)
	require.ErrorIs(t, err, ErrEmptyMACAddress)

	macAddr, err := NewMACAddress(testMacEui48Address)
	require.NoError(t, err)
	require.Equal(t, testMacEui48Address, macAddr.ToNetHardwareAddr())

	macAddr, err = NewMACAddress(testMacEui64Address)
	require.NoError(t, err)
	require.Equal(t, testMacEui64Address, macAddr.ToNetHardwareAddr())
}

func TestMACAddress_RawValue(t *testing.T) {
	macAddr, err := NewMACAddress(testMacEui48Address)
	require.NoError(t, err)

	macAddr.rawValue = ptr.Get(macEui48Raw)
	require.Equal(t, macEui48Raw, *macAddr.RawValue())
}

func TestMACAddress_Clone(t *testing.T) {
	macAddr, err := ParseMACAddressString(macEui64Raw)
	require.NoError(t, err)

	clone := macAddr.Clone()
	*macAddr.rawValue = "rawValue"
	macAddr.mac[0] = 0

	require.NotEqual(t, macAddr.rawValue, clone.rawValue)
	require.NotEqual(t, macAddr.mac, clone.mac)
}

func TestMACAddress_Equal(t *testing.T) {
	macAddr, err := ParseMACAddressString(macEui64Raw)
	require.NoError(t, err)

	for _, testCase := range []struct {
		name     string
		macAddr1 MACAddress
		macAddr2 MACAddress
		equal    bool
	}{
		{
			name:  "empty",
			equal: true,
		},
		{
			name:     "different rawValue 1",
			macAddr1: MACAddress{},
			macAddr2: MACAddress{rawValue: ptr.Get("")},
			equal:    false,
		},
		{
			name:     "different rawValue 2",
			macAddr1: MACAddress{rawValue: ptr.Get("")},
			macAddr2: MACAddress{},
			equal:    false,
		},
		{
			name:     "different rawValue 3",
			macAddr1: MACAddress{rawValue: ptr.Get("hello")},
			macAddr2: MACAddress{rawValue: ptr.Get("world")},
			equal:    false,
		},
		{
			name:     "different hw",
			macAddr1: MACAddress{mac: testMacEui48Address},
			macAddr2: MACAddress{mac: testMacEui64Address},
			equal:    false,
		},
		{
			name:     "equal hw",
			macAddr1: MACAddress{mac: testMacEui48Address},
			macAddr2: MACAddress{mac: testMacEui48Address},
			equal:    true,
		},
		{
			name:     "equal",
			macAddr1: macAddr,
			macAddr2: macAddr,
			equal:    true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.equal, testCase.macAddr1.Equal(testCase.macAddr2))
		})
	}
}

func TestMACAddress_String(t *testing.T) {
	macAddr, err := ParseMACAddressString(macEui64Raw)
	require.NoError(t, err)
	require.Equal(t, testMacEui64Address.String(), macAddr.String())

	macAddr = MACAddress{}
	require.Empty(t, macAddr.String())
}

func TestMACAddress_MarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/marshal_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name    string
		macAddr MACAddress
	}{
		{
			name: "WithRawValue",
			macAddr: MACAddress{
				mac:      testMacEui48Address,
				rawValue: ptr.Get(macEui48Raw),
			},
		},
		{
			name: "WithoutRawValueEUI48",
			macAddr: MACAddress{
				mac: testMacEui48Address,
			},
		},
		{
			name: "WithoutRawValueEUI64",
			macAddr: MACAddress{
				mac: testMacEui64Address,
			},
		},
		{
			name:    "WithNilAddress",
			macAddr: MACAddress{},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := json.Marshal(testCase.macAddr)
			require.NoError(t, err)

			require.NoError(t, fs.WriteFile(testCase.name+".txt", result, 0644))
		})
	}
}

func TestMACAddress_UnmarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/unmarshal_any_mac_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name        string
		rawValue    string
		errExpected bool
	}{
		{
			name:     "eui48",
			rawValue: "08:00:2b:01:02:03",
		},
		{
			name:     "eui48 mixed case",
			rawValue: "Ff:00:2b:01:02:03",
		},
		{
			name:     "eui48 with hyphens",
			rawValue: "08-00-2b-01-02-03",
		},
		{
			name:     "eui48 with one colon separator",
			rawValue: "08002b:010203",
		},
		{
			name:     "eui48 with one hyphen separator",
			rawValue: "08002b-010203",
		},
		{
			name:     "eui48 with double octets through a dot",
			rawValue: "0800.2b01.0203",
		},
		{
			name:     "eui48 with double octets through a hyphen",
			rawValue: "0800-2b01-0203",
		},
		{
			name:     "eui48 without separators",
			rawValue: "08002b010203",
		},
		{
			name:        "invalid eui48 case 1",
			rawValue:    "08:002b010203",
			errExpected: true,
		},
		{
			name:        "invalid eui48 case 2",
			rawValue:    "08-002b010203",
			errExpected: true,
		},
		{
			name:        "invalid eui48 case 3",
			rawValue:    "0800.2b010203",
			errExpected: true,
		},
		{
			name:        "invalid eui48 case 4",
			rawValue:    "0800-2b010203",
			errExpected: true,
		},
		{
			name:        "invalid eui48 case 5",
			rawValue:    "08002b:01020",
			errExpected: true,
		},
		{
			name:        "invalid eui48 case 6",
			rawValue:    "08002b01:020",
			errExpected: true,
		},
		{
			name:        "invalid eui48 case 7",
			rawValue:    "08002b01-020",
			errExpected: true,
		},
		{
			name:        "invalid eui48 case 8",
			rawValue:    "abcdefghijkl",
			errExpected: true,
		},
		{
			name:     "eui64",
			rawValue: "08:00:2b:01:02:03:04:05",
		},
		{
			name:     "eui64 mixed case",
			rawValue: "Ff:00:2b:01:02:03:04:05",
		},
		{
			name:     "eui64 with hyphens",
			rawValue: "08-00-2b-01-02-03-04-05",
		},
		{
			name:     "eui64 with one colon separator case 1",
			rawValue: "08002b:0102030405",
		},
		{
			name:     "eui64 with one colon separator case 2",
			rawValue: "08002b01:02030405",
		},
		{
			name:     "eui64 with one hyphen separator case 1",
			rawValue: "08002b-0102030405",
		},
		{
			name:     "eui64 with one hyphen separator case 2",
			rawValue: "08002b01-02030405",
		},
		{
			name:     "eui64 with double octets through a dot",
			rawValue: "0800.2b01.0203.0405",
		},
		{
			name:     "eui64 with double octets through a hyphen",
			rawValue: "0800-2b01-0203-0405",
		},
		{
			name:     "eui64 without separators",
			rawValue: "08002b0102030405",
		},
		{
			name:        "invalid eui64 case 1",
			rawValue:    "08:002b0102030405",
			errExpected: true,
		},
		{
			name:        "invalid eui48 case 2",
			rawValue:    "08-002b0102030405",
			errExpected: true,
		},
		{
			name:        "invalid eui64 case 3",
			rawValue:    "0800.2b0102030405",
			errExpected: true,
		},
		{
			name:        "invalid eui64 case 4",
			rawValue:    "0800-2b0102030405",
			errExpected: true,
		},
		{
			name:        "invalid eui64 case 5",
			rawValue:    "08002b:010203040",
			errExpected: true,
		},
		{
			name:        "invalid eui64 case 6",
			rawValue:    "08002b-010203040",
			errExpected: true,
		},
		{
			name:        "invalid eui64 case 7",
			rawValue:    "08002b01:0203040",
			errExpected: true,
		},
		{
			name:        "invalid eui64 case 8",
			rawValue:    "08002b01-0203040",
			errExpected: true,
		},
		{
			name:        "invalid eui64 case 9",
			rawValue:    "abcdefghijklmnop",
			errExpected: true,
		},
		{
			name:        "empty",
			rawValue:    "",
			errExpected: true,
		},
		{
			name:        "too long",
			rawValue:    "0000000000000000:0000000000000000000000000000000000000000000000000000",
			errExpected: true,
		},
		{
			name:        "too small",
			rawValue:    "00:00",
			errExpected: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			rawJSON := strconv.Quote(testCase.rawValue)

			var macAddr MACAddress
			err := json.Unmarshal([]byte(rawJSON), &macAddr)
			if testCase.errExpected {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NoError(t, fs.WriteFile(testCase.name+".txt",
				[]byte("RawJSON: "+rawJSON+"\nParsedMACString: "+macAddr.String()), 0644))
		})
	}
}
