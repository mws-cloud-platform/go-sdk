package ipaddress

import (
	"fmt"
	"net"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

// IP6Address wrapper over the IPAddress type for additional parsing validations.
// It can contain only ipv6.
type IP6Address struct {
	IPAddress
}

func NewIP6Address(ip net.IP) (IP6Address, error) {
	ipAddr, err := NewIPAddress(ip)
	if err != nil {
		return IP6Address{}, err
	}

	if ipAddr.ip.To4() != nil {
		return IP6Address{}, fmt.Errorf("%w: ipv6 expected", ErrInvalidIPVersion)
	}

	return IP6Address{
		IPAddress: ipAddr,
	}, nil
}

// Clone returns a clone IP6Address with new pointer values
func (i *IP6Address) Clone() *IP6Address {
	if i == nil {
		return nil
	}

	return &IP6Address{
		IPAddress: ptr.Value(i.IPAddress.Clone()),
	}
}

// Equal checks if the values of i and i2 are equal
func (i IP6Address) Equal(i2 IP6Address) bool {
	return i.IPAddress.Equal(i2.IPAddress)
}

func (i *IP6Address) UnmarshalJSON(bytes []byte) error {
	return i.Decode(jx.DecodeBytes(bytes))
}

func (i *IP6Address) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseIP6AddressString(rawValue)
	if err != nil {
		return err
	}

	i.ip = parsed.ip
	i.rawValue = parsed.rawValue
	return nil
}

// ParseIP6AddressString parses an ipv6 address string.
// Only ipv6 supported: "2001:db8::68"
func ParseIP6AddressString(s string) (IP6Address, error) {
	ipAddress, err := ParseIPAddressString(s)
	if err != nil {
		return IP6Address{}, err
	}

	if ipAddress.ip.To4() != nil {
		return IP6Address{}, fmt.Errorf("%w: %w, ipv6 expected", ErrInvalidIPAddressString, ErrInvalidIPVersion)
	}

	return IP6Address{
		IPAddress: ipAddress,
	}, nil
}

// MustParseIP6AddressString is like ParseIP6AddressString but panics if the string cannot be parsed.
func MustParseIP6AddressString(s string) IP6Address {
	result, err := ParseIP6AddressString(s)
	if err != nil {
		panic(err)
	}
	return result
}
