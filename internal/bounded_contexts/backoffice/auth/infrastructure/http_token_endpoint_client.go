package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/request_id"
)

type HttpTokenEndpointClient struct {
	endpoint     string
	clientId     string
	clientSecret string
	redirectUri  string
	verifier     *GasIdTokenVerifier
}

func NewHttpTokenEndpointClient(endpoint, clientId, clientSecret, redirectUri string, verifier *GasIdTokenVerifier) *HttpTokenEndpointClient {
	return &HttpTokenEndpointClient{
		endpoint:     endpoint,
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectUri:  redirectUri,
		verifier:     verifier,
	}
}

func (c *HttpTokenEndpointClient) ExchangeCode(ctx context.Context, code, codeVerifier, idpReferenceName string) (auth.TokenExchangeResult, error) {
	if strings.Contains(c.redirectUri, ":idp_reference_name") && idpReferenceName == "" {
		return auth.TokenExchangeResult{}, fmt.Errorf("missing idp_reference_name in authorization state; restart login")
	}
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {strings.ReplaceAll(c.redirectUri, ":idp_reference_name", idpReferenceName)},
		"code_verifier": {codeVerifier},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return auth.TokenExchangeResult{}, fmt.Errorf("build token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientId, c.clientSecret)
	if reqID := request_id.FromContext(ctx); reqID != "" {
		req.Header.Set(request_id.Header, reqID)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return auth.TokenExchangeResult{}, fmt.Errorf("call token endpoint: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return auth.TokenExchangeResult{}, fmt.Errorf(
				"token endpoint returned %d: failed to read error response: %w",
				resp.StatusCode,
				err,
			)
		}

		return auth.TokenExchangeResult{}, fmt.Errorf(
			"token endpoint returned %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var body struct {
		IDToken string `json:"id_token"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return auth.TokenExchangeResult{}, fmt.Errorf("decode token response: %w", err)
	}

	claims, err := c.verifier.Verify(ctx, body.IDToken)
	if err != nil {
		return auth.TokenExchangeResult{}, fmt.Errorf("verify id_token: %w", err)
	}

	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	picture, _ := claims["picture"].(string)
	emailVerified, _ := claims["email_verified"].(bool)

	return auth.TokenExchangeResult{
		Sub:           sub,
		Email:         email,
		Picture:       picture,
		EmailVerified: emailVerified,
	}, nil
}
