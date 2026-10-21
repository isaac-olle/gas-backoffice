package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
)

type AuthorizationCodeTokenResponseGenerator struct {
	accessTokenGenerator session.AccessTokenGenerator
	idTokenGenerator     session.IDTokenGenerator
}

func NewAuthorizationCodeTokenResponseGenerator(
	accessTokenGenerator session.AccessTokenGenerator,
	idTokenGenerator session.IDTokenGenerator,
) *AuthorizationCodeTokenResponseGenerator {
	return &AuthorizationCodeTokenResponseGenerator{
		accessTokenGenerator: accessTokenGenerator,
		idTokenGenerator:     idTokenGenerator,
	}
}

func (g *AuthorizationCodeTokenResponseGenerator) Generate(
	ctx context.Context,
	sess *session.AuthorizationCodeSession,
	identity *session.UserIdentityProjection,
) (*session.TokenResponse, error) {
	accessToken, err := g.accessTokenGenerator.GenerateAccessToken(ctx, sess.BaseSession, identity)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	response := &session.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(session.AccessTokenTTL.Seconds()),
		Scope:       sess.TokenClaims().Scope,
	}

	return g.withIDToken(ctx, sess, identity, response)
}

func (g *AuthorizationCodeTokenResponseGenerator) withIDToken(
	ctx context.Context,
	sess *session.AuthorizationCodeSession,
	identity *session.UserIdentityProjection,
	response *session.TokenResponse,
) (*session.TokenResponse, error) {
	if !sess.HasScope("openid") {
		return response, nil
	}

	idToken, err := g.idTokenGenerator.GenerateIDToken(ctx, sess.BaseSession, identity)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}
	response.IDToken = idToken

	return response, nil
}
