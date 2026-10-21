package public_keys

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"

	sharedCrypto "github.com/isaac-olle/gas-backoffice/internal/shared/crypto"
	sharedJwk "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwk"
	sharedJwt "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwt"
)

var rsaHashes = map[string]crypto.Hash{
	"RS256": crypto.SHA256,
	"RS384": crypto.SHA384,
	"RS512": crypto.SHA512,
}

func IsRSAAlgorithm(algorithm string) bool {
	_, ok := rsaHashes[algorithm]
	return ok
}

type rsaPublicKey struct {
	publicKey *rsa.PublicKey
	algorithm string
	hash      crypto.Hash
}

var _ sharedCrypto.VerifyingKey = (*rsaPublicKey)(nil)

func NewRSAPublicKeyFromPEM(pemBytes []byte, algorithm string) (*rsaPublicKey, error) {
	hash, ok := rsaHashes[algorithm]
	if !ok {
		return nil, fmt.Errorf("unsupported rsa algorithm: %s", algorithm)
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode public key PEM")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	rsaKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an rsa public key")
	}
	return &rsaPublicKey{publicKey: rsaKey, algorithm: algorithm, hash: hash}, nil
}

// NewRSAPublicKeyFromJWK builds a verifying key from a JWK (RFC 7517), e.g.
// one entry of a JWKS response, without a PEM round-trip.
func NewRSAPublicKeyFromJWK(jwk map[string]any) (*rsaPublicKey, error) {
	modulus, err := sharedJwk.Bytes(jwk, "n")
	if err != nil {
		return nil, err
	}
	exponent, err := sharedJwk.Bytes(jwk, "e")
	if err != nil {
		return nil, err
	}
	algorithm, err := sharedJwk.String(jwk, "alg")
	if err != nil {
		algorithm = "RS256"
	}
	hash, ok := rsaHashes[algorithm]
	if !ok {
		return nil, fmt.Errorf("unsupported rsa algorithm: %s", algorithm)
	}
	publicKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulus),
		E: int(new(big.Int).SetBytes(exponent).Int64()),
	}
	return &rsaPublicKey{publicKey: publicKey, algorithm: algorithm, hash: hash}, nil
}

func (k *rsaPublicKey) ToPEMFormat() ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(k.publicKey)
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), nil
}

func (k *rsaPublicKey) ToJWK() (map[string]any, error) {
	return map[string]any{
		"kty": "RSA",
		"alg": k.algorithm,
		"n":   base64.RawURLEncoding.EncodeToString(k.publicKey.N.Bytes()),
		"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(k.publicKey.E)).Bytes()),
	}, nil
}

func (k *rsaPublicKey) Algorithm() string {
	return k.algorithm
}

func (k *rsaPublicKey) VerifyBytes(data, signature []byte) error {
	hasher := k.hash.New()
	hasher.Write(data)
	return rsa.VerifyPKCS1v15(k.publicKey, k.hash, hasher.Sum(nil), signature)
}

func (k *rsaPublicKey) VerifyString(tokenString string) (map[string]any, error) {
	return sharedJwt.VerifyToken(k.VerifyBytes, tokenString)
}
