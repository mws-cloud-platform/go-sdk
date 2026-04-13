package credentials_test

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"go.mws.cloud/go-sdk/mws/credentials"
	"go.mws.cloud/go-sdk/mws/iam"
	"go.mws.cloud/go-sdk/pkg/clock/fakeclock"
	"go.mws.cloud/go-sdk/service/iam/client"
	mockclient "go.mws.cloud/go-sdk/service/iam/client/mocks"
	"go.mws.cloud/go-sdk/service/iam/model"
)

func TestServiceAccountProvider(t *testing.T) {
	ctrl := gomock.NewController(t)

	serviceAccount := iam.ServiceAccount{
		Project: "project",
		Name:    "test",
	}
	serviceAccountAuthorizedKey := iam.ServiceAccountAuthorizedKey{
		ServiceAccount: serviceAccount,
		AuthorizedKey: iam.AuthorizedKey{
			Name:       "key",
			PrivateKey: getTestPrivateKey(t),
		},
	}
	token := "token"
	now := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	expected := credentials.Credentials{
		AccessToken: token,
		ExpiresAt:   now.Add(credentials.DefaultTokenTTL),
	}

	serviceAccountTokenIssuer := mockclient.NewMockIssueServiceAccountToken(ctrl)
	serviceAccountTokenIssuer.EXPECT().
		IssueServiceAccountTokenV2(
			gomock.Any(),
			gomock.Cond(matchIssueServiceAccountTokenRequest(serviceAccount.String())),
		).
		Return(&client.IssueServiceAccountTokenV2Response{
			Code:        http.StatusOK,
			Response200: &model.SuccessTokenV2Response{AccessToken: token},
		}, nil)

	clock := fakeclock.NewFake(fakeclock.WithStartAt(now))
	provider := credentials.NewServiceAccountProvider(
		serviceAccountAuthorizedKey,
		serviceAccountTokenIssuer,
		credentials.WithServiceAccountProviderClock(clock),
		credentials.WithServiceAccountProviderCache(credentials.NewNoopCache()),
	)

	actual, err := provider.Provide(t.Context())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func matchIssueServiceAccountTokenRequest(serviceAccount string) func(any) bool {
	return func(x any) bool {
		req, ok := x.(client.IssueServiceAccountTokenV2Request)
		if !ok {
			return false
		}
		return req.Authorization != nil && *req.Authorization != "" &&
			*req.ServiceAccount == serviceAccount
	}
}

func getTestPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	privateKeyRaw, err := os.ReadFile("testdata/private.key")
	require.NoError(t, err)

	decoded, err := base64.StdEncoding.DecodeString(string(privateKeyRaw))
	require.NoError(t, err)

	privateKey, err := x509.ParsePKCS8PrivateKey(decoded)
	require.NoError(t, err)

	return privateKey.(*ecdsa.PrivateKey)
}
