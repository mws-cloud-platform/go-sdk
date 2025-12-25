package cidraddress

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
)

func TestNewCIDR6Address(t *testing.T) {
	_, err := NewCIDR6Address(nil, nil)
	require.ErrorIs(t, err, ErrEmptyIPAddressOrIPNet)

	_, err = NewCIDR6Address(testIP4Address, testIP4Net)
	require.ErrorIs(t, err, ErrInvalidCIDRVersion)

	cidrAddr, err := NewCIDR6Address(testIP6Address, testIP6Net)
	require.NoError(t, err)
	cidrIP, cidrNet := cidrAddr.ToNetCIDR()
	require.Equal(t, testIP6Address, cidrIP)
	require.Equal(t, testIP6Net, cidrNet)
}

func TestCIDR6Address_Clone(t *testing.T) {
	cidrAddr, err := ParseCIDR6AddressString(cidrV6Raw)
	require.NoError(t, err)

	clone := cidrAddr.Clone()
	*cidrAddr.rawValue = "rawValue"
	cidrAddr.ip[0] = 0
	cidrAddr.ipNet.IP[0] = 0
	cidrAddr.ipNet.Mask[0] = 0

	require.NotEqual(t, cidrAddr.rawValue, clone.rawValue)
	require.NotEqual(t, cidrAddr.ip, clone.ip)
	require.NotEqual(t, cidrAddr.ipNet.IP, clone.ipNet.IP)
	require.NotEqual(t, cidrAddr.ipNet.Mask, clone.ipNet.Mask)
}

func TestCIDR6Address_Equal(t *testing.T) {
	cidrAddr, err := ParseCIDR6AddressString(cidrV6Raw)
	require.NoError(t, err)

	newAddr := cidrAddr

	require.True(t, cidrAddr.Equal(newAddr))
	require.False(t, cidrAddr.Equal(CIDR6Address{}))
}

func TestCIDR6Address_UnmarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/unmarshal_cidr6_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name        string
		rawValue    string
		errExpected bool
	}{
		{
			name:        "cidr4Address",
			rawValue:    "192.168.1.1/24",
			errExpected: true,
		},
		{
			name:     "cidr6Address",
			rawValue: "FE80::0202:B3FF:FE1E:8329/24",
		},
		{
			name:        "ipAddress",
			rawValue:    "192.168.1.1",
			errExpected: true,
		},
		{
			name:        "empty",
			rawValue:    "",
			errExpected: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var cidrAddr CIDR6Address
			err := json.Unmarshal([]byte(strconv.Quote(testCase.rawValue)), &cidrAddr)
			if testCase.errExpected {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NoError(t, fs.WriteFile(testCase.name+".txt", []byte(cidrAddr.String()), 0644))
		})
	}
}
