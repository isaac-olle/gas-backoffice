package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"time"

	generateJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt"
)

type BusTokenGenerator struct {
	queryBus query.Bus
}

func NewBusTokenGenerator(queryBus query.Bus) *BusTokenGenerator {
	return &BusTokenGenerator{queryBus: queryBus}
}

func (g *BusTokenGenerator) GenerateJWT(ctx context.Context, keyPairId string, claims map[string]any, expiresIn time.Duration) (string, error) {
	result, err := g.queryBus.Dispatch(ctx, generateJwt.NewQuery(keyPairId, claims, expiresIn))
	if err != nil {
		return "", err
	}
	token, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected result type from generate_jwt query")
	}
	return token, nil
}
