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
) (*session.TokenResponse, error) {
	claims := sess.TokenClaims()
	if claims.Subject == "" {
		return nil, app_errors.InternalError(session.ErrAuthenticatedUserNotResolved).
			WithDescription("cannot generate tokens: authenticated user not resolved")
	}
	if claims.IdpID == "" {
		return nil, app_errors.InternalError(session.ErrInvalidSessionParameters).
			WithDescription("cannot generate tokens: identity provider not resolved")
	}
	accessToken, err := g.accessTokenGenerator.GenerateAccessToken(ctx, claims)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	response := &session.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(session.AccessTokenTTL.Seconds()),
		Scope:       claims.Scope,
	}

	return g.withIDToken(ctx, claims, response)
}

func (g *AuthorizationCodeTokenResponseGenerator) withIDToken(
	ctx context.Context,
	claims session.TokenClaims,
	response *session.TokenResponse,
) (*session.TokenResponse, error) {
	if !claims.HasScope("openid") {
		return response, nil
	}

	idToken, err := g.idTokenGenerator.GenerateIDToken(ctx, claims)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}
	response.IDToken = idToken

	return response, nil
}
