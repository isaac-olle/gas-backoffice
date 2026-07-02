package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets/application/create"
	csmocks "github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets/mocks"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas/internal/shared/buses/event/mocks"
	securityEvents "github.com/isaac-olle/gas/internal/shared/events/security"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestCreateClientSecretService_HappyPath(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	clientId := uuid.New().String()
	const plainSecret = "plaintext-secret"
	const hashedSecret = "$argon2id$hashed"

	var savedCS *client_secrets.ClientSecret
	repo := csmocks.NewMockRepository(t)
	repo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, cs *client_secrets.ClientSecret) {
		savedCS = cs
	}).Return(nil)

	hasher := &mockSecretHasher{hash: hashedSecret}

	var publishedEvts []event.Event
	bus := eventmocks.NewMockBus(t)
	bus.EXPECT().Publish(mock.Anything, mock.Anything).Run(func(_ context.Context, evts []event.Event) {
		publishedEvts = evts
	})

	svc := create.NewService(repo, hasher, bus)

	cmd := create.NewCommand(id, clientId, plainSecret)
	require.NoError(t, svc.Create(context.Background(), cmd))

	require.NotNil(t, savedCS)
	p := savedCS.ToPrimitives()
	assert.Equal(t, id, p.Id)
	assert.Equal(t, clientId, p.ClientId)
	assert.Equal(t, hashedSecret, p.SecretHash)

	require.Len(t, publishedEvts, 1)
	evt, ok := publishedEvts[0].(securityEvents.ClientSecretCreated)
	require.True(t, ok, "expected ClientSecretCreated event")
	assert.Equal(t, clientId, evt.ClientId)
	assert.Equal(t, id, evt.AggregateID())
}

func TestCreateClientSecretService_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		id          string
		clientId    string
		plainSecret string
		setupRepo   func(repo *csmocks.MockRepository)
		setupHasher func(hasher *mockSecretHasher)
		wantErrCode app_errors.ErrorCode
	}{
		{
			name:        "rejects invalid client id",
			id:          uuid.New().String(),
			clientId:    "not-a-uuid",
			plainSecret: "plaintext",
			wantErrCode: client_secrets.ErrClientSecretNotCreated,
		},
		{
			name:        "rejects invalid secret id",
			id:          "not-a-uuid",
			clientId:    uuid.New().String(),
			plainSecret: "plain",
			setupHasher: func(hasher *mockSecretHasher) { hasher.hash = "somehash" },
			wantErrCode: client_secrets.ErrClientSecretNotCreated,
		},
		{
			name:        "propagates hasher error",
			id:          uuid.New().String(),
			clientId:    uuid.New().String(),
			plainSecret: "plain",
			setupHasher: func(hasher *mockSecretHasher) {
				hasher.err = errors.New("hash failed")
			},
			wantErrCode: client_secrets.ErrClientSecretNotCreated,
		},
		{
			name:        "wraps repository save error",
			id:          uuid.New().String(),
			clientId:    uuid.New().String(),
			plainSecret: "plain",
			setupHasher: func(hasher *mockSecretHasher) {
				hasher.hash = "hash"
			},
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
			hasher := &mockSecretHasher{}
			if tc.setupHasher != nil {
				tc.setupHasher(hasher)
			}
			bus := eventmocks.NewMockBus(t)
			svc := create.NewService(repo, hasher, bus)

			cmd := create.NewCommand(tc.id, tc.clientId, tc.plainSecret)
			err := svc.Create(context.Background(), cmd)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
