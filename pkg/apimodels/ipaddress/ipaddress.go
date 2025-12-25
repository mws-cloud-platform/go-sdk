// Package ipaddress provides types and utilities for working with IP addresses.
package ipaddress

import (
	"net"
	"slices"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

// IPAddress wrapper over the standard net.IP type to extend parsing functions.
// To get a standard representation of net.IP, use the method ToNetIP.
// IPAddress can contain any ip version.
type IPAddress struct {
	ip       net.IP
	rawValue *string
}

func NewIPAddress(ip net.IP) (IPAddress, error) {
	if ip == nil {
		return IPAddress{}, ErrEmptyIPAddress
	}

	return IPAddress{
		ip: ip,
	}, nil
}

// ToNetIP converts the IPAddress to the standard type net.IP.
func (i IPAddress) ToNetIP() net.IP {
	return i.ip
}

// RawValue returns a raw value if it was created from a string.
func (i IPAddress) RawValue() *string {
	return i.rawValue
}

// Clone returns a clone IPAddress with new pointer values
func (i *IPAddress) Clone() *IPAddress {
	if i == nil {
		return nil
	}

	clone := *i
	if i.ip != nil {
		clone.ip = make(net.IP, len(i.ip))
		copy(clone.ip, i.ip)
	}
	clone.rawValue = ptr.Clone(i.rawValue)
	return &clone
}

// Equal checks if the values of i and i2 are equal
func (i IPAddress) Equal(i2 IPAddress) bool {
	if !slices.Equal(i.ip, i2.ip) {
		return false
	}

	if !ptr.Equal(i.rawValue, i2.rawValue) {
		return false
	}

	return true
}

// String returns a string representing the IPAddress.
func (i IPAddress) String() string {
	if len(i.ip) == 0 {
		return ""
	}

	return i.ip.String()
}

func (i IPAddress) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	i.Encode(&e)
	return e.Bytes(), nil
}

func (i *IPAddress) Encode(e *jx.Encoder) {
	if i == nil {
		e.Null()
		return
	}
	if i.rawValue != nil {
		e.Str(*i.rawValue)
	} else {
		e.Str(i.String())
	}
}

func (i *IPAddress) UnmarshalJSON(bytes []byte) error {
	return i.Decode(jx.DecodeBytes(bytes))
}

func (i *IPAddress) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseIPAddressString(rawValue)
	if err != nil {
		return err
	}

	i.ip = parsed.ip
	i.rawValue = parsed.rawValue
	return nil
}

// ParseIPAddressString parses an ip address string.
// Inside uses the standard method net.ParseIP.
func ParseIPAddressString(s string) (IPAddress, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return IPAddress{}, ErrInvalidIPAddressString
	}

	return IPAddress{
		ip:       ip,
		rawValue: &s,
	}, nil
}

// MustParseIPAddressString is like ParseIPAddressString but panics if the string cannot be parsed.
func MustParseIPAddressString(s string) IPAddress {
	result, err := ParseIPAddressString(s)
	if err != nil {
		panic(err)
	}
	return result
}
