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
	sess *session.BaseSession,
	identity *session.UserIdentityProjection,
) (string, error) {
	if identity.Subject == "" {
		return "", fmt.Errorf("cannot generate token: authenticated user not resolved")
	}
	claims := sess.ToPrimitives()
	qry := generateJwtByPrId.NewQuery(claims.ProtectedResourceId, identity.Subject, session.AccessTokenTTL).
		WithAudience(claims.ClientID)
	if sess.HasScope("email") {
		qry = qry.WithExtraClaim("email", identity.Email)
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
