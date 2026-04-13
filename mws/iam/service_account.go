package iam

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

var serviceAccountAuthorizedKeyIDRe = regexp.MustCompile(`^projects/(.*)/serviceAccounts/(.*)/authorizedKeys/(.*)$`)

// InvalidServiceAccountAuthorizedKeyIDError is returned when the service
// account authorized key id is invalid.
type InvalidServiceAccountAuthorizedKeyIDError struct {
	ID string
}

func (e InvalidServiceAccountAuthorizedKeyIDError) Error() string {
	return fmt.Sprintf("invalid service account authorized key id: %s", e.ID)
}

// UnsupportedAlgorithmError is returned when the crypto algorithm is not
// supported.
type UnsupportedAlgorithmError struct {
	Algorithm string
}

func (e UnsupportedAlgorithmError) Error() string {
	return fmt.Sprintf("unsupported algorithm: %s", e.Algorithm)
}

type serviceAccountCtxKey struct{}

// WithServiceAccount adds given service account data to the context.
func WithServiceAccount(ctx context.Context, sa ServiceAccount) context.Context {
	return context.WithValue(ctx, serviceAccountCtxKey{}, sa)
}

// ServiceAccountFromContext retrieves service account data from the context.
func ServiceAccountFromContext(ctx context.Context) (ServiceAccount, bool) {
	if ctx == nil {
		return ServiceAccount{}, false
	}
	v, ok := ctx.Value(serviceAccountCtxKey{}).(ServiceAccount)
	if !ok {
		return ServiceAccount{}, false
	}
	return v, true
}

// ServiceAccount contains service account information.
type ServiceAccount struct {
	Project string
	Name    string
}

func (s ServiceAccount) String() string {
	return fmt.Sprintf("projects/%s/serviceAccounts/%s", s.Project, s.Name)
}

func (s ServiceAccount) impersonable() {}

type serviceAccountAuthorizedKeyCtxKey struct{}

// WithServiceAccountAuthorizedKey adds given service account authorized key to
// the context.
func WithServiceAccountAuthorizedKey(ctx context.Context, key ServiceAccountAuthorizedKey) context.Context {
	return context.WithValue(ctx, serviceAccountAuthorizedKeyCtxKey{}, key)
}

// ServiceAccountAuthorizedKeyFromContext retrieves service account authorized
// key from the context.
func ServiceAccountAuthorizedKeyFromContext(ctx context.Context) (ServiceAccountAuthorizedKey, bool) {
	if ctx == nil {
		return ServiceAccountAuthorizedKey{}, false
	}
	v, ok := ctx.Value(serviceAccountAuthorizedKeyCtxKey{}).(ServiceAccountAuthorizedKey)
	if !ok {
		return ServiceAccountAuthorizedKey{}, false
	}
	return v, true
}

// ServiceAccountAuthorizedKey contains a service account and it's authorized key.
type ServiceAccountAuthorizedKey struct {
	ServiceAccount ServiceAccount
	AuthorizedKey  AuthorizedKey
}

// Reference returns a reference to the service account authorized key.
func (k ServiceAccountAuthorizedKey) Reference() string {
	return fmt.Sprintf("projects/%s/serviceAccounts/%s/authorizedKeys/%s",
		k.ServiceAccount.Project, k.ServiceAccount.Name, k.AuthorizedKey.Name)
}

func (k *ServiceAccountAuthorizedKey) UnmarshalJSON(data []byte) error {
	v := struct {
		ID             string `json:"keyId"`
		PrivateKey     string `json:"privateKey"`
		PublicKey      string `json:"publicKey"`
		ExpirationTime string `json:"expirationTime"`
		Algorithm      string `json:"algorithm"`
	}{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.Algorithm != ES256 {
		return UnsupportedAlgorithmError{Algorithm: v.Algorithm}
	}
	k.AuthorizedKey.Algorithm = v.Algorithm

	matches := serviceAccountAuthorizedKeyIDRe.FindStringSubmatch(v.ID)
	if len(matches) != 4 { //nolint:mnd // regex has 4 groups
		return InvalidServiceAccountAuthorizedKeyIDError{ID: v.ID}
	}

	k.ServiceAccount = ServiceAccount{
		Project: matches[1],
		Name:    matches[2],
	}
	k.AuthorizedKey.Name = matches[3]

	decoded, err := base64.StdEncoding.DecodeString(v.PrivateKey)
	if err != nil {
		return fmt.Errorf("decode private key: %w", err)
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(decoded)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}

	var ok bool
	k.AuthorizedKey.PrivateKey, ok = privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return ErrInvalidPrivateKeyFormat
	}

	if v.ExpirationTime != "" {
		k.AuthorizedKey.ExpirationTime, err = time.Parse(time.RFC3339, v.ExpirationTime)
		if err != nil {
			return fmt.Errorf("parse expiration date: %w", err)
		}
	}

	return nil
}
