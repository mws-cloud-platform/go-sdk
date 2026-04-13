package macaddress

import "net"

const (
	lenOneHexWithSeparator = 3
	lenTwoHexWithSeparator = 5
	countBytesInTwoHex     = 2
	countCharsInHex        = 2

	// Bigger than we need, not too big to worry about overflow
	big = 0xFFFFFF

	hexOffset = 10
)

func parseHWAddrOneHex(s string) (net.HardwareAddr, error) {
	n := (len(s) + 1) / lenOneHexWithSeparator
	if n != 6 && n != 8 && n != 20 {
		return nil, ErrInvalidMACString
	}
	hw := make(net.HardwareAddr, n)
	for x, i := 0, 0; i < n; i++ {
		var ok bool
		if hw[i], ok = xtoi2(s[x:], s[2]); !ok {
			return nil, ErrInvalidMACString
		}
		x += lenOneHexWithSeparator
	}
	return hw, nil
}

func parseHWAddrTwoHex(s string) (net.HardwareAddr, error) {
	n := countBytesInTwoHex * (len(s) + 1) / lenTwoHexWithSeparator
	if n != 6 && n != 8 && n != 20 {
		return nil, ErrInvalidMACString
	}
	hw := make(net.HardwareAddr, n)
	for x, i := 0, 0; i < n; i += 2 {
		var ok bool
		if hw[i], ok = xtoi2(s[x:x+countCharsInHex], 0); !ok {
			return nil, ErrInvalidMACString
		}
		if hw[i+1], ok = xtoi2(s[x+countCharsInHex:], s[4]); !ok {
			return nil, ErrInvalidMACString
		}
		x += lenTwoHexWithSeparator
	}
	return hw, nil
}

func parseHWAddrRowHex(s string) (net.HardwareAddr, error) {
	n := len(s) / countCharsInHex
	if n != 6 && n != 8 && n != 20 {
		return nil, ErrInvalidMACString
	}
	hw := make(net.HardwareAddr, n)
	for x, i := 0, 0; i < n; i++ {
		var ok bool
		if hw[i], ok = xtoi2(s[x:x+countCharsInHex], 0); !ok {
			return nil, ErrInvalidMACString
		}
		x += countCharsInHex
	}
	return hw, nil
}

// Hexadecimal to integer.
// Returns number, characters consumed, success.
func xtoi(s string) (n int, i int, ok bool) {
	n = 0
LOOP:
	for i = 0; i < len(s); i++ {
		switch {
		case '0' <= s[i] && s[i] <= '9':
			n *= 16
			n += int(s[i] - '0')
		case 'a' <= s[i] && s[i] <= 'f':
			n *= 16
			n += int(s[i]-'a') + hexOffset
		case 'A' <= s[i] && s[i] <= 'F':
			n *= 16
			n += int(s[i]-'A') + hexOffset
		default:
			break LOOP
		}
		if n >= big {
			return 0, i, false
		}
	}
	if i == 0 {
		return 0, i, false
	}
	return n, i, true
}

// xtoi2 converts the next two hex digits of s into a byte.
// If s is longer than 2 bytes then the third byte must be e.
// If the first two bytes of s are not hex digits or the third byte
// does not match e, false is returned.
func xtoi2(s string, e byte) (byte, bool) {
	if len(s) > 2 && s[2] != e {
		return 0, false
	}
	n, ei, ok := xtoi(s[:2])
	return byte(n), ok && ei == 2 //nolint:gosec // xtoi always returns bytes
}
