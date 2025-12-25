package metadata

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataFromContextPanicsOnType(t *testing.T) {
	ctx := context.WithValue(t.Context(), keyOutgoing{}, 13)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected to panic when value is of the wrong type")
		}
	}()

	FromOutgoingContext(ctx)
}

func TestMetadataFromContextReturnsCopy(t *testing.T) {
	first := FromOutgoingContext(t.Context())
	first["key"] = []string{""}

	second := FromOutgoingContext(t.Context())
	_, ok := second["key"]

	require.NotEqual(t, first, second)
	require.False(t, ok)
}

func TestMetadataFromContextGetValue(t *testing.T) {
	m := Metadata{
		"key": []string{"value"},
	}
	ctx := context.WithValue(t.Context(), keyOutgoing{}, m)

	extracted := FromOutgoingContext(ctx)
	value, ok := extracted["key"]

	require.True(t, ok)
	require.Equal(t, "value", value[0])
}

func TestMetadataMultipleWithMetadataValue(t *testing.T) {
	ctx := t.Context()
	ctx = WithOutgoingMetadataValue(ctx, "key", "value", "another")
	ctx = WithOutgoingMetadataValue(ctx, "key2", "value2")

	m := FromOutgoingContext(ctx)
	value, ok := m["key"]
	value2, ok2 := m["key2"]

	require.True(t, ok)
	require.Equal(t, []string{"value", "another"}, value)

	require.True(t, ok2)
	require.Equal(t, "value2", value2[0])
}

func TestMetadataWithMetadataCopy(t *testing.T) {
	m := map[string][]string{
		"key": {"value"},
	}
	ctx := WithIncomingMetadata(t.Context(), m)
	written := ctx.Value(keyIncoming{}).(Metadata)
	from := FromOutgoingContext(ctx)

	require.EqualValues(t, m, written)
	require.NotEqualValues(t, written, from)
}

func TestMetadataNotSameMaps(t *testing.T) {
	ctx := WithOutgoingMetadataValue(t.Context(), "keyout", "valueout")
	ctx = WithIncomingMetadata(ctx, map[string][]string{"keyin": {"valuein"}})

	in := ctx.Value(keyIncoming{}).(Metadata)
	out := ctx.Value(keyOutgoing{}).(Metadata)

	valueOut := out["keyout"]
	valueIn := in["keyin"]
	require.NotEqualValues(t, in, out)
	require.Equal(t, "valueout", valueOut[0])
	require.Equal(t, "valuein", valueIn[0])

	_, outHasIn := out["keyin"]
	_, inHasOut := in["keyout"]
	require.False(t, outHasIn)
	require.False(t, inHasOut)
}

func TestMetadataGetValue(t *testing.T) {
	ctx := WithOutgoingMetadataValue(t.Context(), "keyout", "valueout")
	ctx = WithIncomingMetadata(ctx, map[string][]string{"keyin": {"valuein"}})

	valueIn, hasIn := GetIncomingMetadataValue(ctx, "keyin")
	valueOut, hasOut := GetIncomingMetadataValue(ctx, "keyout")

	require.True(t, hasIn)
	require.Equal(t, "valuein", valueIn)
	require.False(t, hasOut)
	require.Equal(t, "", valueOut)
}

func TestMetadataIncomingMerge(t *testing.T) {
	ctx := WithIncomingMetadata(t.Context(), map[string][]string{"key1": {"value1"}, "key": {"initial"}})
	ctx = WithIncomingMetadata(ctx, map[string][]string{"key2": {"value2"}, "key": {"replaced"}})

	v1, ok1 := GetIncomingMetadataValue(ctx, "key1")
	v2, ok2 := GetIncomingMetadataValue(ctx, "key2")
	vCommon, okCommon := GetIncomingMetadataValue(ctx, "key")
	vsCommon, okVsCommon := GetIncomingMetadataValues(ctx, "key")

	require.True(t, ok1)
	require.Equal(t, "value1", v1)
	require.True(t, ok2)
	require.Equal(t, "value2", v2)
	require.True(t, okCommon)
	require.Equal(t, "initial", vCommon)
	require.True(t, okVsCommon)
	require.Equal(t, []string{"initial", "replaced"}, vsCommon)
}

func TestMetadataOutgoingAppend(t *testing.T) {
	ctx := t.Context()
	ctx = WithOutgoingMetadataValue(ctx, "key", "value")
	ctx = AppendOutgoingMetadataValue(ctx, "key", "another")

	m := FromOutgoingContext(ctx)
	value, ok := m["key"]

	require.True(t, ok)
	require.Equal(t, []string{"value", "another"}, value)
}

func TestMetadataOutgoingReplace(t *testing.T) {
	ctx := t.Context()
	ctx = WithOutgoingMetadataValue(ctx, "key", "value")
	ctx = WithOutgoingMetadataValue(ctx, "key", "another")

	m := FromOutgoingContext(ctx)
	value, ok := m["key"]

	require.True(t, ok)
	require.Equal(t, []string{"another"}, value)
}

func TestMetadataValueCopyOutgoing(t *testing.T) {
	arr := make([]string, 1)
	arr[0] = "value"
	ctx := WithOutgoingMetadataValue(t.Context(), "key", arr...)
	arr[0] = "changed"

	m := FromOutgoingContext(ctx)

	require.Equal(t, []string{"value"}, m["key"])
}

func TestMetadataValueCopyAppendOutgoing(t *testing.T) {
	arr := make([]string, 1)
	arr[0] = "value"
	ctx := AppendOutgoingMetadataValue(t.Context(), "key", arr...)
	arr[0] = "changed"

	m := FromOutgoingContext(ctx)

	require.Equal(t, []string{"value"}, m["key"])
}

func TestMetadataValueCopyIncoming(t *testing.T) {
	arr := make([]string, 2)
	arr[0] = "value"
	ctx := WithIncomingMetadata(t.Context(), map[string][]string{"key": arr})
	arr[0] = "changed"

	v, _ := GetIncomingMetadataValues(ctx, "key")
	require.Equal(t, []string{"value", ""}, v)

	v[1] = "another"
	v, _ = GetIncomingMetadataValues(ctx, "key")
	require.Equal(t, []string{"value", ""}, v)
}
