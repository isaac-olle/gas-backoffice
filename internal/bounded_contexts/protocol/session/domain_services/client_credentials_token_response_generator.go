package domain_services

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
)

type ClientCredentialsTokenResponseGenerator struct {
	accessTokenGenerator session.AccessTokenGenerator
}

func NewClientCredentialsTokenResponseGenerator(
	accessTokenGenerator session.AccessTokenGenerator,
) *ClientCredentialsTokenResponseGenerator {
	return &ClientCredentialsTokenResponseGenerator{accessTokenGenerator: accessTokenGenerator}
}

func (g *ClientCredentialsTokenResponseGenerator) Generate(
	ctx context.Context,
	sess *session.ClientCredentialsSession,
) (*session.TokenResponse, error) {
	claims := sess.TokenClaims()
	accessToken, err := g.accessTokenGenerator.GenerateAccessToken(ctx, claims)
	if err != nil {
		return nil, err
	}

	return &session.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(session.AccessTokenTTL.Seconds()),
		Scope:       claims.Scope,
	}, nil
}
