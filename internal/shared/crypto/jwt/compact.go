package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// SignCompactToken signs signingInput via signBytes and appends the base64
// signature, producing the final compact JWT. Shared by every SigningKey
// implementation's SignString.
func SignCompactToken(signBytes func([]byte) ([]byte, error), signingInput string) (string, error) {
	sig, err := signBytes([]byte(signingInput))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

// VerifyToken parses a JWT (header.payload.signature), delegating
// signature verification to verifyBytes. Used by VerifyingKey
// implementations to build their VerifyString.
func VerifyToken(verifyBytes func(data, signature []byte) error, tokenString string) (map[string]any, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed token: expected 3 parts, got %d", len(parts))
	}

	signingInput := parts[0] + "." + parts[1]
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("decode signature: %w", err)
	}

	if err = verifyBytes([]byte(signingInput), sig); err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode claims: %w", err)
	}

	var claims map[string]any
	if err = json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("decode claims: %w", err)
	}

	return claims, nil
}
