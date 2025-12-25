package iam

import (
	"crypto/ecdsa"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
	"gopkg.in/yaml.v3"
)

//go:embed testdata/private.key
var privateKey string

func TestAuthorizedKey_JWS(t *testing.T) {
	authorizedKey, err := NewAuthorizedKey("test")
	require.NoError(t, err)

	now := func() time.Time {
		return time.Date(2023, time.February, 14, 15, 30, 0, 0, time.UTC)
	}

	jws, err := authorizedKey.JWS(now, "projects/projectTest/serviceAccounts/saTest")
	require.NoError(t, err)

	parts := strings.Split(jws, ".")
	require.Len(t, parts, 3)

	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	require.NoError(t, err)
	assert.JSONEq(t, `{"alg":"ES256","kid":"test","typ":"JWT"}`, string(header))

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)
	assert.JSONEq(t, `{"sub":"projects/projectTest/serviceAccounts/saTest","iat":1676388600,"exp":1676392200}`, string(payload))

	sign, err := base64.RawURLEncoding.DecodeString(parts[2])
	require.NoError(t, err)

	publicKey := authorizedKey.PrivateKey.Public()
	err = jwt.SigningMethodES256.Verify(parts[0]+"."+parts[1], sign, publicKey)
	require.NoError(t, err)
}

func TestAuthorizedKey_MarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/authorizedKeyJSON/golden"),
		golden.WithRecreateOnUpdate())

	privateKey := getTestPrivateKey(t)

	t.Run("standard", func(t *testing.T) {
		key := AuthorizedKey{
			ID:             "test",
			PrivateKey:     privateKey,
			Algorithm:      ES256,
			ExpirationTime: time.Date(2026, time.February, 14, 15, 30, 0, 0, time.UTC),
		}

		actual, err := json.Marshal(key)
		require.NoError(t, err)
		dir.String(t, "key.json", string(actual))
	})

	t.Run("no_expiration", func(t *testing.T) {
		key := AuthorizedKey{
			ID:         "test",
			PrivateKey: privateKey,
			Algorithm:  ES256,
		}

		actual, err := json.Marshal(key)
		require.NoError(t, err)
		dir.String(t, "key_no_expiration.json", string(actual))
	})
}

func TestAuthorizedKey_UnmarshalJSON(t *testing.T) {
	privateKey := getTestPrivateKey(t)

	t.Run("standard", func(t *testing.T) {
		expected := AuthorizedKey{
			ID:             "test",
			PrivateKey:     privateKey,
			ExpirationTime: time.Date(2026, time.February, 14, 15, 30, 0, 0, time.UTC),
			Algorithm:      ES256,
		}

		data, err := testdata.ReadFile("testdata/authorizedKeyJSON/golden/key.json")
		require.NoError(t, err)

		var actual AuthorizedKey
		err = json.Unmarshal(data, &actual)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("no_expiration", func(t *testing.T) {
		expected := AuthorizedKey{
			ID:         "test",
			PrivateKey: privateKey,
			Algorithm:  ES256,
		}

		data, err := testdata.ReadFile("testdata/authorizedKeyJSON/golden/key_no_expiration.json")
		require.NoError(t, err)

		var actual AuthorizedKey
		err = json.Unmarshal(data, &actual)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestAuthorizedKey_UnmarshalJSON_invalid(t *testing.T) {
	data, err := testdata.ReadFile("testdata/authorizedKeyJSON/invalid_private_key.json")
	require.NoError(t, err)

	var actual AuthorizedKey
	err = json.Unmarshal(data, &actual)
	require.Error(t, err)
}

func TestAuthorizedKey_MarshalYAML(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/authorizedKeyYAML/golden"),
		golden.WithRecreateOnUpdate())

	privateKey := getTestPrivateKey(t)

	t.Run("standard", func(t *testing.T) {
		key := AuthorizedKey{
			ID:             "test",
			PrivateKey:     privateKey,
			ExpirationTime: time.Date(2026, time.February, 14, 15, 30, 0, 0, time.UTC),
			Algorithm:      ES256,
		}

		actual, err := yaml.Marshal(key)
		require.NoError(t, err)
		dir.String(t, "key.yaml", string(actual))
	})

	t.Run("no_expiration", func(t *testing.T) {
		key := AuthorizedKey{
			ID:         "test",
			PrivateKey: privateKey,
			Algorithm:  ES256,
		}

		actual, err := yaml.Marshal(key)
		require.NoError(t, err)
		dir.String(t, "key_no_expiration.yaml", string(actual))
	})
}

func TestAuthorizedKey_UnmarshalYAML(t *testing.T) {
	privateKey := getTestPrivateKey(t)

	t.Run("standard", func(t *testing.T) {
		expected := AuthorizedKey{
			ID:             "test",
			PrivateKey:     privateKey,
			ExpirationTime: time.Date(2026, time.February, 14, 15, 30, 0, 0, time.UTC),
			Algorithm:      ES256,
		}

		data, err := testdata.ReadFile("testdata/authorizedKeyYAML/golden/key.yaml")
		require.NoError(t, err)

		var actual AuthorizedKey
		err = yaml.Unmarshal(data, &actual)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("no_expiration", func(t *testing.T) {
		expected := AuthorizedKey{
			ID:         "test",
			PrivateKey: privateKey,
			Algorithm:  ES256,
		}

		data, err := testdata.ReadFile("testdata/authorizedKeyYAML/golden/key_no_expiration.yaml")
		require.NoError(t, err)

		var actual AuthorizedKey
		err = yaml.Unmarshal(data, &actual)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestAuthorizedKey_UnmarshalYAML_invalid(t *testing.T) {
	data, err := testdata.ReadFile("testdata/authorizedKeyYAML/invalid_private_key.yaml")
	require.NoError(t, err)

	var actual AuthorizedKey
	err = yaml.Unmarshal(data, &actual)
	require.Error(t, err)
}

func getTestPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	key, err := ParsePrivateKey(privateKey, base64.StdEncoding)
	require.NoError(t, err)
	return key
}
