package models

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/pkg/resources/interfaces"
)

func TestAnyResourceRefMethods(t *testing.T) {
	for _, v := range []struct {
		id                   *AnyResourceRef
		expectedSlug         string
		expectedResourceName interfaces.ResourceName
	}{
		{nil, anyServiceSlug, ""},
		{ptr.Get(NewAnyResourceRef("")), anyServiceSlug, ""},
		{ptr.Get(NewAnyResourceRef("hello")), anyServiceSlug, "hello"},
		{ptr.Get(NewAnyResourceRef("/projects")), anyServiceSlug, "projects"},
		{ptr.Get(NewAnyResourceRef("/projects/project")), anyServiceSlug, "project"},
		{ptr.Get(NewAnyResourceRef("/projects/project/object")), anyServiceSlug, "object"},
		{ptr.Get(NewAnyResourceRef("some/projects/project")), anyServiceSlug, "project"},
		{ptr.Get(NewAnyResourceRef("some/projects/project/object")), anyServiceSlug, "object"},
		{ptr.Get(NewAnyResourceRef("some/projects/project/objects/object")), anyServiceSlug, "object"},
	} {
		require.Equal(t, v.expectedSlug, v.id.ServiceSlug())
		require.Equal(t, v.expectedResourceName, v.id.ResourceName())
	}
}
