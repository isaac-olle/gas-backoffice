//go:build integration

package infrastructure_test

import (
	"context"
	"testing"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	signing "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt_by_pr_id"
	keyservices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/domain_services"
	keymocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type tokenSigningBus struct{ handler query.Handler }

func (b tokenSigningBus) Dispatch(ctx context.Context, qry query.Query) (any, error) {
	return b.handler.Handle(ctx, qry)
}
func (tokenSigningBus) Register(query.Type, query.Handler) {}

type capturedTokenSigner struct{ spec key_pairs.TokenSpec }

func (s *capturedTokenSigner) Generate(_ context.Context, _ string, spec key_pairs.TokenSpec, _ *key_pairs.KeyPair, _ time.Duration) (string, error) {
	s.spec = spec
	return "signed-token", nil
}

// Exercise the protocol generator and security query handler together so the
// claims reaching the signer are checked without a database or private keys.
func TestIntegrationAccessTokenClaims(t *testing.T) {
	ctx := context.Background()
	repo := keymocks.NewMockRepository(t)
	repo.EXPECT().SearchByProtectedResourceId(mock.Anything, mock.Anything).Return(&key_pairs.KeyPair{}, nil)
	signer := &capturedTokenSigner{}
	bus := tokenSigningBus{handler: signing.NewQueryHandler(signing.NewService(keyservices.NewKeyPairFinder(repo), signer))}
	generator := infrastructure.NewBusAccessTokenGenerator(bus)
	claims := newResolvedTokenClaims(t)
	_, err := generator.GenerateAccessToken(ctx, claims)
	require.NoError(t, err)
	firstSubject := signer.spec.Sub
	assert.NotEmpty(t, firstSubject)
	assert.Equal(t, claims.Subject, firstSubject)
	assert.NotEqual(t, claims.ClientID, firstSubject)
	assert.Equal(t, claims.ProtectedResourceId, signer.spec.Aud)
	assert.Equal(t, claims.IdpID, signer.spec.ExtraClaims["idp_id"])
	assert.Equal(t, claims.ClientID, signer.spec.ExtraClaims["client_id"])
	assert.Equal(t, claims.Scope, signer.spec.ExtraClaims["scope"])
	_, err = generator.GenerateAccessToken(ctx, claims)
	require.NoError(t, err)
	assert.Equal(t, firstSubject, signer.spec.Sub)
	claims.IdpID = ids.NewIdentityProviderId().String()
	_, err = generator.GenerateAccessToken(ctx, claims)
	require.NoError(t, err)
	assert.Equal(t, firstSubject, signer.spec.Sub)
	assert.Equal(t, claims.IdpID, signer.spec.ExtraClaims["idp_id"])
	claims.IdpID = ""
	_, err = generator.GenerateAccessToken(ctx, claims)
	require.NoError(t, err)
	assert.NotContains(t, signer.spec.ExtraClaims, "idp_id")

	cc, err := session.NewClientCredentialsSession(ids.NewSessionId(), claims.ClientID, "read", claims.ProtectedResourceId)
	require.NoError(t, err)
	ccClaims := cc.TokenClaims()
	_, err = generator.GenerateAccessToken(ctx, ccClaims)
	require.NoError(t, err)
	assert.Equal(t, cc.ClientID().String(), signer.spec.Sub)
	assert.Equal(t, cc.ClientID().String(), signer.spec.ExtraClaims["client_id"])
	assert.Equal(t, cc.ProtectedResourceId(), signer.spec.Aud)
	assert.NotContains(t, signer.spec.ExtraClaims, "idp_id")
}
