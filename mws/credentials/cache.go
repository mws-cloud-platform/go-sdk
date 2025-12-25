package credentials

import (
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

const (
	// ErrEntryNotFound is returned when a cache entry is not found.
	ErrEntryNotFound = consterr.Error("cache entry not found")
	// ErrEntryRejected is returned when a cache entry is rejected.
	ErrEntryRejected = consterr.Error("cache entry rejected")
)

// Cache represents a cache for credentials.
type Cache interface {
	// Load retrieves a cached credentials entry by key.
	Load(string) (Credentials, error)
	// Store stores a credentials entry in the cache.
	Store(string, Credentials) error
	// Delete removes a credentials entry from the cache.
	Delete(string) error
}
