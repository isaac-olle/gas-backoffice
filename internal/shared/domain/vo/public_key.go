package vo

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

type PublicKey struct {
	value []byte
}

func NewPublicKey(value []byte) (PublicKey, error) {
	if len(value) == 0 {
		return PublicKey{}, fmt.Errorf("public key cannot be empty")
	}
	return PublicKey{value: value}, nil
}

func NewPublicKeyFromHexString(value string) (PublicKey, error) {
	if value == "" {
		return PublicKey{}, fmt.Errorf("public key cannot be empty")
	}
	der, err := hex.DecodeString(value)
	if err != nil {
		return PublicKey{}, fmt.Errorf("public key is not valid hex: %w", err)
	}
	if _, err = x509.ParsePKIXPublicKey(der); err != nil {
		return PublicKey{}, fmt.Errorf("public key is not a valid PKIX public key: %w", err)
	}
	return PublicKey{value: der}, nil
}

// NewPublicKeyFromPEMString validates a PEM-encoded PKIX public key and
// stores the PEM text as-is (not the decoded DER), since consumers (e.g. the
// JWT verifier) expect a PEM block.
func NewPublicKeyFromPEMString(value string) (PublicKey, error) {
	if value == "" {
		return PublicKey{}, fmt.Errorf("public key cannot be empty")
	}
	block, _ := pem.Decode([]byte(value))
	if block == nil {
		return PublicKey{}, fmt.Errorf("public key is not valid PEM")
	}
	if _, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return PublicKey{}, fmt.Errorf("public key is not a valid PKIX public key: %w", err)
	}
	return PublicKey{value: []byte(value)}, nil
}

func (p PublicKey) Bytes() []byte  { return p.value }
func (p PublicKey) String() string { return hex.EncodeToString(p.value) }
