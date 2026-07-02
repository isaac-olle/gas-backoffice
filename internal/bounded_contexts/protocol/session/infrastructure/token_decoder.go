package infrastructure

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenDecoder struct{}

func NewJwtTokenDecoder() *JwtTokenDecoder {
	return &JwtTokenDecoder{}
}

func (d *JwtTokenDecoder) DecodeKid(token string) (string, error) {
	unverified, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("cannot parse token: %w", err)
	}
	kid, ok := unverified.Header["kid"].(string)
	if !ok || kid == "" {
		return "", fmt.Errorf("token missing kid header")
	}
	return kid, nil
}
