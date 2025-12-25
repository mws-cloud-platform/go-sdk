// Package iam provides types and utilities for IAM services.
package iam

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

// ES256 is the name of the ECDSA algorithm.
const ES256 = "ES256"

// ErrInvalidPrivateKeyFormat reports that the private key format is invalid.
const ErrInvalidPrivateKeyFormat = consterr.Error("invalid private key format")

// ParsePrivateKey decodes an ECDSA private key from a string.
func ParsePrivateKey(s string, base64Encoding *base64.Encoding) (*ecdsa.PrivateKey, error) {
	decoded, err := base64Encoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}
	parsed, err := x509.ParsePKCS8PrivateKey(decoded)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	key, ok := parsed.(*ecdsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidPrivateKeyFormat
	}

	return key, nil
}

// EncodePublicKey encodes an ECDSA public key to a string.
func EncodePublicKey(key *ecdsa.PrivateKey, encoding *base64.Encoding) (string, error) {
	publicKey, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		return "", fmt.Errorf("marshal public key: %w", err)
	}
	return encoding.EncodeToString(publicKey), nil
}
