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
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	securityEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func newService(repo client_secrets.Repository, hasher client_secrets.SecretHasher, bus event.Bus, projection *client_projection.ClientProjection) *create.Service {
	return create.NewService(&fakeClientProjectionRepository{projection: projection}, repo, hasher, bus)
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

func TestCreateClientSecretService_HappyPath(t *testing.T) {
	t.Parallel()

	const plainSecret = "plaintext-secret"
	const hashedSecret = "$argon2id$hashed"

	adminId := ids2.NewUserId()
	tenantId := ids2.NewUserId()

	tests := []struct {
		name          string
		callerUserId  ids2.UserId
		callerAdminId ids2.UserId
	}{
		{
			name:          "owning tenant creates a secret for its own client",
			callerUserId:  tenantId,
			callerAdminId: adminId,
		},
		{
			name:          "owning admin creates a secret for its tenant's client",
			callerUserId:  adminId,
			callerAdminId: adminId,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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

			svc := newService(repo, &mockSecretHasher{hash: hashedSecret}, bus, aClientProjection(clientId, adminId, tenantId))

			cmd := create.NewCommand(id, clientId, tc.callerUserId, tc.callerAdminId, plainSecret)
			require.NoError(t, svc.Create(context.Background(), cmd))

			require.NotNil(t, savedCS)
			p := savedCS.ToPrimitives()
			assert.Equal(t, id.String(), p.Id)
			assert.Equal(t, clientId.String(), p.ClientId)
			// Ownership is the CLIENT's owner, not the caller's own ids.
			assert.Equal(t, tenantId.String(), p.UserId)
			assert.Equal(t, adminId.String(), p.AdminId)
			assert.Equal(t, hashedSecret, p.SecretHash)

			require.Len(t, publishedEvts, 1)
			evt, ok := publishedEvts[0].(securityEvents.ClientSecretCreated)
			require.True(t, ok, "expected ClientSecretCreated event")
			assert.Equal(t, clientId.String(), evt.ClientId)
			assert.Equal(t, id.String(), evt.AggregateID())
		})
	}
}

func TestCreateClientSecretService_ErrorCases(t *testing.T) {
	t.Parallel()

	adminId := ids2.NewUserId()
	tenantId := ids2.NewUserId()
	otherAdminId := ids2.NewUserId()

	tests := []struct {
		name          string
		callerUserId  ids2.UserId
		callerAdminId ids2.UserId
		plainSecret   string
		hasher        *mockSecretHasher
		noClient      bool
		setupRepo     func(repo *csmocks.MockRepository)
		wantErrCode   app_errors.ErrorCode
	}{
		{
			name:          "returns not found when client does not exist",
			callerUserId:  tenantId,
			callerAdminId: adminId,
			plainSecret:   "plain",
			hasher:        &mockSecretHasher{hash: "hash"},
			noClient:      true,
			wantErrCode:   client_secrets.ErrClientSecretNotCreated,
		},
		{
			name:          "rejects empty plain secret",
			callerUserId:  tenantId,
			callerAdminId: adminId,
			plainSecret:   "",
			hasher:        &mockSecretHasher{hash: "hash"},
			wantErrCode:   client_secrets.ErrClientSecretNotCreated,
		},
		{
			name:          "propagates hasher error",
			callerUserId:  tenantId,
			callerAdminId: adminId,
			plainSecret:   "plain",
			hasher:        &mockSecretHasher{err: errors.New("hash failed")},
			wantErrCode:   client_secrets.ErrClientSecretNotCreated,
		},
		{
			name:          "wraps repository save error",
			callerUserId:  tenantId,
			callerAdminId: adminId,
			plainSecret:   "plain",
			hasher:        &mockSecretHasher{hash: "hash"},
			setupRepo: func(repo *csmocks.MockRepository) {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErrCode: client_secrets.ErrClientSecretNotSaved,
		},
		{
			name:          "returns forbidden when client belongs to an unrelated user",
			callerUserId:  ids2.NewUserId(),
			callerAdminId: ids2.NewUserId(),
			plainSecret:   "plain",
			hasher:        &mockSecretHasher{hash: "hash"},
			wantErrCode:   client_secrets.ErrInvalidOwnership,
		},
		{
			name:          "returns forbidden for a sibling tenant under the same admin",
			callerUserId:  ids2.NewUserId(),
			callerAdminId: adminId,
			plainSecret:   "plain",
			hasher:        &mockSecretHasher{hash: "hash"},
			wantErrCode:   client_secrets.ErrInvalidOwnership,
		},
		{
			name:          "returns forbidden for a different admin",
			callerUserId:  otherAdminId,
			callerAdminId: otherAdminId,
			plainSecret:   "plain",
			hasher:        &mockSecretHasher{hash: "hash"},
			wantErrCode:   client_secrets.ErrInvalidOwnership,
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
				projection = aClientProjection(clientId, adminId, tenantId)
			}

			svc := newService(repo, tc.hasher, bus, projection)

			cmd := create.NewCommand(newSecretId(t), clientId, tc.callerUserId, tc.callerAdminId, tc.plainSecret)
			err := svc.Create(context.Background(), cmd)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
