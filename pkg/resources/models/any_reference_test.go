package models

import (
	"testing"

	"github.com/stretchr/testify/require"

	reserrors "go.mws.cloud/go-sdk/internal/resources/errors"
	"go.mws.cloud/go-sdk/pkg/resources/interfaces"
)

func TestAnyResourceRefMethods(t *testing.T) {
	for _, v := range []struct {
		id                   *AnyResourceRef
		expectedSlug         string
		expectedResourceName interfaces.ResourceName
	}{
		{nil, anyServiceSlug, ""},
		{new(NewMustAnyResourceRef("hello")), anyServiceSlug, "hello"},
		{new(NewMustAnyResourceRef("/projects")), anyServiceSlug, "projects"},
		{new(NewMustAnyResourceRef("/projects/project")), anyServiceSlug, "project"},
		{new(NewMustAnyResourceRef("/projects/project/object")), anyServiceSlug, "object"},
		{new(NewMustAnyResourceRef("some/projects/project")), anyServiceSlug, "project"},
		{new(NewMustAnyResourceRef("some/projects/project/object")), anyServiceSlug, "object"},
		{new(NewMustAnyResourceRef("some/projects/project/objects/object")), anyServiceSlug, "object"},
	} {
		require.Equal(t, v.expectedSlug, v.id.ServiceSlug())
		require.Equal(t, v.expectedResourceName, v.id.ResourceName())
	}
}

func TestAnyResourceRefEmpty(t *testing.T) {
	_, err := NewAnyResourceRef("")
	require.ErrorIs(t, err, reserrors.ErrPathIsEmpty)

	var ref AnyResourceRef

	_, err = ref.MarshalJSON()
	require.ErrorIs(t, err, reserrors.ErrPathIsEmpty)

	require.ErrorIs(t, ref.UnmarshalJSON([]byte(`""`)), reserrors.ErrPathIsEmpty)
}
