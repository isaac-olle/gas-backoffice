package key_factory

import (
	"fmt"

	sharedCrypto "github.com/isaac-olle/gas-backoffice/internal/shared/crypto"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/private_keys"
	"github.com/isaac-olle/gas-backoffice/internal/shared/crypto/public_keys"
)

// NewSigningKeyFromPEM builds the right SigningKey implementation (RSA or
// ECDSA) for the given algorithm.
func NewSigningKeyFromPEM(pemBytes []byte, algorithm string) (sharedCrypto.SigningKey, error) {
	switch {
	case private_keys.IsRSAAlgorithm(algorithm):
		return private_keys.NewRSAPrivateKeyFromPEM(pemBytes, algorithm)
	case private_keys.IsECDSAAlgorithm(algorithm):
		return private_keys.NewECDSAPrivateKeyFromPEM(pemBytes, algorithm)
	default:
		return nil, fmt.Errorf("unsupported signing algorithm: %s", algorithm)
	}
}

// NewVerifyingKeyFromPEM builds the right VerifyingKey implementation (RSA
// or ECDSA) for the given algorithm.
func NewVerifyingKeyFromPEM(pemBytes []byte, algorithm string) (sharedCrypto.VerifyingKey, error) {
	switch {
	case public_keys.IsRSAAlgorithm(algorithm):
		return public_keys.NewRSAPublicKeyFromPEM(pemBytes, algorithm)
	case public_keys.IsECDSAAlgorithm(algorithm):
		return public_keys.NewECDSAPublicKeyFromPEM(pemBytes, algorithm)
	default:
		return nil, fmt.Errorf("unsupported verifying algorithm: %s", algorithm)
	}
}
