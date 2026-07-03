package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys/application/create"
	aakmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	securityEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestCreateAdminAssertionKeyService_HappyPath(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	adminId := uuid.New().String()

	var savedAAK *admin_assertion_keys.AdminAssertionKey
	repo := aakmocks.NewMockRepository(t)
	repo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, aak *admin_assertion_keys.AdminAssertionKey) {
		savedAAK = aak
	}).Return(nil)

	var publishedEvts []event.Event
	bus := eventmocks.NewMockBus(t)
	bus.EXPECT().Publish(mock.Anything, mock.Anything).Run(func(_ context.Context, evts []event.Event) {
		publishedEvts = evts
	})

	svc := create.NewService(repo, bus)

	cmd := create.NewCommand(id, adminId, "RS256", stubPubKeyHex)
	require.NoError(t, svc.Create(context.Background(), cmd))

	require.NotNil(t, savedAAK)
	p := savedAAK.ToPrimitives()
	assert.Equal(t, id, p.Id)
	assert.Equal(t, adminId, p.AdminId)
	assert.NotEmpty(t, p.PublicKey, "saved key must have a non-empty public key")

	require.Len(t, publishedEvts, 1)
	evt, ok := publishedEvts[0].(securityEvents.AdminAssertionKeyCreated)
	require.True(t, ok, "expected AdminAssertionKeyCreated event")
	assert.Equal(t, adminId, evt.AdminId)
	assert.Equal(t, id, evt.AggregateID())
}

func TestCreateAdminAssertionKeyService_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		id          string
		adminId     string
		publicKey   string
		setupRepo   func(repo *aakmocks.MockRepository)
		wantErrCode app_errors.ErrorCode
	}{
		{
			name:        "rejects invalid admin id",
			id:          uuid.New().String(),
			adminId:     "bad-admin-id",
			publicKey:   stubPubKeyHex,
			wantErrCode: admin_assertion_keys.ErrAdminAssertionKeyNotCreated,
		},
		{
			name:        "rejects empty public key",
			id:          uuid.New().String(),
			adminId:     uuid.New().String(),
			publicKey:   "",
			wantErrCode: admin_assertion_keys.ErrAdminAssertionKeyNotCreated,
		},
		{
			name:      "wraps repository save error",
			id:        uuid.New().String(),
			adminId:   uuid.New().String(),
			publicKey: stubPubKeyHex,
			setupRepo: func(repo *aakmocks.MockRepository) {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErrCode: admin_assertion_keys.ErrAdminAssertionKeyNotSaved,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := aakmocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			bus := eventmocks.NewMockBus(t)
			svc := create.NewService(repo, bus)

			cmd := create.NewCommand(tc.id, tc.adminId, "RS256", tc.publicKey)
			err := svc.Create(context.Background(), cmd)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
