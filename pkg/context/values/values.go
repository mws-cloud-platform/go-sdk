// Package values provides utilities for managing key-value pairs in a context.
package values

import (
	"context"
	"fmt"
)

type keyValuesStore struct{}

// With creates a new context with the given key-value pair.
func With(ctx context.Context, key, value string) context.Context {
	store := GetValuesStore(ctx)
	if store == nil {
		store = make(map[string]string)
		ctx = WithValuesStore(ctx, store)
	}
	store[key] = value
	return ctx
}

// From retrieves the value associated with the given key from the context.
func From(ctx context.Context, key string) (string, bool) {
	store := GetValuesStore(ctx)
	if store == nil {
		return "", false
	}
	result, ok := store[key]
	return result, ok
}

// WithValuesStore creates a new context with the given key-value pairs.
func WithValuesStore(ctx context.Context, vs map[string]string) context.Context {
	return context.WithValue(ctx, keyValuesStore{}, vs)
}

// GetValuesStore retrieves the key-value store from the context.
func GetValuesStore(ctx context.Context) map[string]string {
	val := ctx.Value(keyValuesStore{})
	if val == nil {
		return nil
	}
	store, ok := val.(map[string]string)
	if !ok {
		panic(fmt.Sprintf("Unexpected type '%T' for values store", val))
	}
	return store
}
