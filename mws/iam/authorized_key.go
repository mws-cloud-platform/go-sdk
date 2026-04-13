package iam

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultClaimsTTL = time.Hour

// AuthorizedKey represents an authorized key.
type AuthorizedKey struct {
	Name           string
	PrivateKey     *ecdsa.PrivateKey
	ExpirationTime time.Time
	Algorithm      string
}

// NewAuthorizedKey generates a new authorized key.
func NewAuthorizedKey(name string) (AuthorizedKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return AuthorizedKey{}, err
	}
	return AuthorizedKey{
		Name:       name,
		PrivateKey: privateKey,
		Algorithm:  ES256,
	}, nil
}

// JWS creates a JWS token for the specified subject.
func (k AuthorizedKey) JWS(now func() time.Time, subject string) (string, error) {
	issuedAt := now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:  subject,
		IssuedAt: jwt.NewNumericDate(issuedAt),
	}

	expiresAt := k.ExpirationTime
	if expiresAt.IsZero() {
		expiresAt = issuedAt.Add(defaultClaimsTTL)
	}
	claims.ExpiresAt = jwt.NewNumericDate(expiresAt)

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = k.Name

	return token.SignedString(k.PrivateKey)
}

func (k AuthorizedKey) EncodePublicKey(base64Encoder *base64.Encoding) (string, error) {
	publicKey, err := x509.MarshalPKIXPublicKey(k.PrivateKey.Public())
	if err != nil {
		return "", fmt.Errorf("marshal public key: %w", err)
	}
	return base64Encoder.EncodeToString(publicKey), nil
}

func (k AuthorizedKey) EncodePrivateKey(base64Encoder *base64.Encoding) (string, error) {
	privateKey, err := x509.MarshalPKCS8PrivateKey(k.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("marshal private key: %w", err)
	}
	return base64Encoder.EncodeToString(privateKey), nil
}
