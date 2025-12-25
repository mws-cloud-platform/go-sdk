// Package cidraddress provides types and utilities for working with CIDR
// addresses.
package cidraddress

import (
	"fmt"
	"net"
	"slices"
	"strconv"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

// CIDRAddress wrapper over the standard net types describing CIDR to extend parsing functions.
// To get a standard representation, use the method ToNetCIDR.
// CIDRAddress can contain any ip version.
type CIDRAddress struct {
	ip       net.IP
	ipNet    *net.IPNet
	rawValue *string
}

func NewCIDRAddress(ip net.IP, ipNet *net.IPNet) (CIDRAddress, error) {
	if ip == nil || ipNet == nil {
		return CIDRAddress{}, ErrEmptyIPAddressOrIPNet
	}

	return CIDRAddress{
		ip:    ip,
		ipNet: ipNet,
	}, nil
}

// ToNetCIDR converts the CIDRAddress to the standard net types.
func (c CIDRAddress) ToNetCIDR() (net.IP, *net.IPNet) {
	return c.ip, c.ipNet
}

// RawValue returns a raw value if it was created from a string.
func (c CIDRAddress) RawValue() *string {
	return c.rawValue
}

// Clone returns a clone CIDRAddress with new pointer values
func (c *CIDRAddress) Clone() *CIDRAddress {
	if c == nil {
		return nil
	}

	clone := *c
	if c.ip != nil {
		clone.ip = make(net.IP, len(c.ip))
		copy(clone.ip, c.ip)
	}
	if c.ipNet != nil {
		clone.ipNet = new(net.IPNet)
		if c.ipNet.Mask != nil {
			clone.ipNet.Mask = make(net.IPMask, len(c.ipNet.Mask))
			copy(clone.ipNet.Mask, c.ipNet.Mask)
		}
		if c.ipNet.IP != nil {
			clone.ipNet.IP = make(net.IP, len(c.ipNet.IP))
			copy(clone.ipNet.IP, c.ipNet.IP)
		}
	}
	clone.rawValue = ptr.Clone(c.rawValue)
	return &clone
}

// Equal checks if the values of c and c2 are equal
func (c CIDRAddress) Equal(c2 CIDRAddress) bool {
	if (c.ipNet == nil) != (c2.ipNet == nil) {
		return false
	}

	if (c.ipNet != nil && c2.ipNet != nil) && (!slices.Equal(c.ipNet.Mask, c2.ipNet.Mask) || !slices.Equal(c.ipNet.IP, c2.ipNet.IP)) {
		return false
	}

	if !slices.Equal(c.ip, c2.ip) {
		return false
	}

	if !ptr.Equal(c.rawValue, c2.rawValue) {
		return false
	}

	return true
}

// String returns a string representing the IPAddress.
func (c CIDRAddress) String() string {
	if c.ip == nil || c.ipNet == nil {
		return ""
	}

	suffix, _ := c.ipNet.Mask.Size()
	return c.ip.String() + "/" + strconv.Itoa(suffix)
}

func (c CIDRAddress) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	c.Encode(&e)
	return e.Bytes(), nil
}

func (c *CIDRAddress) Encode(e *jx.Encoder) {
	if c == nil {
		e.Null()
		return
	}
	if c.rawValue != nil {
		e.Str(*c.rawValue)
	} else {
		e.Str(c.String())
	}
}

func (c *CIDRAddress) UnmarshalJSON(bytes []byte) error {
	return c.Decode(jx.DecodeBytes(bytes))
}

func (c *CIDRAddress) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseCIDRAddressString(rawValue)
	if err != nil {
		return err
	}

	c.ip = parsed.ip
	c.ipNet = parsed.ipNet
	c.rawValue = parsed.rawValue
	return nil
}

// ParseCIDRAddressString parses a CIDR address string.
// Inside uses the standard method net.ParseCIDR.
func ParseCIDRAddressString(s string) (CIDRAddress, error) {
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return CIDRAddress{}, fmt.Errorf("%w: %w", ErrInvalidCIDRString, err)
	}

	return CIDRAddress{
		ip:       ip,
		ipNet:    ipNet,
		rawValue: &s,
	}, nil
}

// MustParseCIDRAddressString is like ParseCIDRAddressString but panics if the string cannot be parsed.
func MustParseCIDRAddressString(s string) CIDRAddress {
	result, err := ParseCIDRAddressString(s)
	if err != nil {
		panic(err)
	}
	return result
}
