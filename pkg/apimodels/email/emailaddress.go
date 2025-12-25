// Package email provides types for working with email addresses.
package email

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-faster/jx"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

type Email struct {
	username string
	domain   string
	rawValue *string
}

func (e Email) Username() string {
	return e.username
}

func (e Email) Domain() string {
	return e.domain
}

// RawValue returns a raw value if it was created from a string.
func (e Email) RawValue() *string {
	return e.rawValue
}

// Clone returns a clone Email with new pointer values
func (e *Email) Clone() *Email {
	if e == nil {
		return nil
	}

	clone := *e
	clone.rawValue = ptr.Clone(e.rawValue)
	return &clone
}

// Equal checks if the values of e and e2 are equal
func (e Email) Equal(e2 Email) bool {
	if e.username != e2.username || e.domain != e2.domain {
		return false
	}

	if !ptr.Equal(e.rawValue, e2.rawValue) {
		return false
	}

	return true
}

// String returns a string representing of the Email address.
func (e Email) String() string {
	if e.username == "" || e.domain == "" {
		return ""
	}
	return e.username + "@" + e.domain
}

func (e Email) MarshalJSON() ([]byte, error) {
	enc := jx.Encoder{}
	e.Encode(&enc)
	return enc.Bytes(), nil
}

func (e *Email) Encode(enc *jx.Encoder) {
	if e == nil {
		enc.Null()
		return
	}
	if e.rawValue != nil {
		enc.Str(*e.rawValue)
	} else {
		enc.Str(e.String())
	}
}

func (e *Email) UnmarshalJSON(bytes []byte) error {
	return e.Decode(jx.DecodeBytes(bytes))
}

func (e *Email) Decode(d *jx.Decoder) error {
	rawValue, err := d.Str()
	if err != nil {
		return err
	}

	parsed, err := ParseString(rawValue)
	if err != nil {
		return err
	}

	e.username = parsed.username
	e.domain = parsed.domain
	e.rawValue = parsed.rawValue
	return nil
}

var validationRegexp = regexp.MustCompile(`^(\S+)@(\S+)$`)

const (
	maxUsernameLen = 64
	maxDomainLen   = 255
	maxLen         = maxUsernameLen + maxDomainLen + 1

	ErrInvalidEmailString = consterr.Error("invalid email string")
)

// ParseString a method for very simple parsing of an email string according to a regular expression '^(\S+)@(\S+)$'
func ParseString(s string) (Email, error) {
	if len(s) > maxLen || !validationRegexp.MatchString(s) {
		return Email{}, ErrInvalidEmailString
	}

	username, domain, _ := strings.Cut(s, "@")
	if len(username) > maxUsernameLen {
		return Email{}, fmt.Errorf("%w: username is too long", ErrInvalidEmailString)
	}

	if len(domain) > maxDomainLen {
		return Email{}, fmt.Errorf("%w: domain is too long", ErrInvalidEmailString)
	}

	return Email{
		username: username,
		domain:   domain,
		rawValue: &s,
	}, nil
}

// MustParseString is like ParseString but panics if the string cannot be parsed.
func MustParseString(s string) Email {
	result, err := ParseString(s)
	if err != nil {
		panic(err)
	}
	return result
}
