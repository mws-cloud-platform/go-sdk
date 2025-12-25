package client

import (
	"encoding/json"
	"slices"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

func DiffPrimitiveRequired[T comparable](from, to T, nilDiffers bool) Optional[T] {
	return NewDirectOptional[T](to, nilDiffers || from != to)
}

func DiffPrimitiveNonRequired[T comparable](from, to *T, nilDiffers bool) Optional[T] {
	var value T
	if to != nil {
		value = *to
	}
	hasChanges := nilDiffers || (to == nil && from != nil) || (to != nil && from == nil)
	if !hasChanges && to != nil && from != nil {
		hasChanges = *to != *from
	}

	return NewDirectOptional[T](value, hasChanges)
}

func DiffPrimitiveNullable[T comparable](from, to *T, nilDiffers bool) OptionalNil[T] {
	res := DiffPrimitiveNonRequired(from, to, nilDiffers)
	return NewDirectOptionalNil[T](res.Value, res.Set, to == nil)
}

type Equatable[T any] interface {
	Equal(other T) bool
}

func DiffEquatableIfaceRequired[T Equatable[T]](from, to T, nilDiffers bool) Optional[T] {
	return NewDirectOptional[T](to, nilDiffers || !from.Equal(to))
}

func DiffEquatableIfaceNonRequired[T Equatable[T]](from, to *T, nilDiffers bool) Optional[T] {
	var value T
	if to != nil {
		value = *to
	}
	hasChanges := nilDiffers || (to == nil && from != nil) || (to != nil && from == nil)
	if !hasChanges && to != nil && from != nil {
		valTo := *to
		valFrom := *from
		hasChanges = !valTo.Equal(valFrom)
	}

	return NewDirectOptional[T](value, hasChanges)
}

func DiffEquatableIfaceNullable[T Equatable[T]](from, to *T, nilDiffers bool) OptionalNil[T] {
	res := DiffEquatableIfaceNonRequired(from, to, nilDiffers)
	return NewDirectOptionalNil[T](res.Value, res.Set, to == nil)
}

func DiffRawData[T ~[]byte](from, to T) Optional[T] {
	return NewDirectOptional[T](to, !slices.Equal(from, to))
}

func DiffRawDataNullable[T ~[]byte](from, to T) OptionalNil[T] {
	return NewDirectOptionalNil(to, !slices.Equal(from, to), to == nil)
}

func GetChangesArrayPrimitive[T comparable](from, to []T) ([]T, bool) {
	hasChanges := len(from) != len(to)
	if !hasChanges {
		for i, toEl := range to {
			if toEl != from[i] {
				hasChanges = true
				break
			}
		}
	}

	return to, hasChanges
}

func GetChangesArrayRawData(from, to []json.RawMessage) ([]json.RawMessage, bool) {
	hasChanges := len(from) != len(to)
	if !hasChanges {
		for i, toEl := range to {
			if !slices.Equal(from[i], toEl) {
				hasChanges = true
				break
			}
		}
	}

	return to, hasChanges
}

func GetChangesArrayEquatableIface[T Equatable[T]](from, to []T) ([]T, bool) {
	hasChanges := len(from) != len(to)
	if !hasChanges {
		for i, toEl := range to {
			if !toEl.Equal(from[i]) {
				hasChanges = true
				break
			}
		}
	}

	return to, hasChanges
}

func GetChangesArrayObject[T any, Tu updateObject](from, to []T, getDiff func(T, T, bool) Tu) ([]Tu, bool) {
	value := make([]Tu, 0, len(from))
	hasChanges := len(to) != len(from)
	if !hasChanges {
		for i, toItem := range to {
			value = append(value, getDiff(from[i], toItem, false))
		}
	} else {
		for _, toItem := range to {
			value = append(value, getDiff(toItem, toItem, true))
		}
	}

	if !hasChanges {
		for _, item := range value {
			if item.HasChanges() {
				hasChanges = true
				break
			}
		}
	}
	return value, hasChanges
}

func GetChangesArrayObjectError[T any, Tu updateObject](from, to []T, getDiff func(T, T, bool) (Tu, error)) ([]Tu, bool, error) {
	value := make([]Tu, 0, len(from))
	hasChanges := len(to) != len(from)
	if !hasChanges {
		for i, toItem := range to {
			diffValue, err := getDiff(from[i], toItem, false)
			if err != nil {
				return nil, false, err
			}
			value = append(value, diffValue)
		}
	} else {
		for _, toItem := range to {
			diffValue, err := getDiff(toItem, toItem, true)
			if err != nil {
				return nil, false, err
			}
			value = append(value, diffValue)
		}
	}

	if !hasChanges {
		for _, item := range value {
			if item.HasChanges() {
				hasChanges = true
				break
			}
		}
	}
	return value, hasChanges, nil
}

type namedChild interface {
	GetName() string
}

type updateObject interface {
	HasChanges() bool
}

type nameSetter[T any] interface {
	SetName(string)
	*T
}

func ToPointerArray[T any](arr []T) []*T {
	if arr == nil {
		return nil
	}

	res := make([]*T, 0, len(arr))
	for _, v := range arr {
		res = append(res, ptr.Get(v))
	}
	return res
}

const ErrDuplicateKeys = consterr.Error("cannot determine changes, slices have duplicate name keys")

func GetChangesArrayObjectNamed[T namedChild, Tu updateObject, Tns nameSetter[Tu]](from, to []T, getDiff func(T, T, bool) Tu) ([]Tu, bool, error) {
	if arrayObjectNamedHasDuplicateKeys(from) || arrayObjectNamedHasDuplicateKeys(to) {
		return nil, false, ErrDuplicateKeys
	}

	value := make([]Tu, 0, len(from))
	for _, toItem := range to {
		fromItemFound := false
		for _, fromItem := range from {
			if toItem.GetName() == fromItem.GetName() {
				fromItemFound = true
				tmp := Tns(ptr.Get(getDiff(fromItem, toItem, false)))
				tmp.SetName(toItem.GetName())
				value = append(value, *tmp)
				break
			}
		}
		if !fromItemFound {
			tmp := Tns(ptr.Get(getDiff(toItem, toItem, true)))
			tmp.SetName(toItem.GetName())
			value = append(value, *tmp)
		}
	}

	hasChanges := len(to) != len(from)
	if !hasChanges {
		for _, item := range value {
			if item.HasChanges() {
				hasChanges = true
				break
			}
		}
	}

	return value, hasChanges, nil
}

func GetChangesArrayObjectNamedError[T namedChild, Tu updateObject, Tns nameSetter[Tu]](
	from, to []T,
	getDiff func(T, T, bool) (Tu, error),
) ([]Tu, bool, error) {
	if arrayObjectNamedHasDuplicateKeys(from) || arrayObjectNamedHasDuplicateKeys(to) {
		return nil, false, ErrDuplicateKeys
	}

	value := make([]Tu, 0, len(from))
	for _, toItem := range to {
		fromItemFound := false
		for _, fromItem := range from {
			if toItem.GetName() == fromItem.GetName() {
				fromItemFound = true
				diffValue, err := getDiff(fromItem, toItem, false)
				if err != nil {
					return nil, false, err
				}
				tmp := Tns(&diffValue)
				tmp.SetName(toItem.GetName())
				value = append(value, *tmp)
				break
			}
		}
		if !fromItemFound {
			diffValue, err := getDiff(toItem, toItem, true)
			if err != nil {
				return nil, false, err
			}
			tmp := Tns(&diffValue)
			tmp.SetName(toItem.GetName())
			value = append(value, *tmp)
		}
	}

	hasChanges := len(to) != len(from)
	if !hasChanges {
		for _, item := range value {
			if item.HasChanges() {
				hasChanges = true
				break
			}
		}
	}

	return value, hasChanges, nil
}

func GetChangesMapPrimitive[T comparable](from, to map[string]T) (map[string]T, bool) {
	hasChanges := len(from) != len(to)
	if !hasChanges {
		for key, toV := range to {
			if fromV, ok := from[key]; !ok || fromV != toV {
				hasChanges = true
				break
			}
		}
	}
	return to, hasChanges
}

func GetChangesMapEquatableIface[T Equatable[T]](from, to map[string]T) (map[string]T, bool) {
	hasChanges := len(from) != len(to)
	if !hasChanges {
		for key, toV := range to {
			if fromV, ok := from[key]; !ok || !fromV.Equal(toV) {
				hasChanges = true
				break
			}
		}
	}
	return to, hasChanges
}

func GetChangesMapRawData(from, to map[string]json.RawMessage) (map[string]json.RawMessage, bool) {
	hasChanges := len(from) != len(to)
	if !hasChanges {
		for key, toV := range to {
			if fromV, ok := from[key]; !ok || !slices.Equal(fromV, toV) {
				hasChanges = true
				break
			}
		}
	}
	return to, hasChanges
}

func GetChangesMapObject[T any, Tu updateObject](from, to map[string]T, getDiff func(T, T, bool) Tu) (map[string]Tu, bool) {
	value := make(map[string]Tu, 0)
	for key, toItem := range to {
		if from != nil {
			fromItem, ok := from[key]
			if ok {
				value[key] = getDiff(fromItem, toItem, false)
			} else {
				value[key] = getDiff(toItem, toItem, true)
			}
		} else {
			value[key] = getDiff(toItem, toItem, true)
		}
	}

	hasChanges := len(to) != len(from)
	if !hasChanges {
		for _, item := range value {
			if item.HasChanges() {
				hasChanges = true
				break
			}
		}
	}
	return value, hasChanges
}

func arrayObjectNamedHasDuplicateKeys[T namedChild](arr []T) bool {
	m := make(map[string]struct{}, len(arr))

	for _, v := range arr {
		if _, ok := m[v.GetName()]; !ok {
			m[v.GetName()] = struct{}{}
		} else {
			return true
		}
	}

	return false
}
