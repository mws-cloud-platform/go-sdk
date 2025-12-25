package values

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValues(t *testing.T) {
	ctx := t.Context()

	val, ok := From(ctx, "something")
	require.Equal(t, "", val)
	require.False(t, ok)

	ctx = With(ctx, "hello", "world")
	ctx = With(ctx, "foo", "bar")

	val, ok = From(ctx, "hello")
	require.Equal(t, "world", val)
	require.True(t, ok)

	val, ok = From(ctx, "foo")
	require.Equal(t, "bar", val)
	require.True(t, ok)
}
