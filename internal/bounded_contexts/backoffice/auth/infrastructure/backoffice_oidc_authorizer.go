package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/infrastructure/helpers/oidc"
	"io"
	"net/http"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
)

const (
	authorizeTimeout  = 10 * time.Second
	maxErrorBodyBytes = 4096
)

var _ auth.IdpConnector = (auth.IdpConnector)(nil)

// OidcAuthorizer parla amb el nostre propi authorization server per HTTP, encara
// que avui visqui al mateix procés. El backoffice és un client OIDC més del gas
// i el tracta com un servei extern: el dia que /oidc/* se'n vagi a gas-protocol,
// només canvia GAS_PROTOCOL_URL.
//
// /oidc/authorize respon 302 amb el Location cap a l'IdP federat. El client no
// segueix redireccions a propòsit — el que volem és aquesta capçalera, que
// retornem perquè el LoginHandler hi enviï el navegador de l'usuari.
type OidcAuthorizer struct {
	client            *http.Client
	requestBuilder    *oidc.RequestBuilder
	authorizeEndpoint string
	clientId          string
	redirectUri       string
}

func NewOidcAuthorizer(
	requestBuilder *oidc.RequestBuilder,
	authorizeEndpoint string,
	clientId string,
	redirectUri string,
) *OidcAuthorizer {
	return &OidcAuthorizer{
		client: &http.Client{
			Timeout: authorizeTimeout,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		requestBuilder:    requestBuilder,
		authorizeEndpoint: authorizeEndpoint,
		clientId:          clientId,
		redirectUri:       redirectUri,
	}
}

func (a *OidcAuthorizer) Connect(
	ctx context.Context,
	idp *auth.IdpConfig,
) (string, error) {
	requestURL, err := a.requestBuilder.Build(a.redirectUri, a.clientId, idp.IdpReferenceName, a.authorizeEndpoint)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("oidc: build authorize request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("oidc: call authorize endpoint: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusMultipleChoices || resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
		return "", fmt.Errorf("oidc: authorize endpoint returned %d: %s", resp.StatusCode, body)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("oidc: authorize endpoint returned %d without a Location header", resp.StatusCode)
	}

	return location, nil
}
