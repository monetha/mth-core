package mthtoken

import corejwt "gitlab.com/monetha/mth-core/security/jwt"

//TokenBuilderFactory builds new TokenBuilders
type TokenBuilderFactory struct {
	Signer *corejwt.Signer
}

// NewTokenBuilderFactory returns new TokenBuilderFactory
func NewTokenBuilderFactory(signer *corejwt.Signer) *TokenBuilderFactory {
	return &TokenBuilderFactory{signer}
}

// New builds new fluent token builder
func (tokenBuilderFactory *TokenBuilderFactory) New() (builder *TokenBuilder) {
	return NewTokenBuilder(tokenBuilderFactory.Signer)
}
