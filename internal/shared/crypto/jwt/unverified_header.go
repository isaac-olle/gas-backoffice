package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// ExtractKid reads the kid header of a JWT without verifying its signature,
// so callers can pick which key to verify against.
func ExtractKid(tokenString string) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("parse token header: %w", err)
	}
	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return "", fmt.Errorf("token missing kid header")
	}
	return kid, nil
}

// AudienceContains reports whether the aud claim (a string or a list of
// strings, per the JWT spec) contains the given value.
func AudienceContains(aud any, value string) bool {
	switch v := aud.(type) {
	case string:
		return v == value
	case []string:
		for _, a := range v {
			if a == value {
				return true
			}
		}
	case []any:
		for _, a := range v {
			if s, ok := a.(string); ok && s == value {
				return true
			}
		}
	}
	return false
}
