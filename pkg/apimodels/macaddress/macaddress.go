// Package macaddress provides types and utilities for working with MAC
// addresses.
package macaddress

import (
	"net"
	"slices"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

// MACAddress wrapper over the standard net.HardwareAddr type to extend parsing functions.
// To get a standard representation of net.HardwareAddr, use the method ToNetHardwareAddr.
// MACAddress can contain eui-48 and eui-64 formats.
type MACAddress struct {
	mac      net.HardwareAddr
	rawValue *string
}

func NewMACAddress(mac net.HardwareAddr) (MACAddress, error) {
	if len(mac) == 0 {
		return MACAddress{}, ErrEmptyMACAddress
	}

	return MACAddress{
		mac: mac,
	}, nil
}

// ToNetHardwareAddr converts the MacAddress to the standard net type.
func (m MACAddress) ToNetHardwareAddr() net.HardwareAddr {
	return m.mac
}

// RawValue returns a raw value if it was created from a string.
func (m MACAddress) RawValue() *string {
	return m.rawValue
}

// Clone returns a clone MACAddress with new pointer values
func (m *MACAddress) Clone() *MACAddress {
	if m == nil {
		return nil
	}

	clone := *m
	if m.mac != nil {
		clone.mac = make(net.HardwareAddr, len(m.mac))
		copy(clone.mac, m.mac)
	}
	clone.rawValue = ptr.Clone(m.rawValue)
	return &clone
}

// Equal checks if the values of m and m2 are equal
func (m MACAddress) Equal(m2 MACAddress) bool {
	if !slices.Equal(m.mac, m2.mac) {
		return false
	}

	if !ptr.Equal(m.rawValue, m2.rawValue) {
		return false
	}

	return true
}

// String returns a string representing of the MAC address in canonical form with a separator ':' between octets
func (m MACAddress) String() string {
	return m.mac.String()
}

func (m MACAddress) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	m.Encode(&e)
	return e.Bytes(), nil
}

func (m *MACAddress) Encode(e *jx.Encoder) {
	if m == nil {
		e.Null()
		return
	}
	if m.rawValue != nil {
		e.Str(*m.rawValue)
	} else {
		e.Str(m.String())
	}
}

func (m *MACAddress) UnmarshalJSON(bytes []byte) error {
	return m.Decode(jx.DecodeBytes(bytes))
}

func (m *MACAddress) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseMACAddressString(rawValue)
	if err != nil {
		return err
	}

	m.mac = parsed.mac
	m.rawValue = parsed.rawValue
	return nil
}

// ParseMACAddressString parses a mac address string.
//
//	Valid mac addresses:
//	1. EUI48
//	    - "08:00:2b:01:02:03"
//	    - "08-00-2b-01-02-03"
//	    - "08002b:010203"
//	    - "08002b-010203"
//	    - "0800.2b01.0203"
//	    - "0800-2b01-0203"
//	    - "08002b010203"
//	2. EUI64
//	    - "08:00:2b:01:02:03:04:05"
//	    - "08-00-2b-01-02-03-04-05"
//	    - "08002b:0102030405"
//	    - "08002b-0102030405"
//	    - "0800.2b01.0203.0405"
//	    - "0800-2b01-0203-0405"
//	    - "08002b01:02030405"
//	    - "08002b0102030405"
func ParseMACAddressString(s string) (MACAddress, error) {
	var (
		l      = len(s)
		hwAddr net.HardwareAddr
		err    error
	)

	if l > 23 || l < 12 {
		return MACAddress{}, ErrInvalidMACString
	}

	switch {
	case s[2] == ':' || s[2] == '-':
		hwAddr, err = parseHWAddrOneHex(s)
	case s[4] == '.' || s[4] == '-':
		hwAddr, err = parseHWAddrTwoHex(s)
	case s[6] == ':' || s[6] == '-':
		hwAddr, err = parseHWAddrRowHex(s[:6] + s[7:])
	case s[8] == ':' || s[8] == '-':
		hwAddr, err = parseHWAddrRowHex(s[:8] + s[9:])
	case l == 12 || l == 16:
		hwAddr, err = parseHWAddrRowHex(s)
	default:
		return MACAddress{}, ErrInvalidMACString
	}
	if err != nil {
		return MACAddress{}, err
	}

	return MACAddress{
		mac:      hwAddr,
		rawValue: &s,
	}, nil
}

// MustParseMACAddressString is like ParseMACAddressString but panics if the string cannot be parsed.
func MustParseMACAddressString(s string) MACAddress {
	result, err := ParseMACAddressString(s)
	if err != nil {
		panic(err)
	}
	return result
}
