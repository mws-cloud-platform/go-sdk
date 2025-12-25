package iam

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
)

const defaultClaimsTTL = time.Hour

// AuthorizedKey represents an authorized key.
type AuthorizedKey struct {
	ID             string
	PrivateKey     *ecdsa.PrivateKey
	ExpirationTime time.Time
	Algorithm      string
}

// NewAuthorizedKey generates a new authorized key.
func NewAuthorizedKey(id string) (AuthorizedKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return AuthorizedKey{}, err
	}
	return AuthorizedKey{
		ID:         id,
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
	token.Header["kid"] = k.ID

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

func (k AuthorizedKey) MarshalJSON() ([]byte, error) {
	v, err := newAuthorizedKeyOut(k)
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

func (k *AuthorizedKey) UnmarshalJSON(data []byte) error {
	var v authorizedKeyIn
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return v.toAuthorizedKey(k)
}

func (k AuthorizedKey) MarshalYAML() (any, error) {
	return newAuthorizedKeyOut(k)
}

func (k *AuthorizedKey) UnmarshalYAML(node *yaml.Node) error {
	var v authorizedKeyIn
	if err := node.Decode(&v); err != nil {
		return err
	}
	return v.toAuthorizedKey(k)
}

type authorizedKeyIn struct {
	ID             string `json:"id" yaml:"id"`
	PrivateKey     string `json:"private_key" yaml:"private_key"`
	ExpirationTime string `json:"expiration_time" yaml:"expiration_time"`
	Algorithm      string `json:"algorithm" yaml:"algorithm"`
}

func (k *authorizedKeyIn) toAuthorizedKey(key *AuthorizedKey) (err error) {
	key.ID = k.ID
	key.PrivateKey, err = ParsePrivateKey(k.PrivateKey, base64.StdEncoding)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}
	if k.ExpirationTime != "" {
		key.ExpirationTime, err = time.Parse(time.RFC3339, k.ExpirationTime)
		if err != nil {
			return fmt.Errorf("parse expiration date: %w", err)
		}
	}
	key.Algorithm = k.Algorithm
	return nil
}

type authorizedKeyOut struct {
	ID             string `json:"id" yaml:"id"`
	PrivateKey     string `json:"private_key" yaml:"private_key"`
	PublicKey      string `json:"public_key" yaml:"public_key"`
	ExpirationTime string `json:"expiration_time,omitempty" yaml:"expiration_time,omitempty"`
	Algorithm      string `json:"algorithm" yaml:"algorithm"`
}

func newAuthorizedKeyOut(key AuthorizedKey) (out authorizedKeyOut, err error) {
	out.ID = key.ID
	out.PrivateKey, err = key.EncodePrivateKey(base64.StdEncoding)
	if err != nil {
		return authorizedKeyOut{}, fmt.Errorf("marshal private key: %w", err)
	}
	out.PublicKey, err = key.EncodePublicKey(base64.StdEncoding)
	if err != nil {
		return authorizedKeyOut{}, fmt.Errorf("marshal public key: %w", err)
	}
	if !key.ExpirationTime.IsZero() {
		out.ExpirationTime = key.ExpirationTime.Format(time.RFC3339)
	}
	out.Algorithm = key.Algorithm

	return out, nil
}
