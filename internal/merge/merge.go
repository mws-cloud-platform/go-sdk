package merge

import (
	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/pkg/apimodels/sensitive"
)

type GetNameFunc[T any] func(*T) string

type WithChangesFunc[T, U any] func(*T, U) T

// Slice returns a slice that is filled using this algorithm:
//  1. A map dm is created from dst using getNameT as key extractor
//  2. For every value u from src with its name not in dm, withChanges is called on empty object of type T and u, the result is added to out
//  3. For every value u from src with its name in dm, withChanges is called on dm[getNameU(&u)] and u, the result is added to out
func Slice[T, U any](
	dst []T, src []U,
	withChanges WithChangesFunc[T, U],
	getNameT GetNameFunc[T],
	getNameU GetNameFunc[U],
) []T {
	if len(src) == 0 {
		return nil
	}

	dm := make(map[string]T, len(dst))
	for _, t := range dst {
		dm[getNameT(&t)] = t
	}

	out := make([]T, 0, len(src))
	for _, u := range src {
		var t *T
		if m, ok := dm[getNameU(&u)]; ok {
			t = &m
		}
		out = append(out, withChanges(t, u))
	}

	return out
}

// SliceSensitive is equivalent to Slice with the exception that items in dst, src and return value are in sensitive.Sensitive containers
func SliceSensitive[T, U any](
	dst []sensitive.Sensitive[T], src []sensitive.Sensitive[U],
	withChanges WithChangesFunc[T, U],
	getNameT GetNameFunc[T],
	getNameU GetNameFunc[U],
) []sensitive.Sensitive[T] {
	if len(src) == 0 {
		return nil
	}

	dm := make(map[string]sensitive.Sensitive[T], len(dst))
	for _, t := range dst {
		dm[getNameT(ptr.Get(t.Value()))] = t
	}

	out := make([]sensitive.Sensitive[T], 0, len(src))
	for _, u := range src {
		var t *T
		if m, ok := dm[getNameU(ptr.Get(u.Value()))]; ok {
			t = ptr.Get(m.Value())
		}
		out = append(out, sensitive.New(withChanges(t, u.Value())))
	}

	return out
}

// Map returns a map that is filled using this algorithm:
//  1. For every key k that is in src but not in dst, withChanges is called on empty object of type T and src[k], the result is put to out[k]
//  2. For every key k that is both in src and dst, withChanges is called on dst[k] and src[k], the result is put to out[k]
func Map[T, U any](
	dst map[string]T, src map[string]U,
	withChanges WithChangesFunc[T, U],
) map[string]T {
	if len(src) == 0 {
		return nil
	}

	out := make(map[string]T)
	for k, v := range src {
		var d *T
		if dv, ok := dst[k]; ok {
			d = &dv
		}
		out[k] = withChanges(d, v)
	}
	return out
}

// MapSensitive is equivalent to Map with the exception that items in dst, src and return values are in sensitive.Sensitive containers
func MapSensitive[T, U any](
	dst map[string]sensitive.Sensitive[T], src map[string]sensitive.Sensitive[U],
	withChanges WithChangesFunc[T, U],
) map[string]sensitive.Sensitive[T] {
	if len(src) == 0 {
		return nil
	}

	out := make(map[string]sensitive.Sensitive[T])
	for k, v := range src {
		var t *T
		if dv, ok := dst[k]; ok {
			t = ptr.Get(dv.Value())
		}
		out[k] = sensitive.New(withChanges(t, v.Value()))
	}
	return out
}

// InapplicableSlice is needed for cases when elements of the slices don't have a required 'name' field.
// It simply returns a slice of every element from src combined with empty object of type T using withChanges.
func InapplicableSlice[T, U any](src []U, withChanges WithChangesFunc[T, U]) []T {
	if len(src) == 0 {
		return nil
	}

	out := make([]T, 0, len(src))
	for _, u := range src {
		out = append(out, withChanges(nil, u))
	}
	return out
}

// InapplicableSliceSensitive is equivalent to InapplicableSlice with the exception that items in src and return values are in sensitive.Sensitive containers
func InapplicableSliceSensitive[T, U any](src []sensitive.Sensitive[U], withChanges WithChangesFunc[T, U]) []sensitive.Sensitive[T] {
	if len(src) == 0 {
		return nil
	}

	out := make([]sensitive.Sensitive[T], 0, len(src))
	for _, u := range src {
		out = append(out, sensitive.New(withChanges(nil, u.Value())))
	}
	return out
}
