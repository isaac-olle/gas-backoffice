package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/create"
	keyPairInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
	kpmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/mocks"
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

func TestCreateKeyPairService_HappyPath(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	prId := uuid.New().String()

	var savedKP *key_pairs.KeyPair
	repo := kpmocks.NewMockRepository(t)
	repo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, kp *key_pairs.KeyPair) {
		savedKP = kp
	}).Return(nil)

	var publishedEvts []event.Event
	bus := eventmocks.NewMockBus(t)
	bus.EXPECT().Publish(mock.Anything, mock.Anything).Run(func(_ context.Context, evts []event.Event) {
		publishedEvts = evts
	})

	svc := create.NewService(repo, bus, keyPairInfra.NewPEMKeyGenerator())

	cmd := create.NewCommand(id, prId, "RS256")
	require.NoError(t, svc.Create(context.Background(), cmd))

	require.NotNil(t, savedKP)
	p := savedKP.ToPrimitives()
	assert.Equal(t, id, p.Id)
	assert.Equal(t, prId, p.ProtectedResourceId)
	assert.Equal(t, "RS256", p.Algorithm)
	assert.NotEmpty(t, p.PrivateKey)
	assert.NotEmpty(t, p.PublicKey)

	require.Len(t, publishedEvts, 1)
	evt, ok := publishedEvts[0].(securityEvents.KeyPairCreated)
	require.True(t, ok, "expected KeyPairCreated event")
	assert.Equal(t, "RS256", evt.Algorithm)
	assert.Equal(t, id, evt.AggregateID())
}

func TestCreateKeyPairService_ECAlgorithms(t *testing.T) {
	t.Parallel()

	for _, alg := range []string{"ES256", "ES384", "ES512"} {
		alg := alg
		t.Run("creates key pair with algorithm "+alg, func(t *testing.T) {
			t.Parallel()

			var savedKP *key_pairs.KeyPair
			repo := kpmocks.NewMockRepository(t)
			repo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, kp *key_pairs.KeyPair) {
				savedKP = kp
			}).Return(nil)
			bus := eventmocks.NewMockBus(t)
			bus.EXPECT().Publish(mock.Anything, mock.Anything)

			svc := create.NewService(repo, bus, keyPairInfra.NewPEMKeyGenerator())
			cmd := create.NewCommand(uuid.New().String(), uuid.New().String(), alg)
			require.NoError(t, svc.Create(context.Background(), cmd))

			require.NotNil(t, savedKP)
			p := savedKP.ToPrimitives()
			assert.Equal(t, alg, p.Algorithm)
			assert.NotEmpty(t, p.PrivateKey)
			assert.NotEmpty(t, p.PublicKey)
		})
	}
}

func TestCreateKeyPairService_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		id          string
		prId        string
		algorithm   string
		setupRepo   func(repo *kpmocks.MockRepository)
		wantErrCode app_errors.ErrorCode
	}{
		{
			name:        "rejects unsupported algorithm",
			id:          uuid.New().String(),
			prId:        uuid.New().String(),
			algorithm:   "HS256",
			wantErrCode: key_pairs.ErrKeyPairNotCreated,
		},
		{
			name:        "rejects invalid id",
			id:          "bad-id",
			prId:        uuid.New().String(),
			algorithm:   "RS256",
			wantErrCode: key_pairs.ErrKeyPairNotCreated,
		},
		{
			name:      "wraps repository save error",
			id:        uuid.New().String(),
			prId:      uuid.New().String(),
			algorithm: "RS256",
			setupRepo: func(repo *kpmocks.MockRepository) {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErrCode: key_pairs.ErrKeyPairNotSaved,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := kpmocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			bus := eventmocks.NewMockBus(t)
			svc := create.NewService(repo, bus, keyPairInfra.NewPEMKeyGenerator())

			cmd := create.NewCommand(tc.id, tc.prId, tc.algorithm)
			err := svc.Create(context.Background(), cmd)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
