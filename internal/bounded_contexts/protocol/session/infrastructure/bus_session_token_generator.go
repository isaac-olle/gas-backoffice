package infrastructure

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	generateJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt"
)

type BusSessionTokenGenerator struct {
	queryBus     query.Bus
	signingKeyId string
	issuer       string
}

func NewBusSessionTokenGenerator(queryBus query.Bus, signingKeyId, issuer string) *BusSessionTokenGenerator {
	return &BusSessionTokenGenerator{queryBus: queryBus, signingKeyId: signingKeyId, issuer: issuer}
}

func (g *BusSessionTokenGenerator) GenerateFromSession(ctx context.Context, sess *session.AuthorizationCodeSession, client *session.Client) (string, error) {
	return g.dispatch(ctx, map[string]any{
		"iss":       g.issuer,
		"sub":       client.ClientID,
		"scope":     sess.ToPrimitives().Scope,
		"client_id": client.ClientID,
		"tenant_id": client.TenantID,
		"admin_id":  client.AdminID,
	})
}

func (g *BusSessionTokenGenerator) GenerateFromClientCredentials(ctx context.Context, sess *session.ClientCredentialsSession) (string, error) {
	p := sess.ToPrimitives()
	return g.dispatch(ctx, map[string]any{
		"iss":       g.issuer,
		"sub":       p.ClientID,
		"scope":     p.Scope,
		"client_id": p.ClientID,
		"tenant_id": p.TenantID,
		"admin_id":  p.AdminID,
	})
}

func (g *BusSessionTokenGenerator) dispatch(ctx context.Context, claims map[string]any) (string, error) {
	result, err := g.queryBus.Dispatch(ctx, generateJwt.NewQuery(g.signingKeyId, claims, time.Hour))
	if err != nil {
		return "", err
	}
	token, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected result type from generate_jwt query")
	}
	return token, nil
}
