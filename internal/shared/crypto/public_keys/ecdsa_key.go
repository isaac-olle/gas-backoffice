package public_keys

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"

	sharedCrypto "github.com/isaac-olle/gas-backoffice/internal/shared/crypto"
	sharedJwk "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwk"
	sharedJwt "github.com/isaac-olle/gas-backoffice/internal/shared/crypto/jwt"
)

type ecdsaAlgConfig struct {
	curve   elliptic.Curve
	hash    crypto.Hash
	size    int
	crvName string
}

var ecdsaCurves = map[string]ecdsaAlgConfig{
	"ES256": {elliptic.P256(), crypto.SHA256, 32, "P-256"},
	"ES384": {elliptic.P384(), crypto.SHA384, 48, "P-384"},
	"ES512": {elliptic.P521(), crypto.SHA512, 66, "P-521"},
}

func IsECDSAAlgorithm(algorithm string) bool {
	_, ok := ecdsaCurves[algorithm]
	return ok
}

type ecdsaPublicKey struct {
	publicKey *ecdsa.PublicKey
	algorithm string
	hash      crypto.Hash
	size      int
	crvName   string
}

var _ sharedCrypto.VerifyingKey = (*ecdsaPublicKey)(nil)

func NewECDSAPublicKeyFromPEM(pemBytes []byte, algorithm string) (*ecdsaPublicKey, error) {
	cfg, ok := ecdsaCurves[algorithm]
	if !ok {
		return nil, fmt.Errorf("unsupported ecdsa algorithm: %s", algorithm)
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode public key PEM")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	ecKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ecdsa public key")
	}
	return &ecdsaPublicKey{publicKey: ecKey, algorithm: algorithm, hash: cfg.hash, size: cfg.size, crvName: cfg.crvName}, nil
}

// NewECDSAPublicKeyFromJWK builds a verifying key from a JWK (RFC 7517),
// e.g. one entry of a JWKS response, without a PEM round-trip.
func NewECDSAPublicKeyFromJWK(jwk map[string]any) (*ecdsaPublicKey, error) {
	x, err := sharedJwk.Bytes(jwk, "x")
	if err != nil {
		return nil, err
	}
	y, err := sharedJwk.Bytes(jwk, "y")
	if err != nil {
		return nil, err
	}
	crv, err := sharedJwk.String(jwk, "crv")
	if err != nil {
		return nil, err
	}

	var algorithm string
	var cfg ecdsaAlgConfig
	found := false
	for alg, c := range ecdsaCurves {
		if c.crvName == crv {
			algorithm, cfg, found = alg, c, true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("unsupported ecdsa curve: %s", crv)
	}

	return &ecdsaPublicKey{
		publicKey: &ecdsa.PublicKey{Curve: cfg.curve, X: new(big.Int).SetBytes(x), Y: new(big.Int).SetBytes(y)},
		algorithm: algorithm,
		hash:      cfg.hash,
		size:      cfg.size,
		crvName:   cfg.crvName,
	}, nil
}

func (k *ecdsaPublicKey) ToPEMFormat() ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(k.publicKey)
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), nil
}

func (k *ecdsaPublicKey) ToJWK() (map[string]any, error) {
	encode := func(b *big.Int) string {
		out := make([]byte, k.size)
		b.FillBytes(out)
		return base64.RawURLEncoding.EncodeToString(out)
	}

	return map[string]any{
		"kty": "EC",
		"alg": k.algorithm,
		"crv": k.crvName,
		"x":   encode(k.publicKey.X),
		"y":   encode(k.publicKey.Y),
	}, nil
}

func (k *ecdsaPublicKey) Algorithm() string {
	return k.algorithm
}

func (k *ecdsaPublicKey) VerifyBytes(data, signature []byte) error {
	if len(signature) != 2*k.size {
		return fmt.Errorf("invalid signature length")
	}
	r := new(big.Int).SetBytes(signature[:k.size])
	s := new(big.Int).SetBytes(signature[k.size:])
	hasher := k.hash.New()
	hasher.Write(data)
	if !ecdsa.Verify(k.publicKey, hasher.Sum(nil), r, s) {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}

func (k *ecdsaPublicKey) VerifyString(tokenString string) (map[string]any, error) {
	return sharedJwt.VerifyToken(k.VerifyBytes, tokenString)
}
