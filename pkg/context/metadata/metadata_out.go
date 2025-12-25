package metadata

import (
	"context"
)

type keyOutgoing struct{}

// FromOutgoingContext returns outgoing metadata from the context.
func FromOutgoingContext(ctx context.Context) Metadata {
	return getMetadataCopy(ctx, keyOutgoing{})
}

// WithOutgoingMetadataValue returns context with added metadata values for the key
// replacing existing values for that key.
func WithOutgoingMetadataValue(ctx context.Context, key string, values ...string) context.Context {
	m := getMetadata(ctx, keyOutgoing{})
	if m == nil {
		m = make(Metadata)
		ctx = context.WithValue(ctx, keyOutgoing{}, m)
	}
	target := make([]string, len(values))
	copy(target, values)
	m[key] = target
	return ctx
}

// AppendOutgoingMetadataValue returns context with added metadata values for the key.
// New values will be appended to existing ones.
func AppendOutgoingMetadataValue(ctx context.Context, key string, values ...string) context.Context {
	m := getMetadata(ctx, keyOutgoing{})
	if m == nil {
		m = make(Metadata)
		ctx = context.WithValue(ctx, keyOutgoing{}, m)
	}
	mergeValues(m, key, values...)
	return ctx
}
