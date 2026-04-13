package parsers

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/pkg/context/values"
)

var (
	// testService/foos/fo/{foo}/projects/{project}/networks/{network}
	template1 = Template{
		{
			Value:       "network",
			IsConstant:  false,
			SearchAfter: false,
		},
		{
			Value:       "networks",
			IsConstant:  true,
			SearchAfter: false,
		},
		{
			Value:       "project",
			IsConstant:  false,
			SearchAfter: true,
		},
		{
			Value:       "projects",
			IsConstant:  true,
			SearchAfter: false,
		},
		{
			Value:       "foo",
			IsConstant:  false,
			SearchAfter: true,
		},
		{
			Value:       "fo",
			IsConstant:  true,
			SearchAfter: true,
		},
		{
			Value:       "foos",
			IsConstant:  true,
			SearchAfter: false,
		},
		{
			Value:       "testService",
			IsConstant:  true,
			SearchAfter: false,
		},
	}

	// testService/projects/{project}/quotas/serviceAccountCount
	template2 = Template{
		{
			Value:       "serviceAccountCount",
			IsConstant:  true,
			SearchAfter: true,
		},
		{
			Value:       "quotas",
			IsConstant:  true,
			SearchAfter: false,
		},
		{
			Value:       "project",
			IsConstant:  false,
			SearchAfter: true,
		},
		{
			Value:       "projects",
			IsConstant:  true,
			SearchAfter: false,
		},
		{
			Value:       "testService",
			IsConstant:  true,
			SearchAfter: false,
		},
	}

	// testService/{scopeType}/{scopeId}/discounts
	template3 = Template{
		{
			Value:       "discounts",
			IsConstant:  true,
			SearchAfter: false,
		},
		{
			Value:       "scopeId",
			IsConstant:  false,
			SearchAfter: false,
		},
		{
			Value:       "scopeType",
			IsConstant:  false,
			SearchAfter: false,
		},
		{
			Value:       "testService",
			IsConstant:  true,
			SearchAfter: false,
		},
	}
)

func TestReference(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		path        string
		template    Template
		ctxFillFunc func(ctx context.Context) context.Context
		result      map[string]string
		err         error
	}{
		{
			name:     "full reference string without service",
			path:     "foos/fo/fooNAME/projects/projectNAME/networks/networkNAME",
			template: template1,
			result: map[string]string{
				"network": "networkNAME",
				"project": "projectNAME",
				"foo":     "fooNAME",
			},
			err: nil,
		},
		{
			name:     "full reference string with service",
			path:     "testService/foos/fo/fooNAME/projects/projectNAME/networks/networkNAME",
			template: template1,
			result: map[string]string{
				"network": "networkNAME",
				"project": "projectNAME",
				"foo":     "fooNAME",
			},
			err: nil,
		},
		{
			name:     "full reference string with invalid service part",
			path:     "/foos/fo/fooNAME/projects/projectNAME/networks/networkNAME",
			template: template1,
			result:   nil,
			err:      ErrReferenceParsing,
		},
		{
			name:     "reference string too long",
			path:     "/testService/foos/fo/fooNAME/projects/projectNAME/networks/networkNAME",
			template: template1,
			result:   nil,
			err:      ErrReferenceParsing,
		},
		{
			name:     "reference string with empty parts",
			path:     "foos/fo//projects//networks/networkNAME",
			template: template1,
			result:   nil,
			err:      ErrReferenceParsing,
		},
		{
			name:     "reference string is short",
			path:     "fo/fooNAME/projects/projectNAME/networks/networkNAME",
			template: template1,
			result:   nil,
			err:      ErrReferenceParsing,
		},
		{
			name:     "restoring reference from context",
			path:     "networks/networkNAME",
			template: template1,
			ctxFillFunc: func(ctx context.Context) context.Context {
				ctx = values.With(ctx, "project", "projectNAME")
				ctx = values.With(ctx, "foo", "fooNAME")
				return ctx
			},
			result: map[string]string{
				"network": "networkNAME",
				"project": "projectNAME",
				"foo":     "fooNAME",
			},
			err: nil,
		},
		{
			name:     "error restoring reference from context",
			path:     "networks/networkNAME",
			template: template1,
			ctxFillFunc: func(ctx context.Context) context.Context {
				ctx = values.With(ctx, "foo", "fooNAME")
				return ctx
			},
			result: nil,
			err:    ErrReferenceParsing,
		},
		{
			name:     "restoring short reference from context",
			path:     "networkNAME",
			template: template1,
			ctxFillFunc: func(ctx context.Context) context.Context {
				ctx = values.With(ctx, "project", "projectNAME")
				ctx = values.With(ctx, "foo", "fooNAME")
				return ctx
			},
			result: map[string]string{
				"network": "networkNAME",
				"project": "projectNAME",
				"foo":     "fooNAME",
			},
			err: nil,
		},
		{
			name:     "empty reference",
			path:     "",
			template: template1,
			ctxFillFunc: func(ctx context.Context) context.Context {
				ctx = values.With(ctx, "project", "projectNAME")
				ctx = values.With(ctx, "foo", "fooNAME")
				return ctx
			},
			result: nil,
			err:    ErrReferenceParsing,
		},
		{
			name:     "invalid short reference",
			path:     "networksfff/networkNAME",
			template: template1,
			ctxFillFunc: func(ctx context.Context) context.Context {
				ctx = values.With(ctx, "project", "projectNAME")
				ctx = values.With(ctx, "foo", "fooNAME")
				return ctx
			},
			result: nil,
			err:    ErrReferenceParsing,
		},
		{
			name:     "full reference without bearing id",
			path:     "testService/projects/projectNAME/quotas/serviceAccountCount",
			template: template2,
			result: map[string]string{
				"project": "projectNAME",
			},
			err: nil,
		},
		{
			name:     "short reference without bearing id",
			path:     "quotas/serviceAccountCount",
			template: template2,
			ctxFillFunc: func(ctx context.Context) context.Context {
				ctx = values.With(ctx, "project", "projectNAME")
				return ctx
			},
			result: map[string]string{
				"project": "projectNAME",
			},
			err: nil,
		},
		{
			name:     "too short reference without bearing id",
			path:     "serviceAccountCount",
			template: template2,
			ctxFillFunc: func(ctx context.Context) context.Context {
				ctx = values.With(ctx, "project", "projectNAME")
				return ctx
			},
			err: ErrReferenceParsing,
		},
		{
			name:     "full reference parses to scope params",
			path:     "testService/type/id/discounts",
			template: template3,
			result: map[string]string{
				"scopeType": "type",
				"scopeId":   "id",
			},
			err: nil,
		},
		{
			name:     "short reference without context fails",
			path:     "discounts",
			template: template3,
			err:      ErrReferenceParsing,
		},
		{
			name:     "short reference with context succeeds",
			path:     "discounts",
			template: template3,
			ctxFillFunc: func(ctx context.Context) context.Context {
				ctx = values.With(ctx, "scopeType", "type")
				ctx = values.With(ctx, "scopeId", "id")
				return ctx
			},
			result: map[string]string{
				"scopeType": "type",
				"scopeId":   "id",
			},
			err: nil,
		},
		{
			name:     "partial reference with context succeeds",
			path:     "id/discounts",
			template: template3,
			ctxFillFunc: func(ctx context.Context) context.Context {
				ctx = values.With(ctx, "scopeType", "type")
				return ctx
			},
			result: map[string]string{
				"scopeType": "type",
				"scopeId":   "id",
			},
			err: nil,
		},
		{
			name:     "full reference without slug succeeds",
			path:     "type/id/discounts",
			template: template3,
			result: map[string]string{
				"scopeType": "type",
				"scopeId":   "id",
			},
			err: nil,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := t.Context()
			if testCase.ctxFillFunc != nil {
				ctx = testCase.ctxFillFunc(ctx)
			}

			result, err := Reference(ctx, testCase.path, testCase.template)
			if testCase.err != nil {
				require.ErrorIs(t, err, testCase.err)
			} else {
				require.Equal(t, testCase.result, result)
			}
		})
	}
}

