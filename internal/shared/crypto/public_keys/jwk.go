package public_keys

import (
	"fmt"

	sharedCrypto "github.com/isaac-olle/gas-backoffice/internal/shared/crypto"
)

// NewVerifyingKeyFromJWK builds the VerifyingKey described by a JWK.
func NewVerifyingKeyFromJWK(jwk map[string]any) (sharedCrypto.VerifyingKey, error) {
	kty, _ := jwk["kty"].(string)
	switch kty {
	case "RSA":
		return NewRSAPublicKeyFromJWK(jwk)
	case "EC":
		return NewECDSAPublicKeyFromJWK(jwk)
	default:
		return nil, fmt.Errorf("unsupported jwk kty: %s", kty)
	}
}
