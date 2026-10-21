package jwks

import (
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/key_factory"
)

// SigningKeyProjection is protocol's read model of a key the authorization
// server signs tokens with. Only the public half is projected: it is what
// gets published at /oidc/jwks so relying parties can verify our tokens.
// Today it is read straight from security.key_pairs; in the future it will
// be a projection fed by security's events.
type SigningKeyProjection struct {
	Kid          string
	Algorithm    string
	PublicKeyPEM []byte
}

func NewSigningKeyProjection(kid, algorithm string, publicKeyPEM []byte) *SigningKeyProjection {
	return &SigningKeyProjection{Kid: kid, Algorithm: algorithm, PublicKeyPEM: publicKeyPEM}
}

// ToJWK renders the public key as a JWK (RFC 7517), tagged with its kid and
// intended use so verifiers can match it against a token's kid header.
func (p *SigningKeyProjection) ToJWK() (map[string]any, error) {
	key, err := key_factory.NewVerifyingKeyFromPEM(p.PublicKeyPEM, p.Algorithm)
	if err != nil {
		return nil, fmt.Errorf("parse public key %s: %w", p.Kid, err)
	}
	jwk, err := key.ToJWK()
	if err != nil {
		return nil, fmt.Errorf("render jwk %s: %w", p.Kid, err)
	}
	jwk["kid"] = p.Kid
	jwk["use"] = "sig"
	return jwk, nil
}
