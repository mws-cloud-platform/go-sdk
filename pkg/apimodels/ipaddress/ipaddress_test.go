package ipaddress

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
	ipv4Raw = "::ffff:192.0.2.1"
	ipv6Raw = "2001:0db0:0000:123a:0000:0000:0000:0030"
)

var (
	testIP4Address = net.ParseIP("192.0.2.1")
	testIP6Address = net.ParseIP("2001:db0:0:123a:0:0:0:30")
)

func TestNewIPAddress(t *testing.T) {
	_, err := NewIPAddress(nil)
	require.ErrorIs(t, err, ErrEmptyIPAddress)

	ipAddr, err := NewIPAddress(testIP4Address)
	require.NoError(t, err)
	require.Equal(t, testIP4Address, ipAddr.ToNetIP())

	ipAddr, err = NewIPAddress(testIP6Address)
	require.NoError(t, err)
	require.Equal(t, testIP6Address, ipAddr.ToNetIP())
}

func TestIPAddress_RawValue(t *testing.T) {
	ipAddr, err := NewIPAddress(testIP4Address)
	require.NoError(t, err)

	ipAddr.rawValue = ptr.Get(ipv4Raw)
	require.Equal(t, ipv4Raw, *ipAddr.RawValue())
}

func TestIPAddress_Clone(t *testing.T) {
	ipAddr, err := ParseIPAddressString(ipv6Raw)
	require.NoError(t, err)

	clone := ipAddr.Clone()
	*ipAddr.rawValue = "rawValue"
	ipAddr.ip[0] = 0

	require.NotEqual(t, ipAddr.rawValue, clone.rawValue)
	require.NotEqual(t, ipAddr.ip, clone.ip)
}

func TestIPAddress_Equal(t *testing.T) {
	ipAddr, err := ParseIPAddressString(ipv6Raw)
	require.NoError(t, err)

	for _, testCase := range []struct {
		name    string
		ipAddr1 IPAddress
		ipAddr2 IPAddress
		equal   bool
	}{
		{
			name:  "empty",
			equal: true,
		},
		{
			name:    "different rawValue 1",
			ipAddr1: IPAddress{},
			ipAddr2: IPAddress{rawValue: ptr.Get("")},
			equal:   false,
		},
		{
			name:    "different rawValue 2",
			ipAddr1: IPAddress{rawValue: ptr.Get("")},
			ipAddr2: IPAddress{},
			equal:   false,
		},
		{
			name:    "different rawValue 3",
			ipAddr1: IPAddress{rawValue: ptr.Get("hello")},
			ipAddr2: IPAddress{rawValue: ptr.Get("world")},
			equal:   false,
		},
		{
			name:    "different ip",
			ipAddr1: IPAddress{ip: net.ParseIP("192.168.1.1")},
			ipAddr2: IPAddress{ip: net.ParseIP("192.168.1.2")},
			equal:   false,
		},
		{
			name:    "equal ip",
			ipAddr1: IPAddress{ip: net.ParseIP("192.168.1.1")},
			ipAddr2: IPAddress{ip: net.ParseIP("192.168.1.1")},
			equal:   true,
		},
		{
			name:    "equal",
			ipAddr1: ipAddr,
			ipAddr2: ipAddr,
			equal:   true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.equal, testCase.ipAddr1.Equal(testCase.ipAddr2))
		})
	}
}

func TestIPAddress_String(t *testing.T) {
	ipAddr, err := ParseIPAddressString(ipv6Raw)
	require.NoError(t, err)
	require.Equal(t, testIP6Address.String(), ipAddr.String())

	ipAddr = IPAddress{}
	require.Empty(t, ipAddr.String())
}

func TestIPAddress_MarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/marshal_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name   string
		ipAddr IPAddress
	}{
		{
			name: "WithRawValue",
			ipAddr: IPAddress{
				ip:       testIP4Address,
				rawValue: ptr.Get(ipv4Raw),
			},
		},
		{
			name: "WithoutRawValueIPv4",
			ipAddr: IPAddress{
				ip: testIP4Address,
			},
		},
		{
			name: "WithoutRawValueIPv6",
			ipAddr: IPAddress{
				ip: testIP6Address,
			},
		},
		{
			name:   "WithNilAddress",
			ipAddr: IPAddress{},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := json.Marshal(testCase.ipAddr)
			require.NoError(t, err)

			require.NoError(t, fs.WriteFile(testCase.name+".txt", result, 0644))
		})
	}
}

func TestIPAddress_UnmarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/unmarshal_any_ip_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name        string
		rawValue    string
		errExpected bool
	}{
		{
			name:     "ipv4Address",
			rawValue: "192.168.1.1",
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
			var ipAddr IPAddress
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
