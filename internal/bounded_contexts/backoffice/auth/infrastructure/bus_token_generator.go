package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	generateJwtByPrId "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt_by_pr_id"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type BusTokenGenerator struct {
	queryBus            query.Bus
	protectedResourceId string
}

func NewBusTokenGenerator(queryBus query.Bus, protectedResourceId string) *BusTokenGenerator {
	return &BusTokenGenerator{queryBus: queryBus, protectedResourceId: protectedResourceId}
}

func (g *BusTokenGenerator) Generate(ctx context.Context, session *auth.Session) (string, error) {
	p := session.ToPrimitives()
	expiresIn := time.Until(p.ExpiresAt)

	qry := generateJwtByPrId.NewQuery(g.protectedResourceId, p.UserId, expiresIn).
		WithAudience(g.protectedResourceId).
		WithExtraClaim("sid", p.SessionId).
		WithExtraClaim("role", p.Role).
		WithExtraClaim("admin_id", p.AdminId).
		WithExtraClaim("subscription_ids", p.SubscriptionIds)
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
