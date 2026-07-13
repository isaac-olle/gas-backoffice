package key_pairs

import (
	"fmt"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
)

type TokenClaims struct {
	sub  string
	iss  string
	aud  string
	role vo2.Role
	raw  map[string]any
}

func NewTokenClaims(claims map[string]any) (*TokenClaims, error) {
	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("claim 'sub' is invalid")
	}

	iss, ok := claims["iss"].(string)
	if !ok {
		return nil, fmt.Errorf("claim 'iss' is invalid")
	}

	aud, ok := claims["aud"].(string)
	if !ok {
		return nil, fmt.Errorf("claim 'aud' is invalid")
	}

	roleStr, _ := claims["role"].(string)
	role, _ := vo2.NewRole(roleStr)

	return &TokenClaims{sub: sub, iss: iss, aud: aud, role: role, raw: claims}, nil
}

func (t *TokenClaims) Sub() string         { return t.sub }
func (t *TokenClaims) Raw() map[string]any { return t.raw }

func (t *TokenClaims) CheckIssuer(expectedIss string) error {
	if t.iss != expectedIss {
		return fmt.Errorf("invalid issuer: expected '%s', got '%s'", expectedIss, t.iss)
	}
	return nil
}

func (t *TokenClaims) CheckAudience(expectedAud string) error {
	if t.aud != expectedAud {
		return fmt.Errorf("invalid audience: expected '%s', got '%s'", expectedAud, t.aud)
	}
	return nil
}

func (t *TokenClaims) ToAuthor() (vo2.Author, error) {
	return authors.NewAuthorForRole(t.role, t.sub)
}
