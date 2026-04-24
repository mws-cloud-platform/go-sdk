package models

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-faster/jx"

	"go.mws.cloud/go-sdk/internal/conv"
	resifaces "go.mws.cloud/go-sdk/pkg/resources/interfaces"
)

const anyServiceSlug = "any"

// NewAnyResourceID constructor for creating an untyped resource id from its id.
func NewAnyResourceID(name string) AnyResourceID {
	return AnyResourceID{
		resource: name,
	}
}

// ParseAnyResourceID constructor for creating an untyped resource id from its path.
func ParseAnyResourceID(path string) AnyResourceID {
	return NewAnyResourceID(path)
}

// AnyResourceID a container for untyped resource id that can be cast to a desired type using that type's constructor.
type AnyResourceID struct {
	resource string
}

// ServiceSlug returns the service identifier parsed from the resource path.
// For path-based IDs, this is the segment before the first "/".
// If the path has no separator or is malformed, returns "any".
func (n *AnyResourceID) ServiceSlug() string {
	if n == nil {
		return anyServiceSlug
	}
	pos := strings.IndexByte(n.resource, '/')
	if pos == -1 || pos == 0 {
		return anyServiceSlug
	}
	return n.resource[:pos]
}

// ID returns the full resource path as stored in the AnyResourceID.
func (n *AnyResourceID) ID() string {
	if n == nil {
		return ""
	}
	return n.resource
}

// String implements fmt.Stringer. Returns the same value as ID.
func (n *AnyResourceID) String() string {
	return n.ID()
}

// ResourceName returns the trailing segment of the resource path.
// Parses the stored path and extracts everything after the last "/".
// For empty paths or nil receiver, returns empty string.
func (n *AnyResourceID) ResourceName() resifaces.ResourceName {
	if n == nil {
		return ""
	}
	return resourceNameFromResource(n.resource)
}

// Parse is a no-op method that exists for interface compatibility.
// Does not perform any parsing or validation.
func (n *AnyResourceID) Parse(context.Context) error {
	return nil
}

// Clone creates a deep copy of the AnyResourceID.
// Returns nil if the receiver is nil.
func (n *AnyResourceID) Clone() *AnyResourceID {
	if n == nil {
		return nil
	}

	return &AnyResourceID{
		resource: n.resource,
	}
}

// Equal reports whether two resource IDs are identical.
// Compares the underlying resource paths for exact string equality.
func (n AnyResourceID) Equal(n2 AnyResourceID) bool {
	return n.resource == n2.resource
}

// MarshalJSON implements [json.Marshaler].
// Serializes AnyResourceID as a JSON string.
func (n AnyResourceID) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	if err := n.Encode(&e); err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}

// Encode implements [jx.Encoder].
// Writes AnyResourceID as a JSON string.
func (n *AnyResourceID) Encode(e *jx.Encoder) error {
	if n == nil {
		e.Null()
		return nil
	}

	e.Str(n.ID())
	return nil
}

// UnmarshalJSON implements [json.Unmarshaler].
// Parses JSON string into the AnyResourceID.
func (n *AnyResourceID) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &n.resource)
}

// Decode implements [jx.Decoder].
// Reads JSON string into the AnyResourceID.
func (n *AnyResourceID) Decode(d *jx.Decoder) error {
	if n == nil {
		return conv.NewDecodeToNilError("AnyResourceID")
	}

	v, err := d.Str()
	if err != nil {
		return err
	}

	n.resource = v
	return nil
}

func resourceNameFromResource(resource string) resifaces.ResourceName {
	return resifaces.ResourceName(resource[strings.LastIndex(resource, "/")+1:])
}
