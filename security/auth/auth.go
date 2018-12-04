package auth

import (
	"encoding/base64"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type principalClaims struct {
	jwt.StandardClaims
	Reach *string `json:"reach,omitempty"`
}

type secretString string

// GetBytes retrieves a base64 byte array
func (s secretString) GetBytes() ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(string(s))
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret: %v", err)
	}

	return b, nil
}
