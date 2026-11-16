package infrastructure

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	generateJwtByPrId "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt_by_pr_id"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type BusIDTokenGenerator struct {
	queryBus query.Bus
}

func NewBusIDTokenGenerator(queryBus query.Bus) *BusIDTokenGenerator {
	return &BusIDTokenGenerator{queryBus: queryBus}
}

func (g *BusIDTokenGenerator) GenerateIDToken(
	ctx context.Context,
	claims session.TokenClaims,
) (string, error) {
	qry := generateJwtByPrId.NewQuery(claims.ProtectedResourceId, claims.Subject, session.AccessTokenTTL).
		WithAudience(claims.ClientID)
	if claims.HasScope("email") {
		qry = qry.WithExtraClaim("email", claims.Email)
	}
	if claims.HasScope("profile") && claims.Picture != "" {
		qry = qry.WithExtraClaim("picture", claims.Picture)
	}
	result, err := g.queryBus.Dispatch(ctx, qry)
	if err != nil {
		return "", err
	}
	token, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected result type from generate_jwt_by_pr_id query")
	}
	return token, nil
}
