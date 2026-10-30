package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	generateJwtByPrId "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt_by_pr_id"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
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
	qry := g.sessionTokenQuery(p, time.Until(p.ExpiresAt))
	return g.dispatch(ctx, qry, "generate jwt query")
}

func (g *BusTokenGenerator) GenerateRefreshToken(ctx context.Context, session *auth.Session) (string, error) {
	p := session.ToPrimitives()
	qry := g.sessionTokenQuery(p, 30*24*time.Hour).
		WithExtraClaim("typ", "backoffice_refresh")
	return g.dispatch(ctx, qry, "generate refresh jwt query")
}

func (g *BusTokenGenerator) sessionTokenQuery(p *auth.SessionPrimitives, expiresIn time.Duration) *generateJwtByPrId.Query {
	qry := generateJwtByPrId.NewQuery(g.protectedResourceId, p.UserId, expiresIn).
		WithAudience(g.protectedResourceId).
		WithExtraClaim("sid", p.SessionId).
		WithExtraClaim("system_role", roleClaim(p.Role)).
		WithExtraClaim("picture", p.Picture)

	// The authentication middleware expects a stable token contract. System
	// administrators have global access through their role, so they carry an
	// empty memberships list rather than omitting the claim.
	memberships := p.Memberships
	if memberships == nil {
		memberships = []vo.SubscriptionMembership{}
	}
	qry.WithExtraClaim("memberships", memberships)
	return qry
}

func (g *BusTokenGenerator) dispatch(ctx context.Context, qry *generateJwtByPrId.Query, description string) (string, error) {
	result, err := g.queryBus.Dispatch(ctx, qry)
	if err != nil {
		return "", err
	}
	token, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected result type from %s", description)
	}
	return token, nil
}

func roleClaim(role string) string {
	if role == vo.SystemAdministratorRole.String() {
		return "system_administrator"
	}
	return "subscription_user"
}
