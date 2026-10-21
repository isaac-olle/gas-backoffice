package infrastructure

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type PEMKeyGenerator struct{}

func NewPEMKeyGenerator() *PEMKeyGenerator {
	return &PEMKeyGenerator{}
}

func (g *PEMKeyGenerator) Generate(algorithm string) (string, string, error) {
	switch algorithm {
	case "RS256", "RS384", "RS512":
		return generateRSAKeyPair()
	case "ES256", "ES384", "ES512":
		return generateECKeyPair(algorithm)
	default:
		return "", "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func generateRSAKeyPair() (string, string, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}
	return encodeKeyPairPEM(privKey, &privKey.PublicKey)
}

func generateECKeyPair(algorithm string) (string, string, error) {
	var curve elliptic.Curve
	switch algorithm {
	case "ES256":
		curve = elliptic.P256()
	case "ES384":
		curve = elliptic.P384()
	case "ES512":
		curve = elliptic.P521()
	}

	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return "", "", err
	}
	return encodeKeyPairPEM(privKey, &privKey.PublicKey)
}

func encodeKeyPairPEM(privKey, pubKey any) (string, string, error) {
	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", "", err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", "", err
	}

	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})

	return string(privPEM), string(pubPEM), nil
}
