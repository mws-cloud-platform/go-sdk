package client

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

func TestDiffPrimitive(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       string
		to         string
		expected   string
		hasChanges bool
	}{
		{
			name:       "change",
			from:       "hello",
			to:         "world",
			expected:   "world",
			hasChanges: true,
		},
		{
			name:       "no change",
			from:       "hello",
			to:         "hello",
			expected:   "hello",
			hasChanges: false,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			diffReq := DiffPrimitiveRequired(v.from, v.to, false)
			require.Equal(t, v.expected, diffReq.Value)
			require.Equal(t, v.hasChanges, diffReq.Set)

			diffNonReq := DiffPrimitiveNonRequired(&v.from, &v.to, false)
			require.Equal(t, v.expected, diffNonReq.Value)
			require.Equal(t, v.hasChanges, diffNonReq.Set)

			diffNullable := DiffPrimitiveNullable(&v.from, &v.to, false)
			require.Equal(t, v.expected, diffNullable.Value)
			require.Equal(t, v.hasChanges, diffNullable.Set)
			require.False(t, diffNullable.Null)
		})
	}
}

func TestDiffPrimitiveNulls(t *testing.T) {
	for _, v := range []struct {
		name     string
		from     *string
		to       *string
		expected *string
	}{
		{
			name:     "from",
			from:     nil,
			to:       ptr.Get("world"),
			expected: ptr.Get("world"),
		},
		{
			name:     "to",
			from:     ptr.Get("hello"),
			to:       nil,
			expected: nil,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			nilDiffers := v.from == nil && v.to != nil ||
				v.from != nil && v.to == nil
			diffNonReq := DiffPrimitiveNonRequired(v.from, v.to, nilDiffers)
			if v.expected == nil {
				require.Empty(t, diffNonReq.Value)
			} else {
				require.Equal(t, *v.expected, diffNonReq.Value)
			}
			require.True(t, diffNonReq.Set)

			diffNull := DiffPrimitiveNullable(v.from, v.to, nilDiffers)
			if v.expected == nil {
				require.Empty(t, diffNull.Value)
			} else {
				require.Equal(t, *v.expected, diffNull.Value)
			}
			require.True(t, diffNull.Set)
			require.Equal(t, v.to == nil, diffNull.Null)
		})
	}
}

type equatableFoo struct {
	m map[string]string
	s string
}

func (c equatableFoo) Equal(c2 equatableFoo) bool {
	return maps.Equal(c.m, c2.m) && c.s == c2.s
}

func TestDiffEquatableIface(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       equatableFoo
		to         equatableFoo
		expected   equatableFoo
		hasChanges bool
	}{
		{
			name: "change",
			from: equatableFoo{
				m: map[string]string{"hello": "world"},
			},
			to: equatableFoo{
				m: map[string]string{"foo": "bar"},
			},
			expected: equatableFoo{
				m: map[string]string{"foo": "bar"},
			},
			hasChanges: true,
		},
		{
			name: "no change",
			from: equatableFoo{
				s: "hello",
			},
			to: equatableFoo{
				s: "hello",
			},
			expected: equatableFoo{
				s: "hello",
			},
			hasChanges: false,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			diffReq := DiffEquatableIfaceRequired(v.from, v.to, false)
			require.Equal(t, v.expected, diffReq.Value)
			require.Equal(t, v.hasChanges, diffReq.Set)

			diffNonReq := DiffEquatableIfaceNonRequired(&v.from, &v.to, false)
			require.Equal(t, v.expected, diffNonReq.Value)
			require.Equal(t, v.hasChanges, diffNonReq.Set)

			diffNullable := DiffEquatableIfaceNullable(&v.from, &v.to, false)
			require.Equal(t, v.expected, diffNullable.Value)
			require.Equal(t, v.hasChanges, diffNullable.Set)
			require.False(t, diffNullable.Null)
		})
	}
}

