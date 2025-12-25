// Package metadata provides types and utilities for working with incoming and
// outgoing request metadata.
package metadata

import (
	"context"
	"fmt"
)

// Metadata represents maps with data for incoming and outgoing metadata.
// There are two sets of methods for working with metadata. Methods with 'Incoming'
// in the name works with incoming metadata, the same applies for 'Outgoing'.
// Maps for incoming and outgoing metadata are separate which means
// that the data for the same key will be different between incoming and outgoing metadata.
// When retrieving values, the data will be copied, so changing retrieved data
// will not affect stored values. Same with inserting, changing the source array
// after inserting will not affect inserted values.
type Metadata map[string][]string

func getMetadata(ctx context.Context, key any) Metadata {
	val := ctx.Value(key)
	if val == nil {
		return nil
	}
	casted, ok := val.(Metadata)
	if !ok {
		panic(fmt.Sprintf("Unexpected type '%T' for metadata holder", val))
	}
	return casted
}

func getMetadataCopy(ctx context.Context, key any) Metadata {
	m := getMetadata(ctx, key)
	if m == nil {
		return make(Metadata)
	}
	copied := make(Metadata)
	mergeMetadata(copied, m)

	return copied
}

func mergeMetadata(target, other Metadata) {
	for k, v := range other {
		mergeValues(target, k, v...)
	}
}

func mergeValues(m Metadata, key string, values ...string) {
	current, ok := m[key]
	var copied []string
	if ok {
		copied = make([]string, len(current)+len(values))
		copy(copied, append(current, values...))
	} else {
		copied = make([]string, len(values))
		copy(copied, values)
	}
	m[key] = copied
}
