package attribute_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/internal/client/interceptors/attribute"
)

func TestContext(t *testing.T) {
	attributes := []attribute.KeyValue{
		attribute.String("foo", "bar"),
		attribute.String("hello", "world"),
	}

	require.Equal(t, attributes, attribute.FromContext(attribute.WithContext(t.Context(), attributes...)))
}