func TestDiffEquatableIfaceNulls(t *testing.T) {
	for _, v := range []struct {
		name     string
		from     *equatableFoo
		to       *equatableFoo
		expected *equatableFoo
	}{
		{
			name: "from",
			from: nil,
			to: &equatableFoo{
				m: map[string]string{"hello": "world"},
			},
			expected: &equatableFoo{
				m: map[string]string{"hello": "world"},
			},
		},
		{
			name: "to",
			from: &equatableFoo{
				m: map[string]string{"hello": "world"},
			},
			to:       nil,
			expected: nil,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			nilDiffers := v.from == nil && v.to != nil ||
				v.from != nil && v.to == nil
			diffNonReq := DiffEquatableIfaceNonRequired(v.from, v.to, nilDiffers)
			if v.expected == nil {
				require.Empty(t, diffNonReq.Value)
			} else {
				require.Equal(t, *v.expected, diffNonReq.Value)
			}
			require.True(t, diffNonReq.Set)

			diffNull := DiffEquatableIfaceNullable(v.from, v.to, nilDiffers)
			if v.expected == nil {
				require.Empty(t, diffNull.Value)
			} else {
				require.Equal(t, *v.expected, diffNull.Value)
			}
			require.True(t, diffNull.Set)
			require.Equal(t, v.to == nil, diffNull.Null)
		})
	}
}

func TestDiffRawData(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       json.RawMessage
		to         json.RawMessage
		expected   json.RawMessage
		hasChanges bool
	}{
		{
			name:       "change",
			from:       json.RawMessage("hello"),
			to:         json.RawMessage("world"),
			expected:   json.RawMessage("world"),
			hasChanges: true,
		},
		{
			name:       "no change",
			from:       json.RawMessage("hello"),
			to:         json.RawMessage("hello"),
			expected:   json.RawMessage("hello"),
			hasChanges: false,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			diffReq := DiffRawData(v.from, v.to)
			require.Equal(t, v.expected, diffReq.Value)
			require.Equal(t, v.hasChanges, diffReq.Set)
		})
	}
}

func TestDiffRawDataNullable(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       json.RawMessage
		to         json.RawMessage
		expected   json.RawMessage
		hasChanges bool
		isNull     bool
	}{
		{
			name:       "change",
			from:       json.RawMessage("hello"),
			to:         json.RawMessage("world"),
			expected:   json.RawMessage("world"),
			hasChanges: true,
		},
		{
			name:       "no change",
			from:       json.RawMessage("hello"),
			to:         json.RawMessage("hello"),
			expected:   json.RawMessage("hello"),
			hasChanges: false,
		},
		{
			name:       "to nil",
			from:       json.RawMessage("hello"),
			to:         nil,
			expected:   nil,
			hasChanges: true,
			isNull:     true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			diffReq := DiffRawDataNullable(v.from, v.to)
			require.Equal(t, v.expected, diffReq.Value)
			require.Equal(t, v.hasChanges, diffReq.Set)
			require.Equal(t, v.isNull, diffReq.Null)
		})
	}
}

func TestGetChangesArrayPrimitive(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       []string
		to         []string
		hasChanges bool
	}{
		{
			name:       "no change",
			from:       []string{"hello", "world"},
			to:         []string{"hello", "world"},
			hasChanges: false,
		},
		{
			name:       "reorder",
			from:       []string{"hello", "world"},
			to:         []string{"world", "hello"},
			hasChanges: true,
		},
		{
			name:       "add",
			from:       []string{"hello"},
			to:         []string{"hello", "world"},
			hasChanges: true,
		},
		{
			name:       "remove",
			from:       []string{"hello", "world"},
			to:         []string{"world"},
			hasChanges: true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			value, hasChanges := GetChangesArrayPrimitive(v.from, v.to)
			require.ElementsMatch(t, v.to, value)
			require.Equal(t, v.hasChanges, hasChanges)
		})
	}
}

