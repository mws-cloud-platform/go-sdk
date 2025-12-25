package cidraddress

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
)

func TestNewCIDR4Address(t *testing.T) {
	_, err := NewCIDR4Address(nil, nil)
	require.ErrorIs(t, err, ErrEmptyIPAddressOrIPNet)

	cidrAddr, err := NewCIDR4Address(testIP4Address, testIP4Net)
	require.NoError(t, err)
	cidrIP, cidrNet := cidrAddr.ToNetCIDR()
	require.Equal(t, testIP4Address, cidrIP)
	require.Equal(t, testIP4Net, cidrNet)

	_, err = NewCIDR4Address(testIP6Address, testIP6Net)
	require.ErrorIs(t, err, ErrInvalidCIDRVersion)
}

func TestCIDR4Address_Clone(t *testing.T) {
	cidrAddr, err := ParseCIDR4AddressString(cidrV4Raw)
	require.NoError(t, err)

	clone := cidrAddr.Clone()
	*cidrAddr.rawValue = "rawValue"
	cidrAddr.ip[0] = 1
	cidrAddr.ipNet.IP[0] = 1
	cidrAddr.ipNet.Mask[0] = 0

	require.NotEqual(t, cidrAddr.rawValue, clone.rawValue)
	require.NotEqual(t, cidrAddr.ip, clone.ip)
	require.NotEqual(t, cidrAddr.ipNet.IP, clone.ipNet.IP)
	require.NotEqual(t, cidrAddr.ipNet.Mask, clone.ipNet.Mask)
}

func TestCIDR4Address_Equal(t *testing.T) {
	cidrAddr, err := ParseCIDR4AddressString(cidrV4Raw)
	require.NoError(t, err)

	newAddr := cidrAddr

	require.True(t, cidrAddr.Equal(newAddr))
	require.False(t, cidrAddr.Equal(CIDR4Address{}))
}

func TestCIDR4Address_UnmarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/unmarshal_cidr4_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name        string
		rawValue    string
		errExpected bool
	}{
		{
			name:     "cidr4Address",
			rawValue: "192.168.1.1/24",
		},
		{
			name:        "cidr6Address",
			rawValue:    "FE80::0202:B3FF:FE1E:8329/24",
			errExpected: true,
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
			var cidrAddr CIDR4Address
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
