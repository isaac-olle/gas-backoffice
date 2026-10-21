package vo

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

type PrivateKey struct {
	value []byte
}

func NewPrivateKey(value []byte) (PrivateKey, error) {
	if len(value) == 0 {
		return PrivateKey{}, fmt.Errorf("private key cannot be empty")
	}
	return PrivateKey{value: value}, nil
}

func NewPrivateKeyFromHexString(value string) (PrivateKey, error) {
	if value == "" {
		return PrivateKey{}, fmt.Errorf("private key cannot be empty")
	}
	der, err := hex.DecodeString(value)
	if err != nil {
		return PrivateKey{}, fmt.Errorf("private key is not valid hex: %w", err)
	}
	if _, err = x509.ParsePKCS8PrivateKey(der); err != nil {
		return PrivateKey{}, fmt.Errorf("private key is not a valid PKCS8 private key: %w", err)
	}
	return PrivateKey{value: der}, nil
}

// NewPrivateKeyFromPEMString validates a PEM-encoded PKCS8 private key and
// stores the PEM text as-is (not the decoded DER), since consumers (e.g. the
// JWT signer) expect a PEM block.
func NewPrivateKeyFromPEMString(value string) (PrivateKey, error) {
	if value == "" {
		return PrivateKey{}, fmt.Errorf("private key cannot be empty")
	}
	block, _ := pem.Decode([]byte(value))
	if block == nil {
		return PrivateKey{}, fmt.Errorf("private key is not valid PEM")
	}
	if _, err := x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return PrivateKey{}, fmt.Errorf("private key is not a valid PKCS8 private key: %w", err)
	}
	return PrivateKey{value: []byte(value)}, nil
}

func (p PrivateKey) Bytes() []byte  { return p.value }
func (p PrivateKey) String() string { return hex.EncodeToString(p.value) }
