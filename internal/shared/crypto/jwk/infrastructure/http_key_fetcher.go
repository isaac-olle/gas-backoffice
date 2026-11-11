package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sharedCrypto "github.com/isaac-olle/gas-backoffice/internal/shared/crypto"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/public_keys"
)

// HttpKeyFetcher is the KeyFetcher that calls a JWKS HTTP endpoint and builds
// a VerifyingKey out of each JWK it returns.
type HttpKeyFetcher struct{}

func NewHttpKeyFetcher() *HttpKeyFetcher {
	return &HttpKeyFetcher{}
}

func (f *HttpKeyFetcher) FetchByKid(ctx context.Context, jwksUri, kid string) (sharedCrypto.VerifyingKey, error) {
	keys, err := f.FetchAll(ctx, jwksUri)
	if err != nil {
		return nil, err
	}
	key, ok := keys[kid]
	if !ok {
		return nil, fmt.Errorf("kid %q not found in jwks %s", kid, jwksUri)
	}
	return key, nil
}

func (f *HttpKeyFetcher) FetchAll(ctx context.Context, jwksUri string) (map[string]sharedCrypto.VerifyingKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksUri, nil)
	if err != nil {
		return nil, fmt.Errorf("build jwks request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch jwks: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read jwks response: %w", err)
	}

	var parsed struct {
		Keys []map[string]any `json:"keys"`
	}
	if err = json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decode jwks response: %w", err)
	}

	keys := make(map[string]sharedCrypto.VerifyingKey, len(parsed.Keys))
	for _, jwk := range parsed.Keys {
		kid, _ := jwk["kid"].(string)
		if kid == "" {
			continue
		}

		key, err := public_keys.NewVerifyingKeyFromJWK(jwk)
		if err != nil {
			continue
		}
		keys[kid] = key
	}
	return keys, nil
}
