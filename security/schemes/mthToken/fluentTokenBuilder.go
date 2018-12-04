package mthtoken

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	coreJwt "gitlab.com/monetha/mth-core/security/jwt"
)

// TokenBuilder is fluent syntax for building mth-token token
type TokenBuilder struct {
	Signer *coreJwt.Signer
	Claims *principalClaims
}

// NewTokenBuilder creates a new fluent token builder
func NewTokenBuilder(signer *coreJwt.Signer) *TokenBuilder {
	return &TokenBuilder{signer, defaultClaims()}
}

// Build builds JWT token and signs
func (tokenBuilder *TokenBuilder) Build() (tokenString string, err error) {
	token := jwt.NewWithClaims(serviceAuthSigningMethod, tokenBuilder.Claims)
	if tokenString, err = token.SignedString(tokenBuilder.Signer.Bytes); err != nil {
		err = fmt.Errorf("failed to return signed token string: %v", err)
	}
	return
}

// WithSystemReach creates system claims and forms a signed JWT.
func (tokenBuilder *TokenBuilder) WithSystemReach() (builder *TokenBuilder) {
	tokenBuilder.Claims.Reach = &systemAuth
	return tokenBuilder
}
