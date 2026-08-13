package models

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/resources/interfaces"
)

func TestAnyResourceRefMethods(t *testing.T) {
	for _, v := range []struct {
		id                   *AnyResourceRef
		expectedSlug         string
		expectedResourceName interfaces.ResourceName
	}{
		{nil, anyServiceSlug, ""},
		{new(NewAnyResourceRef("")), anyServiceSlug, ""},
		{new(NewAnyResourceRef("hello")), anyServiceSlug, "hello"},
		{new(NewAnyResourceRef("/projects")), anyServiceSlug, "projects"},
		{new(NewAnyResourceRef("/projects/project")), anyServiceSlug, "project"},
		{new(NewAnyResourceRef("/projects/project/object")), anyServiceSlug, "object"},
		{new(NewAnyResourceRef("some/projects/project")), anyServiceSlug, "project"},
		{new(NewAnyResourceRef("some/projects/project/object")), anyServiceSlug, "object"},
		{new(NewAnyResourceRef("some/projects/project/objects/object")), anyServiceSlug, "object"},
	} {
		require.Equal(t, v.expectedSlug, v.id.ServiceSlug())
		require.Equal(t, v.expectedResourceName, v.id.ResourceName())
	}
}
