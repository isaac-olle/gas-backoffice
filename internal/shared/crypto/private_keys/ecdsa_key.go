package private_keys

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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

type ecdsaPrivateKey struct {
	privateKey *ecdsa.PrivateKey
	pemBytes   []byte
	algorithm  string
	hash       crypto.Hash
	size       int
	crvName    string
}

var _ sharedCrypto.SigningKey = (*ecdsaPrivateKey)(nil)

func NewECDSAPrivateKeyFromPEM(pemBytes []byte, algorithm string) (*ecdsaPrivateKey, error) {
	cfg, ok := ecdsaCurves[algorithm]
	if !ok {
		return nil, fmt.Errorf("unsupported ecdsa algorithm: %s", algorithm)
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key PEM")
	}
	privateKey, err := parseECDSAPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &ecdsaPrivateKey{privateKey: privateKey, pemBytes: pemBytes, algorithm: algorithm, hash: cfg.hash, size: cfg.size, crvName: cfg.crvName}, nil
}

func (k *ecdsaPrivateKey) ToPEMFormat() ([]byte, error) {
	if k.pemBytes != nil {
		return k.pemBytes, nil
	}
	der, err := x509.MarshalPKCS8PrivateKey(k.privateKey)
	if err != nil {
		return nil, fmt.Errorf("marshal private key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), nil
}

// NewECDSAPrivateKeyFromJWK builds a signing key from a JWK (RFC 7517/7518)
// that includes the private scalar d.
func NewECDSAPrivateKeyFromJWK(jwk map[string]any) (*ecdsaPrivateKey, error) {
	x, err := sharedJwk.Bytes(jwk, "x")
	if err != nil {
		return nil, err
	}
	y, err := sharedJwk.Bytes(jwk, "y")
	if err != nil {
		return nil, err
	}
	d, err := sharedJwk.Bytes(jwk, "d")
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

	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: cfg.curve, X: new(big.Int).SetBytes(x), Y: new(big.Int).SetBytes(y)},
		D:         new(big.Int).SetBytes(d),
	}

	return &ecdsaPrivateKey{privateKey: privateKey, algorithm: algorithm, hash: cfg.hash, size: cfg.size, crvName: cfg.crvName}, nil
}

// ToJWK returns the private key as a JWK (RFC 7517/7518), including the
// private scalar d alongside the public point x/y.
func (k *ecdsaPrivateKey) ToJWK() (map[string]any, error) {
	encode := func(b *big.Int) string {
		out := make([]byte, k.size)
		b.FillBytes(out)
		return base64.RawURLEncoding.EncodeToString(out)
	}

	return map[string]any{
		"kty": "EC",
		"alg": k.algorithm,
		"crv": k.crvName,
		"x":   encode(k.privateKey.X),
		"y":   encode(k.privateKey.Y),
		"d":   encode(k.privateKey.D),
	}, nil
}

func parseECDSAPrivateKey(der []byte) (*ecdsa.PrivateKey, error) {
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		ecKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an ecdsa private key")
		}
		return ecKey, nil
	}
	return x509.ParseECPrivateKey(der)
}

func (k *ecdsaPrivateKey) SignBytes(data []byte) ([]byte, error) {
	hasher := k.hash.New()
	hasher.Write(data)
	r, s, err := ecdsa.Sign(rand.Reader, k.privateKey, hasher.Sum(nil))
	if err != nil {
		return nil, err
	}
	out := make([]byte, 2*k.size)
	r.FillBytes(out[:k.size])
	s.FillBytes(out[k.size:])
	return out, nil
}

func (k *ecdsaPrivateKey) Algorithm() string {
	return k.algorithm
}

func (k *ecdsaPrivateKey) SignString(signingInput string) (string, error) {
	return sharedJwt.SignCompactToken(k.SignBytes, signingInput)
}