func TestGetChangesArrayEquatableIface(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       []equatableFoo
		to         []equatableFoo
		hasChanges bool
	}{
		{
			name:       "no change",
			from:       []equatableFoo{{s: "hello"}, {s: "hello"}},
			to:         []equatableFoo{{s: "hello"}, {s: "hello"}},
			hasChanges: false,
		},
		{
			name:       "reorder",
			from:       []equatableFoo{{s: "hello"}, {s: "world"}},
			to:         []equatableFoo{{s: "world"}, {s: "hello"}},
			hasChanges: true,
		},
		{
			name:       "add",
			from:       []equatableFoo{{s: "hello"}},
			to:         []equatableFoo{{s: "hello"}, {s: "world"}},
			hasChanges: true,
		},
		{
			name:       "remove",
			from:       []equatableFoo{{s: "hello"}, {s: "world"}},
			to:         []equatableFoo{{s: "hello"}},
			hasChanges: true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			value, hasChanges := GetChangesArrayEquatableIface(v.from, v.to)
			require.ElementsMatch(t, v.to, value)
			require.Equal(t, v.hasChanges, hasChanges)
		})
	}
}

func TestGetChangesArrayRawData(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       []json.RawMessage
		to         []json.RawMessage
		hasChanges bool
	}{
		{
			name:       "no change",
			from:       []json.RawMessage{json.RawMessage("hello"), json.RawMessage("world")},
			to:         []json.RawMessage{json.RawMessage("hello"), json.RawMessage("world")},
			hasChanges: false,
		},
		{
			name:       "reorder",
			from:       []json.RawMessage{json.RawMessage("hello"), json.RawMessage("world")},
			to:         []json.RawMessage{json.RawMessage("world"), json.RawMessage("hello")},
			hasChanges: true,
		},
		{
			name:       "add",
			from:       []json.RawMessage{json.RawMessage("hello")},
			to:         []json.RawMessage{json.RawMessage("hello"), json.RawMessage("world")},
			hasChanges: true,
		},
		{
			name:       "remove",
			from:       []json.RawMessage{json.RawMessage("hello"), json.RawMessage("world")},
			to:         []json.RawMessage{json.RawMessage("world")},
			hasChanges: true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			value, hasChanges := GetChangesArrayRawData(v.from, v.to)
			require.ElementsMatch(t, v.to, value)
			require.Equal(t, v.hasChanges, hasChanges)
		})
	}
}

type element struct {
	name string
	num  int
}

func (e *element) GetName() string {
	return e.name
}

type updateElement struct {
	hasChanges bool
}

func (u updateElement) HasChanges() bool {
	return u.hasChanges
}

func (u *updateElement) SetName(string) {}

func elementDiffFunc(from, to element, fromNil bool) updateElement {
	if fromNil {
		return updateElement{hasChanges: true}
	}
	return updateElement{
		hasChanges: from.num != to.num ||
			from.name != to.name,
	}
}

func elementDiffFuncError(from, to element, fromNil bool) (updateElement, error) {
	return elementDiffFunc(from, to, fromNil), nil
}

func elementDiffFuncRef(from, to *element, fromNil bool) updateElement {
	if fromNil && to != nil ||
		!fromNil && to == nil {
		return updateElement{hasChanges: true}
	}
	if from == nil && to == nil {
		return updateElement{hasChanges: false}
	}
	return elementDiffFunc(*from, *to, fromNil)
}

func elementDiffFuncRefError(from, to *element, fromNil bool) (updateElement, error) {
	return elementDiffFuncRef(from, to, fromNil), nil
}

