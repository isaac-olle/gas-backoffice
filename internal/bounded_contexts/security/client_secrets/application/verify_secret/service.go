package verify_secret

import (
	"context"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets"
)

type Service struct {
	secretVerifier client_secrets.SecretVerifier
}

func NewService(secretVerifier client_secrets.SecretVerifier) *Service {
	return &Service{secretVerifier: secretVerifier}
}

func (s *Service) VerifySecret(ctx context.Context, qry *Query) (any, error) {
	return nil, s.secretVerifier.Verify(ctx, qry.clientID, qry.plainSecret)
}
