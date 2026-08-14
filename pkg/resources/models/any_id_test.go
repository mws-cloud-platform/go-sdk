package models

import (
	"testing"

	"github.com/stretchr/testify/require"

	reserrors "go.mws.cloud/go-sdk/internal/resources/errors"
	"go.mws.cloud/go-sdk/pkg/resources/interfaces"
)

func TestAnyResourceIDMethods(t *testing.T) {
	for _, v := range []struct {
		id                   *AnyResourceID
		expectedSlug         string
		expectedResourceName interfaces.ResourceName
	}{
		{nil, anyServiceSlug, ""},
		{new(NewMustAnyResourceID("hello")), anyServiceSlug, "hello"},
		{new(NewMustAnyResourceID("/projects")), anyServiceSlug, "projects"},
		{new(NewMustAnyResourceID("/projects/project")), anyServiceSlug, "project"},
		{new(NewMustAnyResourceID("/projects/project/object")), anyServiceSlug, "object"},
		{new(NewMustAnyResourceID("some/projects/project")), "some", "project"},
		{new(NewMustAnyResourceID("some/projects/project/object")), "some", "object"},
		{new(NewMustAnyResourceID("some/projects/project/objects/object")), "some", "object"},
	} {
		require.Equal(t, v.expectedSlug, v.id.ServiceSlug())
		require.Equal(t, v.expectedResourceName, v.id.ResourceName())
	}
}

func TestAnyResourceIDEmpty(t *testing.T) {
	_, err := NewAnyResourceID("")
	require.ErrorIs(t, err, reserrors.ErrIDIsEmpty)

	var id AnyResourceID

	_, err = id.MarshalJSON()
	require.ErrorIs(t, err, reserrors.ErrIDIsEmpty)

	require.ErrorIs(t, id.UnmarshalJSON([]byte(`""`)), reserrors.ErrIDIsEmpty)
}
