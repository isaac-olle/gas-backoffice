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

func (d *JwtDecoder) ExtractKeyId(_ context.Context, token string) (string, error) {
	unverified, _, err := d.parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	kid, ok := unverified.Header["kid"].(string)
	if !ok {
		return "", fmt.Errorf("token missing kid header")
	}
	return kid, nil
}
