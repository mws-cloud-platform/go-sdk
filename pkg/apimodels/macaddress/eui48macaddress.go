package macaddress

import (
	"fmt"
	"net"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

// EUI48MACAddress wrapper over the MACAddress type for additional parsing validations.
// It can contain only EUI48 mac address format.
type EUI48MACAddress struct {
	MACAddress
}

func NewEUI48MACAddress(mac net.HardwareAddr) (EUI48MACAddress, error) {
	macAddr, err := NewMACAddress(mac)
	if err != nil {
		return EUI48MACAddress{}, err
	}

	if !isEUI48Format(macAddr.mac) {
		return EUI48MACAddress{}, fmt.Errorf("%w: %d octets mac address are expected", ErrInvalidMACFormat, eui48octetsCount)
	}

	return EUI48MACAddress{
		MACAddress: macAddr,
	}, nil
}

// Clone returns a clone EUI64MACAddress with new pointer values
func (m *EUI48MACAddress) Clone() *EUI48MACAddress {
	if m == nil {
		return nil
	}

	return &EUI48MACAddress{
		MACAddress: ptr.Value(m.MACAddress.Clone()),
	}
}

// Equal checks if the values of m and m2 are equal
func (m EUI48MACAddress) Equal(m2 EUI48MACAddress) bool {
	return m.MACAddress.Equal(m2.MACAddress)
}

func (m *EUI48MACAddress) UnmarshalJSON(bytes []byte) error {
	return m.Decode(jx.DecodeBytes(bytes))
}

func (m *EUI48MACAddress) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseEUI48MACAddressString(rawValue)
	if err != nil {
		return err
	}

	m.mac = parsed.mac
	m.rawValue = parsed.rawValue
	return nil
}

// ParseEUI48MACAddressString parses a eui48 mac address string.
// Only eui48 format supported.
func ParseEUI48MACAddressString(s string) (EUI48MACAddress, error) {
	macAddress, err := ParseMACAddressString(s)
	if err != nil {
		return EUI48MACAddress{}, err
	}

	if !isEUI48Format(macAddress.mac) {
		return EUI48MACAddress{}, fmt.Errorf("%w: %w, %d octets mac address are expected", ErrInvalidMACString, ErrInvalidMACFormat, eui48octetsCount)
	}

	return EUI48MACAddress{
		MACAddress: macAddress,
	}, nil
}

func isEUI48Format(mac net.HardwareAddr) bool {
	return len(mac) == eui48octetsCount
}

// MustParseEUI48MACAddressString is like ParseEUI48MACAddressString but panics if the string cannot be parsed.
func MustParseEUI48MACAddressString(s string) EUI48MACAddress {
	result, err := ParseEUI48MACAddressString(s)
	if err != nil {
		panic(err)
	}
	return result
}
