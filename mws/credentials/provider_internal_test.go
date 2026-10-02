package credentials

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"go.mws.cloud/go-sdk/pkg/clock/fakeclock"
)

func TestProvider_do(t *testing.T) {
	now := time.Date(2025, 9, 4, 12, 1, 5, 0, time.UTC)
	cl := fakeclock.NewFake(fakeclock.WithStartAt(now))
	actualCreds := Credentials{
		ExpiresAt: now.Add(DefaultTokenExpirationDelta * 5),
	}
	expiringCreds := Credentials{
		ExpiresAt: now.Add(DefaultTokenExpirationDelta - time.Minute),
	}
	expiredCreds := Credentials{
		ExpiresAt: now,
	}

	for _, testCase := range []struct {
		name             string
		credentialsCache *Credentials
		mockFunc         func() (Credentials, error)
		expectedCreds    Credentials
		expectedError    bool
	}{
		{
			name:             "cache hit, valid credentials",
			credentialsCache: &actualCreds,
			mockFunc: func() (Credentials, error) {
				t.Fatal("func should not be called for valid cached credentials")
				return Credentials{}, nil
			},
			expectedCreds: actualCreds,
			expectedError: false,
		},
		{
			name:             "credentials with expiration delta has expired, cache hit, service error",
			credentialsCache: &expiringCreds,
			mockFunc: func() (Credentials, error) {
				return Credentials{}, fmt.Errorf("service failed")
			},
			expectedCreds: expiringCreds,
			expectedError: false,
		},
		{
			name:             "credentials has expired, cache hit, service error",
			credentialsCache: &expiredCreds,
			mockFunc: func() (Credentials, error) {
				return Credentials{}, fmt.Errorf("service failed")
			},
			expectedError: true,
		},
		{
			name: "cache miss, service ok",
			mockFunc: func() (Credentials, error) {
				return actualCreds, nil
			},
			expectedCreds: actualCreds,
			expectedError: false,
		},
		{
			name: "cache miss, service error",
			mockFunc: func() (Credentials, error) {
				return Credentials{}, fmt.Errorf("service failed")
			},
			expectedError: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			cache := NewDumbCache(WithClockCache(cl))
			key := "key"

			if testCase.credentialsCache != nil {
				err := cache.Store(key, *testCase.credentialsCache)
				require.NoError(t, err)
			}

			p := &provider{
				clock:  cl,
				logger: zap.NewNop(),
				cache:  cache,
			}
			creds, err := p.do(key, testCase.mockFunc)
			if testCase.expectedError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, testCase.expectedCreds.AccessToken.Value(), creds.AccessToken.Value())
			require.Equal(t, testCase.expectedCreds.ExpiresAt, creds.ExpiresAt)
		})
	}
}
