package jwk

import (
	"encoding/base64"
	"fmt"
)

// String reads a required string member from a JWK map.
func String(jwk map[string]any, member string) (string, error) {
	raw, ok := jwk[member]
	if !ok {
		return "", fmt.Errorf("jwk missing %q member", member)
	}
	s, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("jwk member %q is not a string", member)
	}
	return s, nil
}

// Bytes reads a required base64url-encoded member from a JWK map.
func Bytes(jwk map[string]any, member string) ([]byte, error) {
	s, err := String(jwk, member)
	if err != nil {
		return nil, err
	}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("jwk member %q is not valid base64url: %w", member, err)
	}
	return b, nil
}
