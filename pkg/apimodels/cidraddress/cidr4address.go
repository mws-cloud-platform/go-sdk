package cidraddress

import (
	"fmt"
	"net"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

// CIDR4Address wrapper over the CIDRAddress type for additional parsing validations.
// It can contain only ipv4.
type CIDR4Address struct {
	CIDRAddress
}

func NewCIDR4Address(ip net.IP, ipNet *net.IPNet) (CIDR4Address, error) {
	cidrAddr, err := NewCIDRAddress(ip, ipNet)
	if err != nil {
		return CIDR4Address{}, err
	}

	if cidrAddr.ip.To4() == nil {
		return CIDR4Address{}, fmt.Errorf("%w: ipv4 cidr expected", ErrInvalidCIDRVersion)
	}

	return CIDR4Address{
		cidrAddr,
	}, nil
}

// Clone returns a clone CIDR4Address with new pointer values
func (c *CIDR4Address) Clone() *CIDR4Address {
	if c == nil {
		return nil
	}

	return &CIDR4Address{
		CIDRAddress: ptr.Value(c.CIDRAddress.Clone()),
	}
}

// Equal checks if the values of c and c2 are equal
func (c CIDR4Address) Equal(c2 CIDR4Address) bool {
	return c.CIDRAddress.Equal(c2.CIDRAddress)
}

func (c *CIDR4Address) UnmarshalJSON(bytes []byte) error {
	return c.Decode(jx.DecodeBytes(bytes))
}

func (c *CIDR4Address) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseCIDR4AddressString(rawValue)
	if err != nil {
		return err
	}

	c.ip = parsed.ip
	c.ipNet = parsed.ipNet
	c.rawValue = parsed.rawValue
	return nil
}

// ParseCIDR4AddressString parses a CIDR4Address string.
// Only CIDR v4 supported: "192.0.2.0/24"
func ParseCIDR4AddressString(s string) (CIDR4Address, error) {
	cidrAddress, err := ParseCIDRAddressString(s)
	if err != nil {
		return CIDR4Address{}, err
	}

	if cidrAddress.ip.To4() == nil {
		return CIDR4Address{}, fmt.Errorf("%w: %w, ipv4 cidr expected", ErrInvalidCIDRString, ErrInvalidCIDRVersion)
	}

	return CIDR4Address{
		cidrAddress,
	}, nil
}

// MustParseCIDR4AddressString is like ParseCIDR4AddressString but panics if the string cannot be parsed.
func MustParseCIDR4AddressString(s string) CIDR4Address {
	result, err := ParseCIDR4AddressString(s)
	if err != nil {
		panic(err)
	}
	return result
}
