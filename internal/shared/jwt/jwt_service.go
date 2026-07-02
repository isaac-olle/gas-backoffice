package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenGenerator interface {
	Generate(claims map[string]any, key Key, expiresIn time.Duration) (string, error)
}

type TokenVerifier interface {
	Verify(tokenString string, key Key) (map[string]any, error)
}

type JwtService struct {
	allowedMethods map[string]jwt.SigningMethod
}

func NewJwtService(allowedMethods ...jwt.SigningMethod) *JwtService {
	m := make(map[string]jwt.SigningMethod, len(allowedMethods))
	for _, method := range allowedMethods {
		m[method.Alg()] = method
	}
	return &JwtService{allowedMethods: m}
}

func (s *JwtService) Generate(claims map[string]any, key Key, expiresIn time.Duration) (string, error) {
	for _, required := range []string{"iss", "sub"} {
		if v, ok := claims[required]; !ok || v == "" {
			return "", fmt.Errorf("missing required claim: %s", required)
		}
	}

	signingMethod, err := s.resolveMethod(key.Method())
	if err != nil {
		return "", err
	}

	parsedKey, err := parsePrivateKey(key)
	if err != nil {
		return "", err
	}

	now := time.Now()
	jwtClaims := jwt.MapClaims{}
	for k, v := range claims {
		jwtClaims[k] = v
	}
	jwtClaims["iat"] = now.Unix()
	jwtClaims["exp"] = now.Add(expiresIn).Unix()

	token := jwt.NewWithClaims(signingMethod, jwtClaims)
	token.Header["kid"] = key.ID()
	return token.SignedString(parsedKey)
}

func (s *JwtService) Verify(tokenString string, key Key) (map[string]any, error) {
	if _, err := s.resolveMethod(key.Method()); err != nil {
		return nil, err
	}

	parsedKey, err := parsePublicKey(key)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != key.Method() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return parsedKey, nil
	})
	if err != nil {
		return nil, err
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return map[string]any(mapClaims), nil
}

func (s *JwtService) resolveMethod(method string) (jwt.SigningMethod, error) {
	m, ok := s.allowedMethods[method]
	if !ok {
		return nil, fmt.Errorf("signing method %q is not allowed", method)
	}
	return m, nil
}

func parsePrivateKey(key Key) (any, error) {
	block, _ := pem.Decode(key.Bytes())
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key PEM")
	}
	switch {
	case key.IsRSA():
		if k, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			return k, nil
		}
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case key.IsECDSA():
		if k, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			return k, nil
		}
		return x509.ParseECPrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported signing method: %s", key.Method())
	}
}

func parsePublicKey(key Key) (any, error) {
	block, _ := pem.Decode(key.Bytes())
	if block == nil {
		return nil, fmt.Errorf("failed to decode public key PEM")
	}
	if !key.IsRSA() && !key.IsECDSA() {
		return nil, fmt.Errorf("unsupported signing method: %s", key.Method())
	}
	k, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	return k, nil
}
