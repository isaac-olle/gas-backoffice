package private_keys

import (
	"crypto"
	"crypto/rand"
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

type rsaPrivateKey struct {
	privateKey *rsa.PrivateKey
	pemBytes   []byte
	algorithm  string
	hash       crypto.Hash
}

var _ sharedCrypto.SigningKey = (*rsaPrivateKey)(nil)

func NewRSAPrivateKeyFromPEM(pemBytes []byte, algorithm string) (*rsaPrivateKey, error) {
	hash, ok := rsaHashes[algorithm]
	if !ok {
		return nil, fmt.Errorf("unsupported rsa algorithm: %s", algorithm)
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key PEM")
	}
	privateKey, err := parseRSAPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &rsaPrivateKey{privateKey: privateKey, pemBytes: pemBytes, algorithm: algorithm, hash: hash}, nil
}

func (k *rsaPrivateKey) ToPEMFormat() ([]byte, error) {
	if k.pemBytes != nil {
		return k.pemBytes, nil
	}
	der, err := x509.MarshalPKCS8PrivateKey(k.privateKey)
	if err != nil {
		return nil, fmt.Errorf("marshal private key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), nil
}

// NewRSAPrivateKeyFromJWK builds a signing key from a JWK (RFC 7517/7518)
// that includes the private members (d, p, q).
func NewRSAPrivateKeyFromJWK(jwk map[string]any) (*rsaPrivateKey, error) {
	n, err := sharedJwk.Bytes(jwk, "n")
	if err != nil {
		return nil, err
	}
	e, err := sharedJwk.Bytes(jwk, "e")
	if err != nil {
		return nil, err
	}
	d, err := sharedJwk.Bytes(jwk, "d")
	if err != nil {
		return nil, err
	}
	p, err := sharedJwk.Bytes(jwk, "p")
	if err != nil {
		return nil, err
	}
	q, err := sharedJwk.Bytes(jwk, "q")
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

	privateKey := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: new(big.Int).SetBytes(n),
			E: int(new(big.Int).SetBytes(e).Int64()),
		},
		D:      new(big.Int).SetBytes(d),
		Primes: []*big.Int{new(big.Int).SetBytes(p), new(big.Int).SetBytes(q)},
	}
	privateKey.Precompute()

	return &rsaPrivateKey{privateKey: privateKey, algorithm: algorithm, hash: hash}, nil
}

// ToJWK returns the private key as a JWK (RFC 7517/7518), including the
// private members (d, p, q, dp, dq, qi) alongside n/e.
func (k *rsaPrivateKey) ToJWK() (map[string]any, error) {
	if len(k.privateKey.Primes) != 2 {
		return nil, fmt.Errorf("only 2-prime rsa keys are supported for jwk export")
	}
	k.privateKey.Precompute()

	encode := func(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

	return map[string]any{
		"kty": "RSA",
		"alg": k.algorithm,
		"n":   encode(k.privateKey.N.Bytes()),
		"e":   encode(big.NewInt(int64(k.privateKey.E)).Bytes()),
		"d":   encode(k.privateKey.D.Bytes()),
		"p":   encode(k.privateKey.Primes[0].Bytes()),
		"q":   encode(k.privateKey.Primes[1].Bytes()),
		"dp":  encode(k.privateKey.Precomputed.Dp.Bytes()),
		"dq":  encode(k.privateKey.Precomputed.Dq.Bytes()),
		"qi":  encode(k.privateKey.Precomputed.Qinv.Bytes()),
	}, nil
}

func parseRSAPrivateKey(der []byte) (*rsa.PrivateKey, error) {
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an rsa private key")
		}
		return rsaKey, nil
	}
	return x509.ParsePKCS1PrivateKey(der)
}

func (k *rsaPrivateKey) SignBytes(data []byte) ([]byte, error) {
	hasher := k.hash.New()
	hasher.Write(data)
	return rsa.SignPKCS1v15(rand.Reader, k.privateKey, k.hash, hasher.Sum(nil))
}

func (k *rsaPrivateKey) Algorithm() string {
	return k.algorithm
}

func (k *rsaPrivateKey) SignString(signingInput string) (string, error) {
	return sharedJwt.SignCompactToken(k.SignBytes, signingInput)
}
