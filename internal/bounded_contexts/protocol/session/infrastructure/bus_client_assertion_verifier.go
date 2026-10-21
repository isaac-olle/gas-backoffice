package infrastructure

import (
	"context"
	"fmt"

	verifyAssertion "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/verify_assertion"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

// BusClientAssertionVerifier satisfies
// authenticate_client.ClientAssertionVerifier by dispatching to the security
// bounded context, which is cross-BC and so stays on the query bus.
type BusClientAssertionVerifier struct {
	queryBus query.Bus
}

func NewBusClientAssertionVerifier(queryBus query.Bus) *BusClientAssertionVerifier {
	return &BusClientAssertionVerifier{queryBus: queryBus}
}

func (v *BusClientAssertionVerifier) Verify(ctx context.Context, assertion string) (string, error) {
	result, err := v.queryBus.Dispatch(ctx, verifyAssertion.NewQuery(assertion))
	if err != nil {
		return "", err
	}
	claims, ok := result.(*verifyAssertion.ClientClaims)
	if !ok {
		return "", fmt.Errorf("unexpected result type from verify_assertion query")
	}
	return claims.ClientID, nil
}
