package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlattenJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected string
	}{
		{
			name:     "nil_map",
			input:    nil,
			expected: "",
		},
		{
			name:     "empty_map",
			input:    map[string]any{},
			expected: "",
		},
		{
			name: "structure",
			input: map[string]any{
				"name": "Alice",
				"age":  25,
			},
			expected: "age - 25\nname - Alice",
		},
		{
			name: "nested_structure",
			input: map[string]any{
				"user": map[string]any{
					"name": "Bob",
					"meta": map[string]any{
						"role":   "admin",
						"active": true,
					},
				},
				"version": 1,
			},
			expected: "user.meta.active - true\nuser.meta.role - admin\nuser.name - Bob\nversion - 1",
		},
		{
			name: "different_types",
			input: map[string]any{
				"float":  3.14,
				"bool":   false,
				"null":   nil,
				"array":  []any{1, 2},
				"nested": map[string]any{"a": 1},
			},
			expected: "array - [1 2]\nbool - false\nfloat - 3.14\nnested.a - 1\nnull - <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenJSON(tt.input, "")
			require.Equal(t, tt.expected, result)
		})
	}
}
