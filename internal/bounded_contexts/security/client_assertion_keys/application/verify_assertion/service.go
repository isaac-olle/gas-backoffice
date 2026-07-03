package verify_assertion

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
)

type ClientClaims struct {
	ClientID string
}

type Service struct {
	assertionVerifier client_assertion_keys.AssertionVerifier
}

func NewService(assertionVerifier client_assertion_keys.AssertionVerifier) *Service {
	return &Service{assertionVerifier: assertionVerifier}
}

func (s *Service) VerifyAssertion(ctx context.Context, qry *Query) (any, error) {
	claims, err := s.assertionVerifier.Verify(ctx, qry.assertion)
	if err != nil {
		return nil, err
	}

	sub, _ := claims["sub"].(string)
	return &ClientClaims{ClientID: sub}, nil
}
