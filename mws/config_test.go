package mws_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/os/env"

	"go.mws.cloud/go-sdk/mws"
)

func TestLoadConfig(t *testing.T) {
	for _, tc := range []struct {
		name     string
		opts     []mws.LoadConfigOption
		expected *mws.Config
		err      string
	}{
		{
			name: "default",
			expected: &mws.Config{
				BaseEndpoint: "https://api.mwsapis.ru",
				Zone:         "ru-central1-a",
				Timeout:      5 * time.Second,
			},
		},
		{
			name: "envs",
			opts: []mws.LoadConfigOption{
				mws.LoadConfigWithEnv(env.MapEnv{
					"MWS_BASE_ENDPOINT": "https://www.example.com",
					"MWS_PROJECT":       "my-project",
					"MWS_ZONE":          "ru-central1-b",
					"MWS_TOKEN":         "my-token",
					"MWS_SERVICE_ACCOUNT_AUTHORIZED_KEY_PATH": "/path/to/key.json",
					"MWS_TIMEOUT": "1m",
				}),
			},
			expected: &mws.Config{
				BaseEndpoint:                    "https://www.example.com",
				Project:                         "my-project",
				Zone:                            "ru-central1-b",
				Token:                           "my-token",
				ServiceAccountAuthorizedKeyPath: "/path/to/key.json",
				Timeout:                         time.Minute,
			},
		},
		{
			name: "invalid_timeout",
			opts: []mws.LoadConfigOption{
				mws.LoadConfigWithEnv(env.MapEnv{
					"MWS_TIMEOUT": "invalid",
				}),
			},
			err: `parse "MWS_TIMEOUT": time: invalid duration "invalid"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := mws.LoadConfig(tc.opts...)
			if tc.err != "" {
				require.EqualError(t, err, tc.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}
