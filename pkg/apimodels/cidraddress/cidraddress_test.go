package cidraddress

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
	cidrV4Raw = "::ffff:192.0.2.1/24"
	cidrV6Raw = "2001:0db0:0000:123a:0000:0000:0000:0030/32"
)

var (
	testIP4Address, testIP4Net, _ = net.ParseCIDR("192.0.2.1/24")
	testIP6Address, testIP6Net, _ = net.ParseCIDR("2001:db0:0:123a:0:0:0:30/32")
)

func TestNewCIDRAddress(t *testing.T) {
	_, err := NewCIDRAddress(nil, nil)
	require.ErrorIs(t, err, ErrEmptyIPAddressOrIPNet)

	cidrAddr, err := NewCIDRAddress(testIP4Address, testIP4Net)
	require.NoError(t, err)
	cidrIP, cidrNet := cidrAddr.ToNetCIDR()
	require.Equal(t, testIP4Address, cidrIP)
	require.Equal(t, testIP4Net, cidrNet)

	cidrAddr, err = NewCIDRAddress(testIP6Address, testIP6Net)
	require.NoError(t, err)
	cidrIP, cidrNet = cidrAddr.ToNetCIDR()
	require.Equal(t, testIP6Address, cidrIP)
	require.Equal(t, testIP6Net, cidrNet)
}

func TestCIDRAddress_RawValue(t *testing.T) {
	cidrAddr, err := NewCIDRAddress(testIP4Address, testIP4Net)
	require.NoError(t, err)

	cidrAddr.rawValue = ptr.Get(cidrV4Raw)
	require.Equal(t, cidrV4Raw, *cidrAddr.RawValue())
}

func TestCIDRAddress_Clone(t *testing.T) {
	cidrAddr, err := ParseCIDRAddressString(cidrV6Raw)
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

func TestCIDRAddress_Equal(t *testing.T) {
	cidrAddr, err := ParseCIDRAddressString(cidrV6Raw)
	require.NoError(t, err)

	for _, testCase := range []struct {
		name      string
		cidrAddr1 CIDRAddress
		cidrAddr2 CIDRAddress
		equal     bool
	}{
		{
			name:  "empty",
			equal: true,
		},
		{
			name:      "different rawValue 1",
			cidrAddr1: CIDRAddress{},
			cidrAddr2: CIDRAddress{rawValue: ptr.Get("")},
			equal:     false,
		},
		{
			name:      "different rawValue 2",
			cidrAddr1: CIDRAddress{rawValue: ptr.Get("")},
			cidrAddr2: CIDRAddress{},
			equal:     false,
		},
		{
			name:      "different rawValue 3",
			cidrAddr1: CIDRAddress{rawValue: ptr.Get("hello")},
			cidrAddr2: CIDRAddress{rawValue: ptr.Get("world")},
			equal:     false,
		},
		{
			name:      "different ip",
			cidrAddr1: CIDRAddress{ip: net.ParseIP("192.168.1.1")},
			cidrAddr2: CIDRAddress{ip: net.ParseIP("192.168.1.2")},
			equal:     false,
		},
		{
			name:      "equal ip",
			cidrAddr1: CIDRAddress{ip: net.ParseIP("192.168.1.1")},
			cidrAddr2: CIDRAddress{ip: net.ParseIP("192.168.1.1")},
			equal:     true,
		},
		{
			name:      "different ip net 1",
			cidrAddr1: CIDRAddress{},
			cidrAddr2: CIDRAddress{ipNet: ptr.Get(net.IPNet{})},
			equal:     false,
		},
		{
			name:      "different ip net 2",
			cidrAddr1: CIDRAddress{ipNet: ptr.Get(net.IPNet{})},
			cidrAddr2: CIDRAddress{},
			equal:     false,
		},
		{
			name:      "different ip net 3",
			cidrAddr1: CIDRAddress{ipNet: ptr.Get(net.IPNet{IP: net.ParseIP("192.168.1.1")})},
			cidrAddr2: CIDRAddress{ipNet: ptr.Get(net.IPNet{})},
			equal:     false,
		},
		{
			name:      "different ip net 4",
			cidrAddr1: CIDRAddress{ipNet: ptr.Get(net.IPNet{Mask: []byte{0, 1, 2, 3}})},
			cidrAddr2: CIDRAddress{ipNet: ptr.Get(net.IPNet{Mask: []byte{0, 1, 2}})},
			equal:     false,
		},
		{
			name:      "equal",
			cidrAddr1: cidrAddr,
			cidrAddr2: cidrAddr,
			equal:     true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.equal, testCase.cidrAddr1.Equal(testCase.cidrAddr2))
		})
	}
}

func TestCIDRAddress_String(t *testing.T) {
	cidrAddr, err := ParseCIDRAddressString(cidrV6Raw)
	require.NoError(t, err)
	require.Equal(t, "2001:db0:0:123a::30/32", cidrAddr.String())

	cidrAddr = CIDRAddress{}
	require.Empty(t, cidrAddr.String())
}

func TestCIDRAddress_MarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/marshal_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name     string
		cidrAddr CIDRAddress
	}{
		{
			name: "WithRawValue",
			cidrAddr: CIDRAddress{
				ip:       testIP4Address,
				ipNet:    testIP4Net,
				rawValue: ptr.Get(cidrV4Raw),
			},
		},
		{
			name: "WithoutRawValueCIDRv4",
			cidrAddr: CIDRAddress{
				ip:    testIP4Address,
				ipNet: testIP4Net,
			},
		},
		{
			name: "WithoutRawValueCIDRv6",
			cidrAddr: CIDRAddress{
				ip:    testIP6Address,
				ipNet: testIP6Net,
			},
		},
		{
			name:     "WithNilAddress",
			cidrAddr: CIDRAddress{},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := json.Marshal(testCase.cidrAddr)
			require.NoError(t, err)

			require.NoError(t, fs.WriteFile(testCase.name+".txt", result, 0644))
		})
	}
}

func TestCIDRAddress_UnmarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/unmarshal_any_cidr_json.golden"), golden.WithRecreateOnUpdate())
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
			var cidrAddr CIDRAddress
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
