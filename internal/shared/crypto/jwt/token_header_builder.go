package jwt

import (
	"fmt"

	sharedCrypto "github.com/isaac-olle/gas-backoffice/internal/shared/crypto"
)

// TokenHeaderBuilder builds the JOSE header for a key. alg and kid always
// describe the key itself, so they're fixed at construction; only the
// caller's intent to embed key material (JWK) is a per-call input, and even
// that is built from the key TokenHeaderBuilder already holds — the caller
// never supplies key material directly.
type TokenHeaderBuilder struct {
	key sharedCrypto.Key
	kid string
}

func NewTokenHeaderBuilder(key sharedCrypto.Key, kid string) *TokenHeaderBuilder {
	return &TokenHeaderBuilder{key: key, kid: kid}
}

func (b *TokenHeaderBuilder) Build(includeJWK, includeX5c bool) (map[string]any, error) {
	if includeX5c {
		return nil, fmt.Errorf("x5c header is not supported yet")
	}

	header := map[string]any{
		"alg": b.key.Algorithm(),
		"typ": "JWT",
		"kid": b.kid,
	}

	if includeJWK {
		jwk, err := b.key.ToJWK()
		if err != nil {
			return nil, fmt.Errorf("build jwk header: %w", err)
		}
		header["jwk"] = jwk
	}

	return header, nil
}
