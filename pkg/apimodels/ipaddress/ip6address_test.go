package ipaddress

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
)

func TestNewIP6Address(t *testing.T) {
	_, err := NewIP6Address(nil)
	require.ErrorIs(t, err, ErrEmptyIPAddress)

	_, err = NewIP6Address(testIP4Address)
	require.ErrorIs(t, err, ErrInvalidIPVersion)

	ipAddr, err := NewIP6Address(testIP6Address)
	require.NoError(t, err)
	require.Equal(t, testIP6Address, ipAddr.ToNetIP())
}

func TestIP6Address_Clone(t *testing.T) {
	ipAddr, err := ParseIP6AddressString(ipv6Raw)
	require.NoError(t, err)

	clone := ipAddr.Clone()
	*ipAddr.rawValue = "rawValue"
	ipAddr.ip[0] = 0

	require.NotEqual(t, ipAddr.rawValue, clone.rawValue)
	require.NotEqual(t, ipAddr.ip, clone.ip)
}

func TestIP6Address_Equal(t *testing.T) {
	ipAddr, err := ParseIP6AddressString(ipv6Raw)
	require.NoError(t, err)

	newAddr := ipAddr

	require.True(t, ipAddr.Equal(newAddr))
	require.False(t, ipAddr.Equal(IP6Address{}))
}

func TestIP6Address_UnmarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/unmarshal_ipv6_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name        string
		rawValue    string
		errExpected bool
	}{
		{
			name:        "ipv4Address",
			rawValue:    "192.168.1.1",
			errExpected: true,
		},
		{
			name:     "ipv6Address",
			rawValue: "FE80::0202:B3FF:FE1E:8329",
		},
		{
			name:        "cidr",
			rawValue:    "192.168.1.1/24",
			errExpected: true,
		},
		{
			name:        "empty",
			rawValue:    "",
			errExpected: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var ipAddr IP6Address
			err := json.Unmarshal([]byte(strconv.Quote(testCase.rawValue)), &ipAddr)
			if testCase.errExpected {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NoError(t, fs.WriteFile(testCase.name+".txt", []byte(ipAddr.String()), 0644))
		})
	}
}
