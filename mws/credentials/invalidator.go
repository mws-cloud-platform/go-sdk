package credentials

import (
	"context"
)

// Invalidator represents a credentials provider that supports cached
// credentials invalidation.
//
// Commonly used with type assertion on [Provider]:
//
//	if invalidator, ok := provider.(Invalidator); ok {
//		err = invalidator.InvalidateCredentials(ctx)
//	}
type Invalidator interface {
	InvalidateCredentials(context.Context) error
}
