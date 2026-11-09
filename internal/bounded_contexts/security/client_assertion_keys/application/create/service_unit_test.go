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
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	securityEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func newService(repo client_assertion_keys.Repository, bus event.Bus, projection *client_projection.ClientProjection) *create.Service {
	return create.NewService(&fakeClientProjectionRepository{projection: projection}, repo, bus)
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

func TestCreateClientAssertionKeyService_HappyPath(t *testing.T) {
	t.Parallel()

	adminId := ids2.NewUserId()
	tenantId := ids2.NewUserId()

	tests := []struct {
		name          string
		callerUserId  ids2.UserId
		callerAdminId ids2.UserId
	}{
		{
			name:          "owning tenant creates a key for its own client",
			callerUserId:  tenantId,
			callerAdminId: adminId,
		},
		{
			name:          "owning admin creates a key for its tenant's client",
			callerUserId:  adminId,
			callerAdminId: adminId,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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

			svc := newService(repo, bus, aClientProjection(clientId, adminId, tenantId))

			cmd := create.NewCommand(cakId, clientId, tc.callerUserId, tc.callerAdminId, "RS256", validPubKeyHex)
			require.NoError(t, svc.Create(context.Background(), cmd))

			require.NotNil(t, savedCAK)
			p := savedCAK.ToPrimitives()
			assert.Equal(t, cakId.String(), p.Id)
			assert.Equal(t, clientId.String(), p.ClientId)
			// Ownership is the CLIENT's owner, not the caller's own ids.
			assert.Equal(t, tenantId.String(), p.UserId)
			assert.Equal(t, adminId.String(), p.AdminId)
			assert.NotEmpty(t, p.PublicKey, "saved key must have a non-empty public key")

			require.Len(t, publishedEvts, 1)
			evt, ok := publishedEvts[0].(securityEvents.ClientAssertionKeyCreated)
			require.True(t, ok, "expected ClientAssertionKeyCreated event")
			assert.Equal(t, clientId.String(), evt.ClientId)
			assert.Equal(t, cakId.String(), evt.AggregateID())
		})
	}
}

func TestCreateClientAssertionKeyService_ErrorCases(t *testing.T) {
	t.Parallel()

	adminId := ids2.NewUserId()
	tenantId := ids2.NewUserId()
	otherAdminId := ids2.NewUserId()

	tests := []struct {
		name          string
		callerUserId  ids2.UserId
		callerAdminId ids2.UserId
		publicKey     string
		noClient      bool
		setupRepo     func(repo *cakmocks.MockClientAssertionKeyRepository)
		wantErrCode   app_errors.ErrorCode
	}{
		{
			name:          "returns not found when client does not exist",
			callerUserId:  tenantId,
			callerAdminId: adminId,
			publicKey:     validPubKeyHex,
			noClient:      true,
			wantErrCode:   client_assertion_keys.ErrClientAssertionKeyNotCreated,
		},
		{
			name:          "rejects empty public key",
			callerUserId:  tenantId,
			callerAdminId: adminId,
			publicKey:     "",
			wantErrCode:   client_assertion_keys.ErrClientAssertionKeyNotCreated,
		},
		{
			name:          "wraps repository save error",
			callerUserId:  tenantId,
			callerAdminId: adminId,
			publicKey:     validPubKeyHex,
			setupRepo: func(repo *cakmocks.MockClientAssertionKeyRepository) {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErrCode: client_assertion_keys.ErrClientAssertionKeyNotSaved,
		},
		{
			name:          "returns forbidden when client belongs to an unrelated user",
			callerUserId:  ids2.NewUserId(),
			callerAdminId: ids2.NewUserId(),
			publicKey:     validPubKeyHex,
			wantErrCode:   client_assertion_keys.ErrInvalidOwnership,
		},
		{
			name:          "returns forbidden for a sibling tenant under the same admin",
			callerUserId:  ids2.NewUserId(),
			callerAdminId: adminId,
			publicKey:     validPubKeyHex,
			wantErrCode:   client_assertion_keys.ErrInvalidOwnership,
		},
		{
			name:          "returns forbidden for a different admin",
			callerUserId:  otherAdminId,
			callerAdminId: otherAdminId,
			publicKey:     validPubKeyHex,
			wantErrCode:   client_assertion_keys.ErrInvalidOwnership,
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
				projection = aClientProjection(clientId, adminId, tenantId)
			}

			svc := newService(repo, bus, projection)

			cmd := create.NewCommand(newCakId(t), clientId, tc.callerUserId, tc.callerAdminId, "RS256", tc.publicKey)
			err := svc.Create(context.Background(), cmd)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
