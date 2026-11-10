package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/_shared/client_projection"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/create"
	cakmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/domain_services/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	securityEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func newService(repo client_assertion_keys.Repository, bus event.Bus, projection *client_projection.ClientProjection) *create.Service {
	return create.NewService(&fakeClientProjectionRepository{projection: projection}, repo, bus, sharedDomainServices.NewSubscriptionAuthorizer())
}

func newClientId(t *testing.T) ids2.ClientId {
	t.Helper()
	id, err := ids2.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	return id
}

func newCakId(t *testing.T) ids2.ClientAssertionKeyId {
	t.Helper()
	id, err := ids2.NewClientAssertionKeyIdFromString(uuid.New().String())
	require.NoError(t, err)
	return id
}

func authorFor(t *testing.T, subscriptionId ids2.SubscriptionId) *vo.UserAuthor {
	t.Helper()
	author, err := vo.NewUserAuthorWithMemberships(ids2.NewUserId().String(), ids2.NewUserId().String(), vo.SubscriptionUserRole,
		[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
	require.NoError(t, err)
	return author
}

func TestCreateClientAssertionKeyService_HappyPath(t *testing.T) {
	t.Parallel()

	subscriptionId := ids2.NewSubscriptionId()

	cakId := newCakId(t)
	clientId := newClientId(t)

	var savedCAK *client_assertion_keys.ClientAssertionKey
	repo := cakmocks.NewMockClientAssertionKeyRepository(t)
	repo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, cak *client_assertion_keys.ClientAssertionKey) {
		savedCAK = cak
	}).Return(nil)

	var publishedEvts []event.Event
	bus := eventmocks.NewMockBus(t)
	bus.EXPECT().Publish(mock.Anything, mock.Anything).Run(func(_ context.Context, evts []event.Event) {
		publishedEvts = evts
	})

	svc := newService(repo, bus, aClientProjection(clientId, subscriptionId))

	cmd := create.NewCommand(cakId, clientId, "RS256", validPubKeyHex, authorFor(t, subscriptionId))
	require.NoError(t, svc.Create(context.Background(), cmd))

	require.NotNil(t, savedCAK)
	p := savedCAK.ToPrimitives()
	assert.Equal(t, cakId.String(), p.Id)
	assert.Equal(t, clientId.String(), p.ClientId)
	// Ownership is the CLIENT's subscription, not necessarily the caller's own.
	assert.Equal(t, subscriptionId.String(), p.SubscriptionId)
	assert.NotEmpty(t, p.PublicKey, "saved key must have a non-empty public key")

	require.Len(t, publishedEvts, 1)
	evt, ok := publishedEvts[0].(securityEvents.ClientAssertionKeyCreated)
	require.True(t, ok, "expected ClientAssertionKeyCreated event")
	assert.Equal(t, clientId.String(), evt.ClientId)
	assert.Equal(t, cakId.String(), evt.AggregateID())
}

func TestCreateClientAssertionKeyService_ErrorCases(t *testing.T) {
	t.Parallel()

	subscriptionId := ids2.NewSubscriptionId()

	tests := []struct {
		name        string
		author      func(t *testing.T) *vo.UserAuthor
		publicKey   string
		noClient    bool
		setupRepo   func(repo *cakmocks.MockClientAssertionKeyRepository)
		wantErrCode app_errors.ErrorCode
	}{
		{
			name:        "returns not found when client does not exist",
			author:      func(t *testing.T) *vo.UserAuthor { return authorFor(t, subscriptionId) },
			publicKey:   validPubKeyHex,
			noClient:    true,
			wantErrCode: client_assertion_keys.ErrClientAssertionKeyNotCreated,
		},
		{
			name:        "returns forbidden when author lacks subscription access",
			author:      func(t *testing.T) *vo.UserAuthor { return authorFor(t, ids2.NewSubscriptionId()) },
			publicKey:   validPubKeyHex,
			wantErrCode: sharedDomainServices.SubscriptionAccessDenied,
		},
		{
			name:        "rejects empty public key",
			author:      func(t *testing.T) *vo.UserAuthor { return authorFor(t, subscriptionId) },
			publicKey:   "",
			wantErrCode: client_assertion_keys.ErrClientAssertionKeyNotCreated,
		},
		{
			name:      "wraps repository save error",
			author:    func(t *testing.T) *vo.UserAuthor { return authorFor(t, subscriptionId) },
			publicKey: validPubKeyHex,
			setupRepo: func(repo *cakmocks.MockClientAssertionKeyRepository) {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErrCode: client_assertion_keys.ErrClientAssertionKeyNotSaved,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := cakmocks.NewMockClientAssertionKeyRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			bus := eventmocks.NewMockBus(t)

			clientId := newClientId(t)
			var projection *client_projection.ClientProjection
			if !tc.noClient {
				projection = aClientProjection(clientId, subscriptionId)
			}

			svc := newService(repo, bus, projection)

			cmd := create.NewCommand(newCakId(t), clientId, "RS256", tc.publicKey, tc.author(t))
			err := svc.Create(context.Background(), cmd)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