func TestTemplate_String(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		template Template
		result   string
	}{
		{
			name: "ok",
			template: Template{
				{
					Value:      "world",
					IsConstant: false,
				},
				{
					Value:      "hello",
					IsConstant: true,
				},
			},
			result: "hello/{world}",
		},
		{
			name:     "empty",
			template: Template{},
			result:   "",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.result, testCase.template.String())
		})
	}
}

func TestTemplate_AsID(t *testing.T) {
	idTemplate := template1.AsID()

	for index, pathElement := range idTemplate {
		require.True(t, pathElement.SearchAfter)
		require.Equal(t, template1[index].Value, pathElement.Value)
		require.Equal(t, template1[index].IsConstant, pathElement.IsConstant)
	}
}

func TestTemplate_Pattern(t *testing.T) {
	t.Parallel()

	template := Template{
		{
			Value:       "disk",
			Pattern:     regexp.MustCompile("^[a-z]([a-z0-9-_]{0,61}[a-z0-9])?$"),
			IsConstant:  false,
			SearchAfter: false,
		},
		{
			Value:       "disks",
			IsConstant:  true,
			SearchAfter: false,
		},
		{
			Value:       "project",
			IsConstant:  false,
			SearchAfter: true,
		},
		{
			Value:       "projects",
			IsConstant:  false,
			SearchAfter: true,
		},
	}
	for _, testCase := range []struct {
		name string
		path string
		err  error
	}{
		{
			name: "empty",
			path: "",
			err:  ErrReferenceParsing,
		},
		{
			name: "success",
			path: "projects/project123/disks/disk123",
		},
		{
			name: "error",
			path: "projects/project123/disks/Disk123",
			err:  ErrReferenceParsing,
		},
	} {
		t.Run(testCase.name, func(tt *testing.T) {
			tt.Parallel()

			_, err := Reference(tt.Context(), testCase.path, template)
			if testCase.err != nil {
				require.Error(tt, err, testCase.err)
			} else {
				require.NoError(tt, err)
			}
		})
	}
}