func TestGetChangesArrayObject(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       []element
		to         []element
		expected   []updateElement
		hasChanges bool
	}{
		{
			name:       "no changes",
			from:       []element{{num: 1}, {num: 2}},
			to:         []element{{num: 1}, {num: 2}},
			expected:   []updateElement{{hasChanges: false}, {hasChanges: false}},
			hasChanges: false,
		},
		{
			name:       "reorder",
			from:       []element{{num: 1}, {num: 2}},
			to:         []element{{num: 2}, {num: 1}},
			expected:   []updateElement{{hasChanges: true}, {hasChanges: true}},
			hasChanges: true,
		},
		{
			name:       "add",
			from:       []element{{num: 1}},
			to:         []element{{num: 1}, {num: 2}},
			expected:   []updateElement{{hasChanges: true}, {hasChanges: true}},
			hasChanges: true,
		},
		{
			name:       "remove",
			from:       []element{{num: 1}, {num: 2}},
			to:         []element{{num: 2}},
			expected:   []updateElement{{hasChanges: true}},
			hasChanges: true,
		},
		{
			name:       "compare with default",
			from:       []element{{num: 1}, {num: 0}},
			to:         []element{{num: 0}},
			expected:   []updateElement{{hasChanges: true}},
			hasChanges: true,
		},
		{
			name:       "from nil",
			from:       nil,
			to:         []element{{num: 1}, {num: 2}},
			expected:   []updateElement{{hasChanges: true}, {hasChanges: true}},
			hasChanges: true,
		},
		{
			name:       "to nil",
			from:       []element{{num: 1}, {num: 2}},
			to:         nil,
			expected:   []updateElement{},
			hasChanges: true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			actual, hasChanges := GetChangesArrayObject(v.from, v.to, elementDiffFunc)
			require.Equal(t, v.expected, actual)
			require.Equal(t, v.hasChanges, hasChanges)
		})
	}
}

func TestGetChangesArrayObjectError(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       []element
		to         []element
		expected   []updateElement
		hasChanges bool
		hasError   bool
		diffFunc   func(from, to element, fromNil bool) (updateElement, error)
	}{
		{
			name:       "no changes",
			from:       []element{{num: 1}, {num: 2}},
			to:         []element{{num: 1}, {num: 2}},
			expected:   []updateElement{{hasChanges: false}, {hasChanges: false}},
			hasChanges: false,
			diffFunc:   elementDiffFuncError,
		},
		{
			name:       "reorder",
			from:       []element{{num: 1}, {num: 2}},
			to:         []element{{num: 2}, {num: 1}},
			expected:   []updateElement{{hasChanges: true}, {hasChanges: true}},
			hasChanges: true,
			diffFunc:   elementDiffFuncError,
		},
		{
			name:       "add",
			from:       []element{{num: 1}},
			to:         []element{{num: 1}, {num: 2}},
			expected:   []updateElement{{hasChanges: true}, {hasChanges: true}},
			hasChanges: true,
			diffFunc:   elementDiffFuncError,
		},
		{
			name:       "remove",
			from:       []element{{num: 1}, {num: 2}},
			to:         []element{{num: 2}},
			expected:   []updateElement{{hasChanges: true}},
			hasChanges: true,
			diffFunc:   elementDiffFuncError,
		},
		{
			name:       "compare with default",
			from:       []element{{num: 1}, {num: 0}},
			to:         []element{{num: 0}},
			expected:   []updateElement{{hasChanges: true}},
			hasChanges: true,
			diffFunc:   elementDiffFuncError,
		},
		{
			name:       "from nil",
			from:       nil,
			to:         []element{{num: 1}, {num: 2}},
			expected:   []updateElement{{hasChanges: true}, {hasChanges: true}},
			hasChanges: true,
			diffFunc:   elementDiffFuncError,
		},
		{
			name:       "to nil",
			from:       []element{{num: 1}, {num: 2}},
			to:         nil,
			expected:   []updateElement{},
			hasChanges: true,
			diffFunc:   elementDiffFuncError,
		},
		{
			name:       "diff func returned error",
			from:       []element{{num: 1}, {num: 2}},
			to:         nil,
			expected:   []updateElement{},
			hasChanges: true,
			diffFunc: func(_, _ element, _ bool) (updateElement, error) {
				return updateElement{hasChanges: false}, fmt.Errorf("diff func error")
			},
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			actual, hasChanges, err := GetChangesArrayObjectError(v.from, v.to, v.diffFunc)
			if v.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, v.expected, actual)
				require.Equal(t, v.hasChanges, hasChanges)
			}
		})
	}
}

