package models

import (
	"context"
	"maps"

	"github.com/go-faster/jx"

	"go.mws.cloud/go-sdk/internal/conv"
	"go.mws.cloud/go-sdk/pkg/context/values"
	resifaces "go.mws.cloud/go-sdk/pkg/resources/interfaces"
)

// NewAnyResourceRef constructor for creating an untyped resource id from its id.
func NewAnyResourceRef(ref string) AnyResourceRef {
	return AnyResourceRef{
		resource: ref,
	}
}

// NewAnyResourceRefErr is an error-returning variant of NewAnyResourceRef.
// Always succeeds - error is always nil.
func NewAnyResourceRefErr(ref string) (AnyResourceRef, error) {
	return NewAnyResourceRef(ref), nil
}

// ParseAnyResourceRef creates an AnyResourceRef from a path string and context.
func ParseAnyResourceRef(ctx context.Context, path string) (AnyResourceRef, error) {
	n := NewAnyResourceRef(path)
	if err := n.Parse(ctx); err != nil {
		return AnyResourceRef{}, err
	}
	return n, nil
}

// MustParseAnyResourceRef is like ParseAnyResourceRef but panics if an error has occurred.
// Intended for use in initialization code where failures are unrecoverable.
// Prefer ParseAnyResourceRef for runtime path processing.
func MustParseAnyResourceRef(ctx context.Context, path string) AnyResourceRef {
	result, err := ParseAnyResourceRef(ctx, path)
	if err != nil {
		panic(err)
	}
	return result
}

// AnyResourceRef a container for untyped resource reference that can be cast to a desired type using that type's constructor.
type AnyResourceRef struct {
	pathValues map[string]string
	resource   string
}

// GetPathValues returns path parameters extracted during parsing.
// These values are populated from the context when Parse is called.
func (n *AnyResourceRef) GetPathValues() map[string]string {
	if n.pathValues == nil {
		return nil
	}
	return n.pathValues
}

// Parse extracts path values from context and stores them in AnyResourceRef.
// These values are used later to convert this untyped reference to a typed one.
func (n *AnyResourceRef) Parse(ctx context.Context) error {
	if n == nil {
		return nil
	}

	n.pathValues = values.GetValuesStore(ctx)
	return nil
}

// ServiceSlug always returns "any" as this untyped reference may contain relative paths.
// This method provides interface compatibility. For the actual service slug,
// convert this untyped reference to a typed one.
func (n *AnyResourceRef) ServiceSlug() string {
	return anyServiceSlug
}

// IDPath returns the full resource path stored in this untyped reference.
func (n *AnyResourceRef) IDPath() string {
	if n == nil {
		return ""
	}

	return n.resource
}

// Path returns the same value as IDPath for interface compatibility.
func (n *AnyResourceRef) Path() string {
	return n.IDPath()
}

// String implements fmt.Stringer. Returns the same value as IDPath.
func (n *AnyResourceRef) String() string {
	return n.IDPath()
}

// Clone creates a deep copy of the AnyResourceRef.
// Returns nil if the receiver is nil.
func (n *AnyResourceRef) Clone() *AnyResourceRef {
	if n == nil {
		return nil
	}

	return &AnyResourceRef{
		resource: n.resource,
	}
}

// ResourceName returns the trailing segment of the resource path.
// Parses the stored path and extracts everything after the last "/".
// For empty paths or nil receiver, returns empty string.
func (n *AnyResourceRef) ResourceName() resifaces.ResourceName {
	if n == nil {
		return ""
	}

	return resourceNameFromResource(n.resource)
}

// Equal reports whether two resource refs are identical.
// Compares the underlying resource paths for exact string equality.
func (n AnyResourceRef) Equal(n2 AnyResourceRef) bool {
	return maps.Equal(n.pathValues, n2.pathValues) && n.resource == n2.resource
}

// MarshalJSON implements [json.Marshaler].
// Serializes AnyResourceRef as a JSON string.
func (n AnyResourceRef) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	if err := n.Encode(&e); err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}

// Encode implements [jx.Encoder].
// Writes AnyResourceRef as a JSON string.
func (n *AnyResourceRef) Encode(e *jx.Encoder) error {
	if n == nil {
		e.Null()
		return nil
	}

	e.Str(n.Path())
	return nil
}

// UnmarshalJSON implements [json.Unmarshaler].
// Parses JSON string into the AnyResourceRef.
func (n *AnyResourceRef) UnmarshalJSON(b []byte) error {
	return n.Decode(jx.DecodeBytes(b))
}

// Decode implements [jx.Decoder].
// Reads JSON string into the AnyResourceRef.
func (n *AnyResourceRef) Decode(d *jx.Decoder) error {
	if n == nil {
		return conv.NewDecodeToNilError("AnyResourceRef")
	}

	v, err := d.Str()
	if err != nil {
		return err
	}

	n.resource = v
	return nil
}
