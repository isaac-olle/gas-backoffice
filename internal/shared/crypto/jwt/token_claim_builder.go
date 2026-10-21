package jwt

import "time"

// TokenClaimBuilder builds the payload claims for a JWT. iss is fixed at
// construction because it's resolved by whoever holds the key (not the
// caller); sub/aud/jti/extraClaims come from the caller's request as-is,
// since TokenClaimBuilder has no business interpreting them.
type TokenClaimBuilder struct {
	iss         string
	sub         string
	aud         string
	jti         string
	extraClaims map[string]any
}

func NewTokenClaimBuilder(iss, sub, aud, jti string, extraClaims map[string]any) *TokenClaimBuilder {
	return &TokenClaimBuilder{iss: iss, sub: sub, aud: aud, jti: jti, extraClaims: extraClaims}
}

// Build returns the claims map, including iat/exp computed from expiresIn.
func (b *TokenClaimBuilder) Build(expiresIn time.Duration) map[string]any {
	claims := make(map[string]any, len(b.extraClaims)+4)
	for k, v := range b.extraClaims {
		claims[k] = v
	}
	claims["iss"] = b.iss
	claims["sub"] = b.sub
	if b.aud != "" {
		claims["aud"] = b.aud
	}
	if b.jti != "" {
		claims["jti"] = b.jti
	}

	now := time.Now()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(expiresIn).Unix()

	return claims
}
