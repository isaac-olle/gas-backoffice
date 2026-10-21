package symmetric_keys

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	sharedCrypto "github.com/isaac-olle/gas-backoffice/internal/shared/crypto"
	sharedJwk "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwk"
	sharedJwt "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwt"
)

var hmacHashes = map[string]crypto.Hash{
	"HS256": crypto.SHA256,
	"HS384": crypto.SHA384,
	"HS512": crypto.SHA512,
}

func IsHMACAlgorithm(algorithm string) bool {
	_, ok := hmacHashes[algorithm]
	return ok
}

// hmacKey implements both SigningKey and VerifyingKey: the same secret
// signs and verifies, so unlike RSA/ECDSA there is no public/private split.
type hmacKey struct {
	secret    []byte
	algorithm string
	hash      crypto.Hash
}

var (
	_ sharedCrypto.SigningKey   = (*hmacKey)(nil)
	_ sharedCrypto.VerifyingKey = (*hmacKey)(nil)
)

func NewHMACKey(secret []byte, algorithm string) (*hmacKey, error) {
	hash, ok := hmacHashes[algorithm]
	if !ok {
		return nil, fmt.Errorf("unsupported hmac algorithm: %s", algorithm)
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("secret cannot be empty")
	}
	return &hmacKey{secret: secret, algorithm: algorithm, hash: hash}, nil
}

// ToPEMFormat is not a real PEM/ASN.1 structure — an HMAC secret is just raw
// bytes, not a standard key format — but wraps it as PEM for interface
// consistency with the asymmetric keys.
func (k *hmacKey) ToPEMFormat() ([]byte, error) {
	return pem.EncodeToMemory(&pem.Block{Type: "HMAC KEY", Bytes: k.secret}), nil
}

// NewHMACKeyFromJWK builds a key from a JWK (RFC 7517) of kty "oct".
func NewHMACKeyFromJWK(jwk map[string]any) (*hmacKey, error) {
	secret, err := sharedJwk.Bytes(jwk, "k")
	if err != nil {
		return nil, err
	}
	algorithm, err := sharedJwk.String(jwk, "alg")
	if err != nil {
		return nil, err
	}
	return NewHMACKey(secret, algorithm)
}

func (k *hmacKey) ToJWK() (map[string]any, error) {
	return map[string]any{
		"kty": "oct",
		"alg": k.algorithm,
		"k":   base64.RawURLEncoding.EncodeToString(k.secret),
	}, nil
}

func (k *hmacKey) SignBytes(data []byte) ([]byte, error) {
	mac := hmac.New(k.hash.New, k.secret)
	mac.Write(data)
	return mac.Sum(nil), nil
}

func (k *hmacKey) Algorithm() string {
	return k.algorithm
}

func (k *hmacKey) SignString(signingInput string) (string, error) {
	return sharedJwt.SignCompactToken(k.SignBytes, signingInput)
}

func (k *hmacKey) VerifyBytes(data, signature []byte) error {
	expected, err := k.SignBytes(data)
	if err != nil {
		return err
	}
	if !hmac.Equal(expected, signature) {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}

func (k *hmacKey) VerifyString(tokenString string) (map[string]any, error) {
	return sharedJwt.VerifyToken(k.VerifyBytes, tokenString)
}
