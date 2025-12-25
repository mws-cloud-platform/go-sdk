package macaddress

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
)

func TestNewEUI48MACAddress(t *testing.T) {
	_, err := NewEUI48MACAddress(nil)
	require.ErrorIs(t, err, ErrEmptyMACAddress)

	ipAddr, err := NewEUI48MACAddress(testMacEui48Address)
	require.NoError(t, err)
	require.Equal(t, testMacEui48Address, ipAddr.ToNetHardwareAddr())

	_, err = NewEUI48MACAddress(testMacEui64Address)
	require.ErrorIs(t, err, ErrInvalidMACFormat)
}

func TestEUI48MACAddress_Clone(t *testing.T) {
	macAddr, err := ParseEUI48MACAddressString(macEui48Raw)
	require.NoError(t, err)

	clone := macAddr.Clone()
	*macAddr.rawValue = "rawValue"
	macAddr.mac[0] = 0

	require.NotEqual(t, macAddr.rawValue, clone.rawValue)
	require.NotEqual(t, macAddr.mac, clone.mac)
}

func TestEUI48MACAddress_Equal(t *testing.T) {
	macAddr, err := ParseEUI48MACAddressString(macEui48Raw)
	require.NoError(t, err)

	newMac := macAddr

	require.True(t, macAddr.Equal(newMac))
	require.False(t, macAddr.Equal(EUI48MACAddress{}))
}

func TestEUI48MACAddress_UnmarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/unmarshal_eui48_mac_json.golden"), golden.WithRecreateOnUpdate())
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
			name:        "eui64",
			rawValue:    "08:00:2b:01:02:03:04:05",
			errExpected: true,
		},
		{
			name:        "empty",
			rawValue:    "",
			errExpected: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			rawJSON := strconv.Quote(testCase.rawValue)

			var macAddr EUI48MACAddress
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
