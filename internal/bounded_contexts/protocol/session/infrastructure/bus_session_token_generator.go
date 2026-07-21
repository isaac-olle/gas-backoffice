package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	generateJwtByPrId "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt_by_pr_id"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type BusSessionTokenGenerator struct {
	queryBus query.Bus
	issuer   string
}

func NewBusSessionTokenGenerator(queryBus query.Bus, issuer string) *BusSessionTokenGenerator {
	return &BusSessionTokenGenerator{queryBus: queryBus, issuer: issuer}
}

func (g *BusSessionTokenGenerator) GenerateFromSession(ctx context.Context, sess *session.AuthorizationCodeSession, client *session.Client, protectedResourceId string) (string, error) {
	return g.dispatch(ctx, map[string]any{
		"iss":       g.issuer,
		"sub":       client.ClientID,
		"scope":     sess.ToPrimitives().Scope,
		"client_id": client.ClientID,
		"user_id":   client.UserId,
	}, protectedResourceId)
}

func (g *BusSessionTokenGenerator) GenerateFromClientCredentials(ctx context.Context, sess *session.ClientCredentialsSession, protectedResourceId string) (string, error) {
	p := sess.ToPrimitives()
	return g.dispatch(ctx, map[string]any{
		"iss":       g.issuer,
		"sub":       p.ClientID,
		"scope":     p.Scope,
		"client_id": p.ClientID,
		"user_id":   p.UserId,
	}, protectedResourceId)
}

func (g *BusSessionTokenGenerator) dispatch(ctx context.Context, claims map[string]any, protectedResourceId string) (string, error) {
	result, err := g.queryBus.Dispatch(ctx, generateJwtByPrId.NewQuery(protectedResourceId, claims, time.Hour))
	if err != nil {
		return "", err
	}
	token, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected result type from generate_jwt_by_pr_id query")
	}
	return token, nil
}
