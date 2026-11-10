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
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/application/create"
	csmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/mocks"
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

func newService(repo client_secrets.Repository, hasher client_secrets.SecretHasher, bus event.Bus, projection *client_projection.ClientProjection) *create.Service {
	return create.NewService(&fakeClientProjectionRepository{projection: projection}, repo, hasher, bus, sharedDomainServices.NewSubscriptionAuthorizer())
}

func newClientId(t *testing.T) ids2.ClientId {
	t.Helper()
	id, err := ids2.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	return id
}

func newSecretId(t *testing.T) ids2.ClientSecretId {
	t.Helper()
	id, err := ids2.NewClientSecretIdFromString(uuid.New().String())
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

func TestCreateClientSecretService_HappyPath(t *testing.T) {
	t.Parallel()

	const plainSecret = "plaintext-secret"
	const hashedSecret = "$argon2id$hashed"

	subscriptionId := ids2.NewSubscriptionId()

	id := newSecretId(t)
	clientId := newClientId(t)

	var savedCS *client_secrets.ClientSecret
	repo := csmocks.NewMockRepository(t)
	repo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, cs *client_secrets.ClientSecret) {
		savedCS = cs
	}).Return(nil)

	var publishedEvts []event.Event
	bus := eventmocks.NewMockBus(t)
	bus.EXPECT().Publish(mock.Anything, mock.Anything).Run(func(_ context.Context, evts []event.Event) {
		publishedEvts = evts
	})

	svc := newService(repo, &mockSecretHasher{hash: hashedSecret}, bus, aClientProjection(clientId, subscriptionId))

	cmd := create.NewCommand(id, clientId, plainSecret, authorFor(t, subscriptionId))
	require.NoError(t, svc.Create(context.Background(), cmd))

	require.NotNil(t, savedCS)
	p := savedCS.ToPrimitives()
	assert.Equal(t, id.String(), p.Id)
	assert.Equal(t, clientId.String(), p.ClientId)
	// Ownership is the CLIENT's subscription, not necessarily the caller's own.
	assert.Equal(t, subscriptionId.String(), p.SubscriptionId)
	assert.Equal(t, hashedSecret, p.SecretHash)

	require.Len(t, publishedEvts, 1)
	evt, ok := publishedEvts[0].(securityEvents.ClientSecretCreated)
	require.True(t, ok, "expected ClientSecretCreated event")
	assert.Equal(t, clientId.String(), evt.ClientId)
	assert.Equal(t, id.String(), evt.AggregateID())
}

func TestCreateClientSecretService_ErrorCases(t *testing.T) {
	t.Parallel()

	subscriptionId := ids2.NewSubscriptionId()

	tests := []struct {
		name        string
		author      func(t *testing.T) *vo.UserAuthor
		plainSecret string
		hasher      *mockSecretHasher
		noClient    bool
		setupRepo   func(repo *csmocks.MockRepository)
		wantErrCode app_errors.ErrorCode
	}{
		{
			name:        "returns not found when client does not exist",
			author:      func(t *testing.T) *vo.UserAuthor { return authorFor(t, subscriptionId) },
			plainSecret: "plain",
			hasher:      &mockSecretHasher{hash: "hash"},
			noClient:    true,
			wantErrCode: client_secrets.ErrClientSecretNotCreated,
		},
		{
			name:        "returns forbidden when author lacks subscription access",
			author:      func(t *testing.T) *vo.UserAuthor { return authorFor(t, ids2.NewSubscriptionId()) },
			plainSecret: "plain",
			hasher:      &mockSecretHasher{hash: "hash"},
			wantErrCode: sharedDomainServices.SubscriptionAccessDenied,
		},
		{
			name:        "rejects empty plain secret",
			author:      func(t *testing.T) *vo.UserAuthor { return authorFor(t, subscriptionId) },
			plainSecret: "",
			hasher:      &mockSecretHasher{hash: "hash"},
			wantErrCode: client_secrets.ErrClientSecretNotCreated,
		},
		{
			name:        "propagates hasher error",
			author:      func(t *testing.T) *vo.UserAuthor { return authorFor(t, subscriptionId) },
			plainSecret: "plain",
			hasher:      &mockSecretHasher{err: errors.New("hash failed")},
			wantErrCode: client_secrets.ErrClientSecretNotCreated,
		},
		{
			name:        "wraps repository save error",
			author:      func(t *testing.T) *vo.UserAuthor { return authorFor(t, subscriptionId) },
			plainSecret: "plain",
			hasher:      &mockSecretHasher{hash: "hash"},
			setupRepo: func(repo *csmocks.MockRepository) {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErrCode: client_secrets.ErrClientSecretNotSaved,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := csmocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			bus := eventmocks.NewMockBus(t)

			clientId := newClientId(t)
			var projection *client_projection.ClientProjection
			if !tc.noClient {
				projection = aClientProjection(clientId, subscriptionId)
			}

			svc := newService(repo, tc.hasher, bus, projection)

			cmd := create.NewCommand(newSecretId(t), clientId, tc.plainSecret, tc.author(t))
			err := svc.Create(context.Background(), cmd)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