func TestGetChangesArrayObjectNamed(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       []element
		to         []element
		expected   []updateElement
		hasChanges bool
		hasError   bool
	}{
		{
			name:       "no changes",
			from:       []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			to:         []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			expected:   []updateElement{{hasChanges: false}, {hasChanges: false}},
			hasChanges: false,
		},
		{
			name:       "reorder",
			from:       []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			to:         []element{{num: 2, name: "2"}, {num: 1, name: "1"}},
			expected:   []updateElement{{hasChanges: false}, {hasChanges: false}},
			hasChanges: false,
		},
		{
			name:       "add",
			from:       []element{{num: 1, name: "1"}},
			to:         []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			expected:   []updateElement{{hasChanges: false}, {hasChanges: true}},
			hasChanges: true,
		},
		{
			name:       "remove",
			from:       []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			to:         []element{{num: 2, name: "2"}},
			expected:   []updateElement{{hasChanges: false}},
			hasChanges: true,
		},
		{
			name:       "compare with default",
			from:       []element{{num: 1, name: "1"}, {num: 0, name: "0"}},
			to:         []element{{num: 0, name: "0"}},
			expected:   []updateElement{{hasChanges: false}},
			hasChanges: true,
		},
		{
			name:       "from nil",
			from:       nil,
			to:         []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			expected:   []updateElement{{hasChanges: true}, {hasChanges: true}},
			hasChanges: true,
		},
		{
			name:       "to nil",
			from:       []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			to:         nil,
			expected:   []updateElement{},
			hasChanges: true,
		},
		{
			name:     "to duplicate key",
			from:     nil,
			to:       []element{{num: 1}, {num: 2}},
			hasError: true,
		},
		{
			name:     "from duplicate key",
			from:     []element{{num: 1}, {num: 2}},
			to:       nil,
			hasError: true,
		},
		{
			name:       "from nil to nil",
			from:       nil,
			to:         nil,
			expected:   []updateElement{},
			hasChanges: false,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			actual, hasChanges, err := GetChangesArrayObjectNamed(ToPointerArray(v.from), ToPointerArray(v.to), elementDiffFuncRef)
			if v.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, v.expected, actual)
				require.Equal(t, v.hasChanges, hasChanges)
			}
		})
	}
}

