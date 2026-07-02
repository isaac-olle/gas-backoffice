package infrastructure

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtDecoder struct {
	parser *jwt.Parser
}

func NewJwtDecoder() *JwtDecoder {
	return &JwtDecoder{parser: jwt.NewParser()}
}

func (d *JwtDecoder) ExtractKeyId(_ context.Context, assertion string) (string, error) {
	unverified, _, err := d.parser.ParseUnverified(assertion, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	kid, ok := unverified.Header["kid"].(string)
	if !ok {
		return "", fmt.Errorf("assertion JWT missing kid header")
	}
	return kid, nil
}
