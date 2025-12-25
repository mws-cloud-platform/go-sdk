package metadata

import (
	"context"
)

type keyIncoming struct{}

// FromIncomingContext returns outgoing metadata from the context.
func FromIncomingContext(ctx context.Context) Metadata {
	return getMetadataCopy(ctx, keyIncoming{})
}

// WithIncomingMetadata returns context with added incoming metadata.
// If metadata is already exists in the context it will be merged.
func WithIncomingMetadata(ctx context.Context, values Metadata) context.Context {
	copied := getMetadataCopy(ctx, keyIncoming{})
	mergeMetadata(copied, values)
	return context.WithValue(ctx, keyIncoming{}, copied)
}

// GetIncomingMetadataValue returns first value for the key from incoming metadata
// and a boolean flag indicating if the value is found.
func GetIncomingMetadataValue(ctx context.Context, key string) (string, bool) {
	m := getMetadata(ctx, keyIncoming{})
	v, ok := m[key]
	if ok && len(v) > 0 {
		return v[0], ok
	}
	return "", ok
}

// GetIncomingMetadataValues returns array of values for the key from incoming metadata
// and a boolean flag indicating if the value is found.
func GetIncomingMetadataValues(ctx context.Context, key string) ([]string, bool) {
	m := getMetadata(ctx, keyIncoming{})
	v, ok := m[key]
	if !ok {
		return nil, false
	}
	return append([]string{}, v...), ok
}
