package imds_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/os/env"
	"go.uber.org/mock/gomock"

	"go.mws.cloud/go-sdk/internal/imds"
	mockimds "go.mws.cloud/go-sdk/internal/imds/mocks"
)

func TestGetVMServiceAccount(t *testing.T) {
	e := env.MapEnv{
		"MWS_COMPUTE_METADATA_HOST": "host",
	}
	for _, v := range []struct {
		Name     string
		Prepare  func(*mockimds.MockHTTPClient)
		Expected string
		Error    error
	}{
		{
			Name: "error",
			Prepare: func(client *mockimds.MockHTTPClient) {
				client.EXPECT().Do(matchRequest("host/computeMetadata/v1/instance/service-accounts/?recursive=true")).Return(nil, errRequest)
			},
			Error: errRequest,
		},
		{
			Name: "status not ok",
			Prepare: func(client *mockimds.MockHTTPClient) {
				client.EXPECT().Do(matchRequest("host/computeMetadata/v1/instance/service-accounts/?recursive=true")).Return(&http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader("bad request")),
				}, nil)
			},
			Error: imds.Error{Code: http.StatusBadRequest, Message: "bad request"},
		},
		{
			Name: "no sa",
			Prepare: func(client *mockimds.MockHTTPClient) {
				client.EXPECT().Do(matchRequest("host/computeMetadata/v1/instance/service-accounts/?recursive=true")).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("null")),
				}, nil)
			},
			Error: imds.ErrServiceAccountNotFound,
		},
		{
			Name: "success",
			Prepare: func(client *mockimds.MockHTTPClient) {
				client.EXPECT().Do(matchRequest("host/computeMetadata/v1/instance/service-accounts/?recursive=true")).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"serviceAccounts/testServiceAccount": 42}`)),
				}, nil)
			},
			Expected: "serviceAccounts/testServiceAccount",
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			client := mockimds.NewMockHTTPClient(ctrl)
			v.Prepare(client)

			actual, err := imds.GetVMServiceAccount(t.Context(), client, e)
			if err != nil {
				require.ErrorIs(t, err, v.Error)
				return
			}

			require.NoError(t, err)
			require.Equal(t, v.Expected, actual)
		})
	}
}