func TestGetChangesArrayObjectNamedError(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       []element
		to         []element
		expected   []updateElement
		diffFunc   func(from, to *element, fromNil bool) (updateElement, error)
		hasChanges bool
		hasError   bool
	}{
		{
			name:       "no changes",
			from:       []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			to:         []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			expected:   []updateElement{{hasChanges: false}, {hasChanges: false}},
			diffFunc:   elementDiffFuncRefError,
			hasChanges: false,
		},
		{
			name:       "reorder",
			from:       []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			to:         []element{{num: 2, name: "2"}, {num: 1, name: "1"}},
			expected:   []updateElement{{hasChanges: false}, {hasChanges: false}},
			diffFunc:   elementDiffFuncRefError,
			hasChanges: false,
		},
		{
			name:       "add",
			from:       []element{{num: 1, name: "1"}},
			to:         []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			expected:   []updateElement{{hasChanges: false}, {hasChanges: true}},
			diffFunc:   elementDiffFuncRefError,
			hasChanges: true,
		},
		{
			name:       "remove",
			from:       []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			to:         []element{{num: 2, name: "2"}},
			expected:   []updateElement{{hasChanges: false}},
			diffFunc:   elementDiffFuncRefError,
			hasChanges: true,
		},
		{
			name:       "compare with default",
			from:       []element{{num: 1, name: "1"}, {num: 0, name: "0"}},
			to:         []element{{num: 0, name: "0"}},
			expected:   []updateElement{{hasChanges: false}},
			diffFunc:   elementDiffFuncRefError,
			hasChanges: true,
		},
		{
			name:       "from nil",
			from:       nil,
			to:         []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			expected:   []updateElement{{hasChanges: true}, {hasChanges: true}},
			diffFunc:   elementDiffFuncRefError,
			hasChanges: true,
		},
		{
			name:       "to nil",
			from:       []element{{num: 1, name: "1"}, {num: 2, name: "2"}},
			to:         nil,
			expected:   []updateElement{},
			diffFunc:   elementDiffFuncRefError,
			hasChanges: true,
		},
		{
			name:     "to duplicate key",
			from:     nil,
			to:       []element{{num: 1}, {num: 2}},
			diffFunc: elementDiffFuncRefError,
			hasError: true,
		},
		{
			name:     "from duplicate key",
			from:     []element{{num: 1}, {num: 2}},
			to:       nil,
			diffFunc: elementDiffFuncRefError,
			hasError: true,
		},
		{
			name:       "from nil to nil",
			from:       nil,
			to:         nil,
			expected:   []updateElement{},
			diffFunc:   elementDiffFuncRefError,
			hasChanges: false,
		},
		{
			name:     "diff func returned error",
			from:     []element{{num: 1, name: "1"}, {num: 0, name: "0"}},
			to:       []element{{num: 0, name: "0"}},
			expected: []updateElement{{hasChanges: false}},
			diffFunc: func(_, _ *element, _ bool) (updateElement, error) {
				return updateElement{hasChanges: false}, fmt.Errorf("diff func error")
			},
			hasError: true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			actual, hasChanges, err := GetChangesArrayObjectNamedError(ToPointerArray(v.from), ToPointerArray(v.to), v.diffFunc)
			if v.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, v.expected, actual)
				require.Equal(t, v.hasChanges, hasChanges)
			}
		})
	}
}

func TestGetChangesMapPrimitive(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       map[string]string
		to         map[string]string
		hasChanges bool
	}{
		{
			name:       "no change",
			from:       map[string]string{"hello": "world", "hi": "mark"},
			to:         map[string]string{"hi": "mark", "hello": "world"},
			hasChanges: false,
		},
		{
			name:       "add",
			from:       map[string]string{"hello": "world"},
			to:         map[string]string{"hello": "world", "hi": "mark"},
			hasChanges: true,
		},
		{
			name:       "remove",
			from:       map[string]string{"hello": "world", "hi": "mark"},
			to:         map[string]string{"hello": "world"},
			hasChanges: true,
		},
		{
			name:       "change",
			from:       map[string]string{"hello": "world"},
			to:         map[string]string{"hi": "world"},
			hasChanges: true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			actual, hasChanges := GetChangesMapPrimitive(v.from, v.to)

			require.True(t, reflect.DeepEqual(v.to, actual))
			require.Equal(t, v.hasChanges, hasChanges)
		})
	}
}

func TestGetChangesMapEquatableIface(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       map[string]equatableFoo
		to         map[string]equatableFoo
		hasChanges bool
	}{
		{
			name:       "no change",
			from:       map[string]equatableFoo{"hello": {s: "world"}, "hi": {s: "mark"}},
			to:         map[string]equatableFoo{"hi": {s: "mark"}, "hello": {s: "world"}},
			hasChanges: false,
		},
		{
			name:       "add",
			from:       map[string]equatableFoo{"hello": {s: "world"}},
			to:         map[string]equatableFoo{"hello": {s: "world"}, "hi": {s: "mark"}},
			hasChanges: true,
		},
		{
			name:       "remove",
			from:       map[string]equatableFoo{"hello": {s: "world"}, "hi": {s: "mark"}},
			to:         map[string]equatableFoo{"hello": {s: "world"}},
			hasChanges: true,
		},
		{
			name:       "change",
			from:       map[string]equatableFoo{"hello": {s: "world"}},
			to:         map[string]equatableFoo{"hi": {s: "world"}},
			hasChanges: true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			actual, hasChanges := GetChangesMapEquatableIface(v.from, v.to)

			require.True(t, reflect.DeepEqual(v.to, actual))
			require.Equal(t, v.hasChanges, hasChanges)
		})
	}
}

