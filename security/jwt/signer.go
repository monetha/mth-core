package jwt

import (
	"encoding/base64"
	"fmt"
)

// Signer is a key to sign JWT tokens
type Signer struct {
	SecretString string
	Bytes        []byte
}

// GetBytes retrieves a base64 byte array
func getBytes(secret string) ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret: %v", err)
	}

	return b, nil
}

// NewSigner creates a new service auth object.
func NewSigner(secret string) (*Signer, error) {
	b, err := getBytes(secret)
	if err != nil {
		return nil, err
	}
	return &Signer{secret, b}, nil
}
