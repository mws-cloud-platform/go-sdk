package models

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/resources/interfaces"
)

func TestAnyResourceIDMethods(t *testing.T) {
	for _, v := range []struct {
		id                   *AnyResourceID
		expectedSlug         string
		expectedResourceName interfaces.ResourceName
	}{
		{nil, anyServiceSlug, ""},
		{new(NewAnyResourceID("")), anyServiceSlug, ""},
		{new(NewAnyResourceID("hello")), anyServiceSlug, "hello"},
		{new(NewAnyResourceID("/projects")), anyServiceSlug, "projects"},
		{new(NewAnyResourceID("/projects/project")), anyServiceSlug, "project"},
		{new(NewAnyResourceID("/projects/project/object")), anyServiceSlug, "object"},
		{new(NewAnyResourceID("some/projects/project")), "some", "project"},
		{new(NewAnyResourceID("some/projects/project/object")), "some", "object"},
		{new(NewAnyResourceID("some/projects/project/objects/object")), "some", "object"},
	} {
		require.Equal(t, v.expectedSlug, v.id.ServiceSlug())
		require.Equal(t, v.expectedResourceName, v.id.ResourceName())
	}
}
