package create_test

import (
	"context"
	"errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/_shared/client_projection"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// mockSecretHasher stubs client_secrets.SecretHasher.
type mockSecretHasher struct {
	hash string
	err  error
}

func (m *mockSecretHasher) Hash(_ string) (string, error) { return m.hash, m.err }

// ── Client projection Mother ──────────────────────────────────────────────────

// fakeClientProjectionRepository stubs client_projection.Repository. A nil
// projection makes Find fail as the real repository does for unknown clients.
type fakeClientProjectionRepository struct {
	projection *client_projection.ClientProjection
}

func (f *fakeClientProjectionRepository) Find(_ context.Context, _ ids2.ClientId) (*client_projection.ClientProjection, error) {
	if f.projection == nil {
		return nil, errors.New("this client projection does not exist")
	}
	return f.projection, nil
}

func aClientProjection(clientId ids2.ClientId, subscriptionId ids2.SubscriptionId) *client_projection.ClientProjection {
	return &client_projection.ClientProjection{
		ClientId:       clientId.String(),
		SubscriptionId: subscriptionId.String(),
	}
}
