package macaddress

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
)

func TestNewEUI64MACAddress(t *testing.T) {
	_, err := NewEUI64MACAddress(nil)
	require.ErrorIs(t, err, ErrEmptyMACAddress)

	_, err = NewEUI64MACAddress(testMacEui48Address)
	require.ErrorIs(t, err, ErrInvalidMACFormat)

	ipAddr, err := NewEUI64MACAddress(testMacEui64Address)
	require.NoError(t, err)
	require.Equal(t, testMacEui64Address, ipAddr.ToNetHardwareAddr())
}

func TestEUI64MACAddress_Clone(t *testing.T) {
	macAddr, err := ParseEUI64MACAddressString(macEui64Raw)
	require.NoError(t, err)

	clone := macAddr.Clone()
	*macAddr.rawValue = "rawValue"
	macAddr.mac[0] = 0

	require.NotEqual(t, macAddr.rawValue, clone.rawValue)
	require.NotEqual(t, macAddr.mac, clone.mac)
}

func TestEUI64MACAddress_Equal(t *testing.T) {
	macAddr, err := ParseEUI64MACAddressString(macEui64Raw)
	require.NoError(t, err)

	newMac := macAddr

	require.True(t, macAddr.Equal(newMac))
	require.False(t, macAddr.Equal(EUI64MACAddress{}))
}

func TestEUI64MACAddress_UnmarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/unmarshal_eui64_mac_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name        string
		rawValue    string
		errExpected bool
	}{
		{
			name:        "eui48",
			rawValue:    "08:00:2b:01:02:03",
			errExpected: true,
		},
		{
			name:     "eui64",
			rawValue: "08:00:2b:01:02:03:04:05",
		},
		{
			name:        "empty",
			rawValue:    "",
			errExpected: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			rawJSON := strconv.Quote(testCase.rawValue)

			var macAddr EUI64MACAddress
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
