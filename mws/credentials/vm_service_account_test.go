package credentials_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.uber.org/mock/gomock"

	"go.mws.cloud/go-sdk/mws/credentials"
	mockcredentials "go.mws.cloud/go-sdk/mws/credentials/mocks"
	"go.mws.cloud/go-sdk/pkg/clock/fakeclock"
)

func TestVMServiceAccountProvider(t *testing.T) {
	clock := fakeclock.NewFake()

	for _, v := range []struct {
		Name        string
		Prepare     func(*mockcredentials.MockMetadataProvider)
		Expected    credentials.Credentials
		ExpectedErr require.ErrorAssertionFunc
	}{
		{
			Name: "metadata provider error",
			Prepare: func(p *mockcredentials.MockMetadataProvider) {
				p.EXPECT().GetWithContext(gomock.Any(), gomock.Any()).Return("", errMetadataProvider)
			},
			ExpectedErr: errorIs(errMetadataProvider),
		},
		{
			Name: "unmarshal error",
			Prepare: func(p *mockcredentials.MockMetadataProvider) {
				p.EXPECT().GetWithContext(gomock.Any(), gomock.Any()).Return("hi", nil)
			},
			ExpectedErr: errorContains(`invalid character 'h' looking for beginning of value`),
		},
		{
			Name: "success",
			Prepare: func(p *mockcredentials.MockMetadataProvider) {
				p.EXPECT().GetWithContext(gomock.Any(), gomock.Any()).Return(`{"access_token": "token", "token_type": "Bearer", "expires_in": 3600}`, nil)
			},
			Expected: credentials.Credentials{
				AccessToken: "token",
				ExpiresAt:   clock.Now().Add(time.Hour),
			},
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()
			cache := credentials.NewDumbCache(credentials.WithClockCache(clock))

			ctrl := gomock.NewController(t)
			metadataProvider := mockcredentials.NewMockMetadataProvider(ctrl)
			provider := credentials.NewVMServiceAccountProvider(
				metadataProvider,
				credentials.WithVMServiceAccountProviderClock(clock),
				credentials.WithVMServiceAccountProviderCache(cache),
			)

			v.Prepare(metadataProvider)
			actual, err := provider.Provide(t.Context())
			if v.ExpectedErr != nil {
				v.ExpectedErr(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, v.Expected, actual)
		})
	}
}

func TestVMServiceAccountProviderCached(t *testing.T) {
	ctrl := gomock.NewController(t)
	metadataProvider := mockcredentials.NewMockMetadataProvider(ctrl)
	metadataProvider.EXPECT().GetWithContext(gomock.Any(), gomock.Any()).
		Return(`{"access_token": "token", "token_type": "Bearer", "expires_in": 3600}`, nil)

	clock := fakeclock.NewFake()
	cache := credentials.NewDumbCache(credentials.WithClockCache(clock))
	provider := credentials.NewVMServiceAccountProvider(
		metadataProvider,
		credentials.WithVMServiceAccountProviderClock(clock),
		credentials.WithVMServiceAccountProviderCache(cache),
	)

	const times = 5
	for range times {
		creds, err := provider.Provide(t.Context())
		require.NoError(t, err)
		require.Equal(t, "token", creds.AccessToken)
	}

	clock.Advance(time.Hour)
	metadataProvider.EXPECT().GetWithContext(gomock.Any(), gomock.Any()).
		Return(`{"access_token": "new_token", "token_type": "Bearer", "expires_in": 3600}`, nil)
	for range times {
		creds, err := provider.Provide(t.Context())
		require.NoError(t, err)
		require.Equal(t, "new_token", creds.AccessToken)
	}
}

func TestVMServiceAccountProviderNotCached(t *testing.T) {
	const times = 5

	ctrl := gomock.NewController(t)
	metadataProvider := mockcredentials.NewMockMetadataProvider(ctrl)
	metadataProvider.EXPECT().GetWithContext(gomock.Any(), gomock.Any()).
		Return(`{"access_token": "token", "token_type": "Bearer", "expires_in": 3600}`, nil).
		Times(times)

	clock := fakeclock.NewFake()
	provider := credentials.NewVMServiceAccountProvider(metadataProvider,
		credentials.WithVMServiceAccountProviderClock(clock),
		credentials.WithVMServiceAccountProviderCache(credentials.NewNoopCache()),
	)

	for range times {
		creds, err := provider.Provide(t.Context())
		require.NoError(t, err)
		require.Equal(t, "token", creds.AccessToken)
	}
}

func TestOnComputeVM(t *testing.T) {
	for _, v := range []struct {
		Name     string
		Result   string
		Error    error
		Expected bool
	}{
		{
			Name:     "error",
			Error:    consterr.Error("fail"),
			Expected: false,
		},
		{
			Name:     "empty",
			Expected: false,
		},
		{
			Name:     "success",
			Result:   "vm-name",
			Expected: true,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			metadataProvider := mockcredentials.NewMockMetadataProvider(ctrl)
			metadataProvider.EXPECT().GetWithContext(gomock.Any(), gomock.Any()).Return(v.Result, v.Error)

			actual := credentials.OnComputeVM(t.Context(), metadataProvider)
			require.Equal(t, v.Expected, actual)
		})
	}
}

const errMetadataProvider = consterr.Error("metadata provider error")

func errorIs(target error) require.ErrorAssertionFunc {
	return func(t require.TestingT, err error, args ...any) {
		require.ErrorIs(t, err, target, args...)
	}
}

func errorContains(s string) require.ErrorAssertionFunc {
	return func(t require.TestingT, err error, args ...any) {
		require.ErrorContains(t, err, s, args...)
	}
}
