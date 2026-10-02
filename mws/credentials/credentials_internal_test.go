package credentials_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/mws/credentials"
)

func TestAccessToken(t *testing.T) {
	for _, tc := range []struct {
		name     string
		value    string
		redacted string
	}{
		{
			name:     "short",
			value:    "secret",
			redacted: "*******",
		},
		{
			name:     "long",
			value:    "prefix-" + strings.Repeat("x", 32),
			redacted: "prefix-x*******",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			token := credentials.NewAccessToken(tc.value)
			require.Equal(t, tc.redacted, token.String())
			require.Equal(t, tc.value, token.Value())
		})
	}
}

func requireCredentialsEqual(t *testing.T, expected, actual credentials.Credentials) {
	t.Helper()
	require.Equal(t, expected.AccessToken.Value(), actual.AccessToken.Value())
	require.Equal(t, expected.ExpiresAt, actual.ExpiresAt)
}
