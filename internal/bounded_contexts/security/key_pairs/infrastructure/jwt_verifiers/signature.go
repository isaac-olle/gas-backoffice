package jwt_verifiers

import (
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/key_factory"
)

// Signature checks only that a JWT's signature is valid for a given
// KeyPair's public key, returning the raw claims. It doesn't know which
// claims should be present or what they should contain.
type Signature struct{}

func NewSignatureVerifier() *Signature {
	return &Signature{}
}

func (v *Signature) Verify(token string, kp *key_pairs.KeyPair) (map[string]any, error) {
	p := kp.ToPrimitives()

	verifyingKey, err := key_factory.NewVerifyingKeyFromPEM([]byte(p.PublicKey), p.Algorithm)
	if err != nil {
		return nil, fmt.Errorf("invalid public key in key pair: %w", err)
	}

	claims, err := verifyingKey.VerifyString(token)
	if err != nil {
		return nil, fmt.Errorf("token signature verification failed: %w", err)
	}

	return claims, nil
}
