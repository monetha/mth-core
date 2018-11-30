package mthtoken

import (
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

type principalClaims struct {
	jwt.StandardClaims
	Reach *string `json:"reach,omitempty"`
}

func defaultClaims() *principalClaims {
	now := time.Now().UTC()
	claims := principalClaims{}
	claims.IssuedAt = now.Unix()
	claims.NotBefore = now.Unix()
	claims.ExpiresAt = now.Add(defaultServiceAuthExpiration).Unix()
	return &claims
}
