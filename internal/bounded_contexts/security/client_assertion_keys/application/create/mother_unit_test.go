package create_test

import (
	"context"
	"errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/_shared/client_projection"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// validPubKeyHex is a valid PEM-encoded RSA 2048-bit PKIX public key for tests.
const validPubKeyHex = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArRj90RkbqGVRcp1xYi63
L+ZKaqMi3IVFhm6adSM+SNJsdiEkcMjrXu8iKlxN19kfJFN/fDlKwZuDgdKLq2/X
QqKpQHItF6kgNbuzfkNWbLkzebWn5sokC9BqgStH4Inob5Plg8T4v2Po7JIn8098
bGtxNPLObG/+heX7JQNZDyocWpQWD4G4BrvgimfKGH4azvBmLHR4dayCrii7GR2I
W1kY4MmhX84TNLCqFnzZdd0GswyGJHVnyrpUvEJFaHmr+97I8vDBvFPJ2w2dEXN1
+BC5dJHqmLjr6qV+6Dqq56Gn8DcestRxtiS20np4Z2cMSCgqtgp1l1zHSjMO/rGX
WQIDAQAB
-----END PUBLIC KEY-----
`

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
