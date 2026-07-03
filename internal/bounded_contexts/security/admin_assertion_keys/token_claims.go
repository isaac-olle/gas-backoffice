package admin_assertion_keys

import (
	"fmt"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/actors"
)

type TokenClaims struct {
	sub string
	iss string
	aud string
	raw map[string]any
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

	return &TokenClaims{
		sub: sub,
		iss: iss,
		aud: aud,
		raw: claims,
	}, nil
}

func (t *TokenClaims) Sub() string { return t.sub }

func (t *TokenClaims) CheckIssuer(key *AdminAssertionKey) error {
	if t.iss != key.ownership.AdminId.String() {
		return fmt.Errorf(
			"invalid issuer: expected '%s', got '%s'",
			key.id.String(),
			t.iss,
		)
	}
	return nil
}

func (t *TokenClaims) CheckAudience(expectedValue string) error {
	if t.aud != expectedValue {
		return fmt.Errorf(
			"invalid audience: expected '%s', got '%s'",
			expectedValue,
			t.aud,
		)
	}
	return nil
}

func (t *TokenClaims) ToActor() (vo2.Actor, error) {
	return actors.NewActorForRole(vo2.AdminRole, t.sub)
}
