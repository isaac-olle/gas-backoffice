package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	sharedCrypto "github.com/isaac-olle/gas-backoffice/internal/shared/crypto"
)

// TokenSigner signs an already-built header and claims map. It knows
// nothing about JWT semantics (claims, kid, iat/exp) — that's
// TokenHeaderBuilder's and TokenClaimBuilder's job; TokenSigner only
// encodes and signs bytes.
type TokenSigner struct {
	key sharedCrypto.SigningKey
}

func NewTokenSigner(key sharedCrypto.SigningKey) *TokenSigner {
	return &TokenSigner{key: key}
}

func (s *TokenSigner) Sign(header, claims map[string]any) (string, error) {
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("encode header: %w", err)
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("encode claims: %w", err)
	}

	signingInput := base64.RawURLEncoding.EncodeToString(headerJSON) + "." + base64.RawURLEncoding.EncodeToString(claimsJSON)
	return s.key.SignString(signingInput)
}
