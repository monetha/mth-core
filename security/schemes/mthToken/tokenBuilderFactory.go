package mthtoken

import coreJwt "gitlab.com/monetha/mth-core/security/jwt"

//TokenBuilderFactory builds new TokenBuilders
type TokenBuilderFactory struct {
	Signer *coreJwt.Signer
}

// NewTokenBuilderFactory returns new TokenBuilderFactory
func NewTokenBuilderFactory(signer *coreJwt.Signer) *TokenBuilderFactory {
	return &TokenBuilderFactory{signer}
}

// New builds new fluent token builder
func (tokenBuilderFactory *TokenBuilderFactory) New() (builder *TokenBuilder) {
	return NewTokenBuilder(tokenBuilderFactory.Signer)
}
