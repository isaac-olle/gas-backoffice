package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/buses/query"

	"github.com/golang-jwt/jwt/v5"
	verifyJwt "github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
)

type BusTokenVerifier struct {
	queryBus query.Bus
}

func NewBusTokenVerifier(queryBus query.Bus) *BusTokenVerifier {
	return &BusTokenVerifier{queryBus: queryBus}
}

func (v *BusTokenVerifier) VerifyJWT(ctx context.Context, keyPairId string, token string) (jwt.MapClaims, error) {
	result, err := v.queryBus.Dispatch(ctx, verifyJwt.NewQuery(token))
	if err != nil {
		return nil, err
	}
	claims, ok := result.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from verify_jwt query")
	}
	return claims, nil
}
