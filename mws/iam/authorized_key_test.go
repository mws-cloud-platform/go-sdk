package iam

import (
	"crypto/ecdsa"
	_ "embed"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func getTestPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	key, err := ParsePrivateKey(privateKey, base64.StdEncoding)
	require.NoError(t, err)
	return key
}
