package imds_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/os/env"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.uber.org/mock/gomock"

	"go.mws.cloud/go-sdk/internal/imds"
	mockimds "go.mws.cloud/go-sdk/internal/imds/mocks"
)

func TestClient(t *testing.T) {
	type mockObjects struct {
		env        env.MapEnv
		httpClient *mockimds.MockHTTPClient
	}

	for _, v := range []struct {
		Name     string
		Prepare  func(*mockObjects)
		Expected string
		Error    error
	}{
		{
			Name: "invalid base url",
			Prepare: func(o *mockObjects) {
				o.env[metadataHostEnv] = "http://[fe80::1%en0]:8080/"
			},
			Error: url.EscapeError("%en"),
		},
		{
			Name: "request error",
			Prepare: func(o *mockObjects) {
				o.env[metadataHostEnv] = "otherHost"
				o.httpClient.EXPECT().Do(matchRequest("otherHost/computeMetadata/v1/key")).Return(nil, errRequest)
			},
			Error: errRequest,
		},
		{
			Name: "body read error",
			Prepare: func(o *mockObjects) {
				o.httpClient.EXPECT().Do(matchRequest("http://169.254.169.254/computeMetadata/v1/key")).Return(&http.Response{
					Body: io.NopCloser(failReader{}),
				}, nil)
			},
			Error: errRead,
		},
		{
			Name: "status not ok",
			Prepare: func(o *mockObjects) {
				o.httpClient.EXPECT().Do(matchRequest("http://169.254.169.254/computeMetadata/v1/key")).Return(&http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader("Key not found")),
				}, nil)
			},
			Error: imds.Error{
				Code:    http.StatusNotFound,
				Message: "Key not found",
			},
		},
		{
			Name: "body close error",
			Prepare: func(o *mockObjects) {
				o.httpClient.EXPECT().Do(matchRequest("http://169.254.169.254/computeMetadata/v1/key")).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       failCloser{strings.NewReader("value")},
				}, nil)
			},
			Error: errClose,
		},
		{
			Name: "status ok",
			Prepare: func(o *mockObjects) {
				o.httpClient.EXPECT().Do(matchRequest("http://169.254.169.254/computeMetadata/v1/key")).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"access_token": "token"}`)),
				}, nil)
			},
			Expected: `{"access_token": "token"}`,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mocks := &mockObjects{
				env:        env.MapEnv{},
				httpClient: mockimds.NewMockHTTPClient(ctrl),
			}
			v.Prepare(mocks)

			client := imds.NewClient(mocks.httpClient, mocks.env)
			actual, err := client.GetWithContext(t.Context(), "key")
			if v.Error != nil {
				require.ErrorIs(t, err, v.Error)
				return
			}

			require.NoError(t, err)
			require.Equal(t, v.Expected, actual)
		})
	}
}

const (
	metadataHostEnv = "MWS_COMPUTE_METADATA_HOST"

	errRequest = consterr.Error("request error")
	errRead    = consterr.Error("read error")
	errClose   = consterr.Error("close error")
)

type failReader struct{}

func (failReader) Read([]byte) (int, error) {
	return 0, errRead
}

type failCloser struct {
	io.Reader
}

func (failCloser) Close() error {
	return errClose
}

func matchRequest(addr string) gomock.Matcher {
	return gomock.Cond(func(in any) bool {
		req, ok := in.(*http.Request)
		if !ok {
			return false
		}

		return req.URL.String() == addr && req.Header.Get("Metadata-Flavor") == "Google"
	})
}
