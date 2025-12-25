package endpoints_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"go.mws.cloud/go-sdk/mws/endpoints"
	"go.mws.cloud/go-sdk/mws/endpoints/mocks"
)

const errResponse = consterr.Error("some response error")

func TestDiscoveryEndpointResolver(t *testing.T) {
	for _, tt := range []struct {
		Name        string
		Endpoints   endpoints.DiscoveryEndpoints
		Service     endpoints.ServiceName
		ResponseErr error
		ErrMsg      string
		Err         error
		Expected    endpoints.Endpoint
	}{
		{
			Name:        "endpoints_error",
			ResponseErr: errResponse,
			Err:         errResponse,
			ErrMsg:      "get endpoints list: some response error",
		},
		{
			Name:      "exists_compute",
			Endpoints: map[string]endpoints.DiscoveryEndpoint{"vpc": {"vpc.mws.ru"}, "compute": {"compute.mws.ru"}},
			Service:   "compute/instance",
			Expected:  "compute.mws.ru",
		},
		{
			Name:      "exists_some",
			Endpoints: map[string]endpoints.DiscoveryEndpoint{"some": {"thing"}},
			Service:   "some",
			Expected:  "thing",
		},
		{
			Name:      "exists_other",
			Endpoints: map[string]endpoints.DiscoveryEndpoint{"other": {"other.endpoint"}, "vpc": {"vpc.mws.ru"}},
			Service:   "other/end/point",
			Expected:  "other.endpoint",
		},
		{
			Name:      "nonexistent",
			Service:   "vpc/nonexistent",
			Endpoints: map[string]endpoints.DiscoveryEndpoint{"other": {"other.endpoint"}, "compute": {"compute.mws.ru"}},
			Err:       endpoints.ErrEndpointNotFound,
			ErrMsg:    "service \"vpc/nonexistent\": endpoint not found",
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			client := mocks.NewMockDiscoveryClient(ctrl)
			client.EXPECT().Endpoints(gomock.Any()).Return(tt.Endpoints, tt.ResponseErr)

			resolver := endpoints.NewDiscoveryServiceEndpointResolver(client)
			actual, err := resolver.Resolve(t.Context(), tt.Service)
			if tt.Err != nil {
				require.ErrorIs(t, err, tt.Err)
				require.ErrorContains(t, err, tt.ErrMsg)
				require.Empty(t, actual)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.Expected, actual)
		})
	}
}

func TestHTTPDiscoveryClient(t *testing.T) {
	for _, v := range []struct {
		Name           string
		ResponseBody   any
		ResponseError  error
		Expected       endpoints.DiscoveryEndpoints
		UnmarshalError *json.UnmarshalTypeError
	}{
		{
			Name:          "response error",
			ResponseError: errResponse,
		},
		{
			Name: "invalid json",
			ResponseBody: struct {
				Name string `json:"name"`
			}{
				Name: "test",
			},
			UnmarshalError: &json.UnmarshalTypeError{},
		},
		{
			Name: "valid json",
			ResponseBody: []struct {
				ID      string `json:"id"`
				Address string `json:"address"`
			}{
				{ID: "iam", Address: "iam.api.mws"},
				{ID: "compute", Address: "compute.api.mws"},
				{ID: "vpc", Address: "vpc.api.mws"},
			},
			Expected: map[string]endpoints.DiscoveryEndpoint{
				"vpc":     {Address: "vpc.api.mws"},
				"compute": {Address: "compute.api.mws"},
				"iam":     {Address: "iam.api.mws"},
			},
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			ctrl := gomock.NewController(t)
			httpClient := mocks.NewMockHTTPClient(ctrl)
			client := endpoints.NewHTTPDiscoveryClient(zap.NewNop(), httpClient, "")

			if v.ResponseError != nil {
				httpClient.EXPECT().Do(gomock.Any()).Return(nil, v.ResponseError)
				_, err := client.Endpoints(ctx)
				require.ErrorIs(t, err, v.ResponseError)
				return
			}

			data, err := json.Marshal(v.ResponseBody)
			require.NoError(t, err)

			response := &http.Response{Body: io.NopCloser(bytes.NewReader(data))}
			httpClient.EXPECT().Do(gomock.Any()).Return(response, nil)

			actual, err := client.Endpoints(ctx)
			if v.UnmarshalError != nil {
				require.ErrorAs(t, err, &v.UnmarshalError)
			} else {
				require.NoError(t, err)
				require.Equal(t, v.Expected, actual)
			}
		})
	}
}
