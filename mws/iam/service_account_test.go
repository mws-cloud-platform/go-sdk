package iam

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceAccountKey_UnmarshalJSON(t *testing.T) {
	actual := ServiceAccountKey{}

	data, err := testdata.ReadFile("testdata/serviceAccountKey/key.json")
	require.NoError(t, err)

	err = json.Unmarshal(data, &actual)
	require.NoError(t, err)

	require.Equal(t, ServiceAccount{
		Project: "foo",
		Name:    "bar",
	}, actual.ServiceAccount)
	require.Equal(t, "default", actual.AuthorizedKey.ID)
	require.Equal(t, getTestPrivateKey(t), actual.AuthorizedKey.PrivateKey)
	require.Equal(t, "2026-01-01 00:00:00 +0000 UTC", actual.AuthorizedKey.ExpirationTime.String())
	require.Equal(t, ES256, actual.AuthorizedKey.Algorithm)
}

func TestServiceAccountKey_UnmarshalJSON_noExpiration(t *testing.T) {
	actual := ServiceAccountKey{}

	data, err := testdata.ReadFile("testdata/serviceAccountKey/key_no_expiration.json")
	require.NoError(t, err)

	err = json.Unmarshal(data, &actual)
	require.NoError(t, err)

	require.Equal(t, ServiceAccount{
		Project: "foo",
		Name:    "bar",
	}, actual.ServiceAccount)
	require.Equal(t, "default", actual.AuthorizedKey.ID)
	require.Equal(t, getTestPrivateKey(t), actual.AuthorizedKey.PrivateKey)
	require.True(t, actual.AuthorizedKey.ExpirationTime.IsZero())
	require.Equal(t, ES256, actual.AuthorizedKey.Algorithm)
}
