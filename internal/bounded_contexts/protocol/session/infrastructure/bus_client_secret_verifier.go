package infrastructure

import (
	"context"

	verifySecret "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/application/verify_secret"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

// BusClientSecretVerifier satisfies authenticate_client.ClientSecretVerifier
// by dispatching to the security bounded context, which is cross-BC and so
// stays on the query bus.
type BusClientSecretVerifier struct {
	queryBus query.Bus
}

func NewBusClientSecretVerifier(queryBus query.Bus) *BusClientSecretVerifier {
	return &BusClientSecretVerifier{queryBus: queryBus}
}

func (v *BusClientSecretVerifier) Verify(ctx context.Context, clientID, secret string) error {
	_, err := v.queryBus.Dispatch(ctx, verifySecret.NewQuery(clientID, secret))
	return err
}
