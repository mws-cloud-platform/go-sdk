package ipaddress

import (
	"fmt"
	"net"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

// IP4Address wrapper over the IPAddress type for additional parsing validations.
// It can contain only ipv4.
type IP4Address struct {
	IPAddress
}

func NewIP4Address(ip net.IP) (IP4Address, error) {
	ipAddr, err := NewIPAddress(ip)
	if err != nil {
		return IP4Address{}, err
	}

	if ipAddr.ip.To4() == nil {
		return IP4Address{}, fmt.Errorf("%w: ipv4 expected", ErrInvalidIPVersion)
	}

	return IP4Address{
		IPAddress: ipAddr,
	}, nil
}

// Clone returns a clone IP4Address with new pointer values
func (i *IP4Address) Clone() *IP4Address {
	if i == nil {
		return nil
	}

	return &IP4Address{
		IPAddress: ptr.Value(i.IPAddress.Clone()),
	}
}

// Equal checks if the values of i and i2 are equal
func (i IP4Address) Equal(i2 IP4Address) bool {
	return i.IPAddress.Equal(i2.IPAddress)
}

func (i *IP4Address) UnmarshalJSON(bytes []byte) error {
	return i.Decode(jx.DecodeBytes(bytes))
}

func (i *IP4Address) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseIP4AddressString(rawValue)
	if err != nil {
		return err
	}

	i.ip = parsed.ip
	i.rawValue = parsed.rawValue
	return nil
}

// ParseIP4AddressString parses an ipv4 address string.
// Only ipv4 supported: "192.168.1.1"
func ParseIP4AddressString(s string) (IP4Address, error) {
	ipAddress, err := ParseIPAddressString(s)
	if err != nil {
		return IP4Address{}, err
	}

	if ipAddress.ip.To4() == nil {
		return IP4Address{}, fmt.Errorf("%w: %w, ipv4 expected", ErrInvalidIPAddressString, ErrInvalidIPVersion)
	}

	return IP4Address{
		IPAddress: ipAddress,
	}, nil
}

// MustParseIP4AddressString is like ParseIP4AddressString but panics if the string cannot be parsed.
func MustParseIP4AddressString(s string) IP4Address {
	result, err := ParseIP4AddressString(s)
	if err != nil {
		panic(err)
	}
	return result
}
