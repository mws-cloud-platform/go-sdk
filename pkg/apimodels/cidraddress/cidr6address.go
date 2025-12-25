package cidraddress

import (
	"fmt"
	"net"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

// CIDR6Address wrapper over the CIDRAddress type for additional parsing validations.
// It can contain only ipv6.
type CIDR6Address struct {
	CIDRAddress
}

func NewCIDR6Address(ip net.IP, ipNet *net.IPNet) (CIDR6Address, error) {
	cidrAddr, err := NewCIDRAddress(ip, ipNet)
	if err != nil {
		return CIDR6Address{}, err
	}

	if cidrAddr.ip.To4() != nil {
		return CIDR6Address{}, fmt.Errorf("%w: ipv6 cidr expected", ErrInvalidCIDRVersion)
	}

	return CIDR6Address{
		cidrAddr,
	}, nil
}

// Clone returns a clone CIDR6Address with new pointer values
func (c *CIDR6Address) Clone() *CIDR6Address {
	if c == nil {
		return nil
	}

	return &CIDR6Address{
		CIDRAddress: ptr.Value(c.CIDRAddress.Clone()),
	}
}

// Equal checks if the values of c and c2 are equal
func (c CIDR6Address) Equal(c2 CIDR6Address) bool {
	return c.CIDRAddress.Equal(c2.CIDRAddress)
}

func (c *CIDR6Address) UnmarshalJSON(bytes []byte) error {
	return c.Decode(jx.DecodeBytes(bytes))
}

func (c *CIDR6Address) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseCIDR6AddressString(rawValue)
	if err != nil {
		return err
	}

	c.ip = parsed.ip
	c.ipNet = parsed.ipNet
	c.rawValue = parsed.rawValue
	return nil
}

// ParseCIDR6AddressString parses a CIDR6Address string.
// Only CIDR v6 supported: "2001:db8::/32"
func ParseCIDR6AddressString(s string) (CIDR6Address, error) {
	cidrAddress, err := ParseCIDRAddressString(s)
	if err != nil {
		return CIDR6Address{}, err
	}

	if cidrAddress.ip.To4() != nil {
		return CIDR6Address{}, fmt.Errorf("%w: %w, ipv6 cidr expected", ErrInvalidCIDRString, ErrInvalidCIDRVersion)
	}

	return CIDR6Address{
		cidrAddress,
	}, nil
}

// MustParseCIDR6AddressString is like ParseCIDR6AddressString but panics if the string cannot be parsed.
func MustParseCIDR6AddressString(s string) CIDR6Address {
	result, err := ParseCIDR6AddressString(s)
	if err != nil {
		panic(err)
	}
	return result
}