func TestGetChangesMapRawData(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       map[string]json.RawMessage
		to         map[string]json.RawMessage
		hasChanges bool
	}{
		{
			name:       "no change",
			from:       map[string]json.RawMessage{"hello": json.RawMessage(`"world"`), "hi": json.RawMessage(`"mark"`)},
			to:         map[string]json.RawMessage{"hi": json.RawMessage(`"mark"`), "hello": json.RawMessage(`"world"`)},
			hasChanges: false,
		},
		{
			name:       "add",
			from:       map[string]json.RawMessage{"hello": json.RawMessage(`"world"`)},
			to:         map[string]json.RawMessage{"hello": json.RawMessage(`"world"`), "hi": json.RawMessage(`"mark"`)},
			hasChanges: true,
		},
		{
			name:       "remove",
			from:       map[string]json.RawMessage{"hello": json.RawMessage(`"world"`), "hi": json.RawMessage(`"mark"`)},
			to:         map[string]json.RawMessage{"hello": json.RawMessage(`"world"`)},
			hasChanges: true,
		},
		{
			name:       "change",
			from:       map[string]json.RawMessage{"hello": json.RawMessage(`"world"`)},
			to:         map[string]json.RawMessage{"hi": json.RawMessage(`"world"`)},
			hasChanges: true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			actual, hasChanges := GetChangesMapRawData(v.from, v.to)

			require.True(t, reflect.DeepEqual(v.to, actual))
			require.Equal(t, v.hasChanges, hasChanges)
		})
	}
}

func TestGetChangesMapObject(t *testing.T) {
	for _, v := range []struct {
		name       string
		from       map[string]element
		to         map[string]element
		expected   map[string]updateElement
		hasChanges bool
	}{
		{
			name:       "no change",
			from:       map[string]element{"hello": {num: 1}, "hi": {num: 2}},
			to:         map[string]element{"hello": {num: 1}, "hi": {num: 2}},
			expected:   map[string]updateElement{"hello": {hasChanges: false}, "hi": {hasChanges: false}},
			hasChanges: false,
		},
		{
			name:       "add",
			from:       map[string]element{"hello": {num: 1}},
			to:         map[string]element{"hello": {num: 1}, "hi": {num: 2}},
			expected:   map[string]updateElement{"hello": {hasChanges: false}, "hi": {hasChanges: true}},
			hasChanges: true,
		},
		{
			name:       "remove",
			from:       map[string]element{"hello": {num: 1}, "hi": {num: 2}},
			to:         map[string]element{"hello": {num: 1}},
			expected:   map[string]updateElement{"hello": {hasChanges: false}},
			hasChanges: true,
		},
		{
			name:       "change",
			from:       map[string]element{"hello": {num: 1}},
			to:         map[string]element{"hello": {num: 2}},
			expected:   map[string]updateElement{"hello": {hasChanges: true}},
			hasChanges: true,
		},
		{
			name:       "change default",
			from:       map[string]element{},
			to:         map[string]element{"hello": {num: 0}},
			expected:   map[string]updateElement{"hello": {hasChanges: true}},
			hasChanges: true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			actual, hasChanges := GetChangesMapObject(v.from, v.to, elementDiffFunc)

			require.True(t, reflect.DeepEqual(v.expected, actual))
			require.Equal(t, v.hasChanges, hasChanges)
		})
	}
}
