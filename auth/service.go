package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	defaultServiceAuthExpiration = 1 * time.Minute
)

var (
	serviceAuthSigningMethod = jwt.SigningMethodHS256
	systemAuth               = "system"
)

// ServiceAuth is for making safe calls to other Monetha API instances.
type ServiceAuth struct {
	secret []byte
}

// NewServiceAuth creates a new service auth object.
func NewServiceAuth(secret string) (*ServiceAuth, error) {
	b, err := secretString(secret).GetBytes()
	if err != nil {
		return nil, err
	}
	return &ServiceAuth{b}, nil
}

// NewSystemAuth creates system claims and forms a signed JWT.
func (sauth *ServiceAuth) NewSystemAuth() (s string, err error) {
	now := time.Now().UTC()
	claims := principalClaims{}
	claims.IssuedAt = now.Unix()
	claims.NotBefore = now.Unix()
	claims.ExpiresAt = now.Add(defaultServiceAuthExpiration).Unix()
	claims.Reach = &systemAuth

	token := jwt.NewWithClaims(serviceAuthSigningMethod, claims)
	if s, err = token.SignedString(sauth.secret); err != nil {
		err = fmt.Errorf("failed to return signed token string: %v", err)
	}
	return
}
