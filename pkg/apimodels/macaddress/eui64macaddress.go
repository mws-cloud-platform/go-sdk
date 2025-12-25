package macaddress

import (
	"fmt"
	"net"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

// EUI64MACAddress wrapper over the MACAddress type for additional parsing validations.
// It can contain only EUI64 mac address format.
type EUI64MACAddress struct {
	MACAddress
}

func NewEUI64MACAddress(mac net.HardwareAddr) (EUI64MACAddress, error) {
	macAddr, err := NewMACAddress(mac)
	if err != nil {
		return EUI64MACAddress{}, err
	}

	if !isEUI64Format(macAddr.mac) {
		return EUI64MACAddress{}, fmt.Errorf("%w: %d octets mac address are expected", ErrInvalidMACFormat, eui64octetsCount)
	}

	return EUI64MACAddress{
		MACAddress: macAddr,
	}, nil
}

// Clone returns a clone EUI64MACAddress with new pointer values
func (m *EUI64MACAddress) Clone() *EUI64MACAddress {
	if m == nil {
		return nil
	}

	return &EUI64MACAddress{
		MACAddress: ptr.Value(m.MACAddress.Clone()),
	}
}

// Equal checks if the values of m and m2 are equal
func (m EUI64MACAddress) Equal(m2 EUI64MACAddress) bool {
	return m.MACAddress.Equal(m2.MACAddress)
}

func (m *EUI64MACAddress) UnmarshalJSON(bytes []byte) error {
	return m.Decode(jx.DecodeBytes(bytes))
}

func (m *EUI64MACAddress) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseEUI64MACAddressString(rawValue)
	if err != nil {
		return err
	}

	m.mac = parsed.mac
	m.rawValue = parsed.rawValue
	return nil
}

// ParseEUI64MACAddressString parses a eui64 mac address string.
// Only eui64 format supported.
func ParseEUI64MACAddressString(s string) (EUI64MACAddress, error) {
	macAddress, err := ParseMACAddressString(s)
	if err != nil {
		return EUI64MACAddress{}, err
	}

	if !isEUI64Format(macAddress.mac) {
		return EUI64MACAddress{}, fmt.Errorf("%w: %w, %d octets mac address are expected", ErrInvalidMACString, ErrInvalidMACFormat, eui64octetsCount)
	}

	return EUI64MACAddress{
		MACAddress: macAddress,
	}, nil
}

func isEUI64Format(mac net.HardwareAddr) bool {
	return len(mac) == eui64octetsCount
}

// MustParseEUI64MACAddressString is like ParseEUI64MACAddressString but panics if the string cannot be parsed.
func MustParseEUI64MACAddressString(s string) EUI64MACAddress {
	result, err := ParseEUI64MACAddressString(s)
	if err != nil {
		panic(err)
	}
	return result
}
