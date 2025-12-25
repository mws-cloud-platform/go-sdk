package models

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/pkg/resources/interfaces"
)

func TestAnyResourceIDMethods(t *testing.T) {
	for _, v := range []struct {
		id                   *AnyResourceID
		expectedSlug         string
		expectedResourceName interfaces.ResourceName
	}{
		{nil, anyServiceSlug, ""},
		{ptr.Get(NewAnyResourceID("")), anyServiceSlug, ""},
		{ptr.Get(NewAnyResourceID("hello")), anyServiceSlug, "hello"},
		{ptr.Get(NewAnyResourceID("/projects")), anyServiceSlug, "projects"},
		{ptr.Get(NewAnyResourceID("/projects/project")), anyServiceSlug, "project"},
		{ptr.Get(NewAnyResourceID("/projects/project/object")), anyServiceSlug, "object"},
		{ptr.Get(NewAnyResourceID("some/projects/project")), "some", "project"},
		{ptr.Get(NewAnyResourceID("some/projects/project/object")), "some", "object"},
		{ptr.Get(NewAnyResourceID("some/projects/project/objects/object")), "some", "object"},
	} {
		require.Equal(t, v.expectedSlug, v.id.ServiceSlug())
		require.Equal(t, v.expectedResourceName, v.id.ResourceName())
	}
}
