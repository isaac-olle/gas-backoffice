package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure/helpers"
	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
)

// OidcIdentityProviderClient exchanges an authorization code for tokens with
// an OIDC provider following RFC 6749 / OpenID Connect Core.
type OidcIdentityProviderClient struct {
	verifier *helpers.OidcTokenVerifier
}

func NewOidcIdentityProviderClient(verifier *helpers.OidcTokenVerifier) *OidcIdentityProviderClient {
	return &OidcIdentityProviderClient{verifier: verifier}
}

func (c *OidcIdentityProviderClient) ExchangeCode(
	ctx context.Context,
	idp *protocolShared.IdentityProviderConfig,
	code string,
) (session.IdpAuthentication, error) {
	body := url.Values{}
	body.Set("code", code)
	body.Set("client_id", idp.ClientID)
	body.Set("client_secret", idp.ClientSecret)
	body.Set("redirect_uri", idp.BackofficeRedirectUri)
	body.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		idp.TokenEndpoint,
		strings.NewReader(body.Encode()),
	)
	if err != nil {
		return session.IdpAuthentication{}, fmt.Errorf(
			"build token request: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return session.IdpAuthentication{}, fmt.Errorf(
			"token exchange with oidc provider: %w",
			err,
		)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return session.IdpAuthentication{}, fmt.Errorf(
			"oidc token endpoint returned status %d",
			resp.StatusCode,
		)
	}

	var tokenResp struct {
		IDToken      string `json:"id_token"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return session.IdpAuthentication{}, fmt.Errorf(
			"decode oidc token response: %w",
			err,
		)
	}

	if tokenResp.IDToken == "" {
		return session.IdpAuthentication{}, fmt.Errorf(
			"oidc token response missing id_token",
		)
	}

	if tokenResp.AccessToken == "" {
		return session.IdpAuthentication{}, fmt.Errorf(
			"oidc token response missing access_token",
		)
	}

	verifiedClaims, err := c.verifier.Verify(ctx, idp, tokenResp.IDToken)
	if err != nil {
		return session.IdpAuthentication{}, fmt.Errorf("verify id_token: %w", err)
	}

	sub, _ := verifiedClaims["sub"].(string)
	email, _ := verifiedClaims["email"].(string)
	if sub == "" {
		return session.IdpAuthentication{}, fmt.Errorf("id_token missing sub claim")
	}
	if email == "" {
		return session.IdpAuthentication{}, fmt.Errorf("id_token missing email claim")
	}
	claims := sessionvos.IdpClaims{Sub: sub, Email: email}

	return session.IdpAuthentication{
		Tokens: session.IdpTokens{
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.RefreshToken,
			IDToken:      tokenResp.IDToken,
			ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		},
		Claims: claims,
	}, nil
}
