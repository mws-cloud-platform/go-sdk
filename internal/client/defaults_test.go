package client_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/internal/client"
)

func TestDefaultsInjector(t *testing.T) {
	tests := []struct {
		name    string
		req     *requestWithDefaults
		project string
		zone    string
		want    *requestWithDefaults
	}{
		{
			name:    "empty",
			req:     &requestWithDefaults{},
			project: "",
			zone:    "",
			want:    &requestWithDefaults{},
		},
		{
			name:    "only defaults",
			req:     &requestWithDefaults{},
			project: "project",
			zone:    "zone",
			want:    &requestWithDefaults{Project: "project", Zone: "zone"},
		},
		{
			name:    "already set without defaults",
			req:     &requestWithDefaults{Project: "project", Zone: "zone"},
			project: "",
			zone:    "",
			want:    &requestWithDefaults{Project: "project", Zone: "zone"},
		},
		{
			name:    "already set with defaults",
			req:     &requestWithDefaults{Project: "project", Zone: "zone"},
			project: "foo",
			zone:    "bar",
			want:    &requestWithDefaults{Project: "project", Zone: "zone"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			injector := client.DefaultsInjector(tt.project, tt.zone)
			_ = injector(t.Context(), tt.req, nil, noopInvoke)
			require.Equal(t, tt.want, tt.req)
		})
	}
}

type requestWithDefaults struct {
	Project string
	Zone    string
}

func (r *requestWithDefaults) SetProject(project string) {
	r.Project = project
}

func (r *requestWithDefaults) GetProject() string {
	return r.Project
}

func (r *requestWithDefaults) SetZone(zone string) {
	r.Zone = zone
}

func (r *requestWithDefaults) GetZone() string {
	return r.Zone
}
