package infrastructure

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	generateJwtByPrId "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt_by_pr_id"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type BusAccessTokenGenerator struct {
	queryBus query.Bus
}

func NewBusAccessTokenGenerator(queryBus query.Bus) *BusAccessTokenGenerator {
	return &BusAccessTokenGenerator{queryBus: queryBus}
}

func (g *BusAccessTokenGenerator) GenerateAccessToken(
	ctx context.Context,
	sess *session.BaseSession,
	identity *session.UserIdentityProjection,
) (string, error) {
	if identity.Subject == "" {
		return "", fmt.Errorf("cannot generate token: authenticated user not resolved")
	}
	claims := sess.ToPrimitives()
	qry := generateJwtByPrId.NewQuery(claims.ProtectedResourceId, claims.ClientID, session.AccessTokenTTL).
		WithExtraClaim("scope", claims.Scope)
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
