package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

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

func newCmd(t *testing.T, cakIdStr, clientIdStr, algorithm, publicKey string) *create.Command {
	t.Helper()
	cakId, err := ids2.NewClientAssertionKeyIdFromString(cakIdStr)
	require.NoError(t, err)
	clientId, err := ids2.NewClientIdFromString(clientIdStr)
	require.NoError(t, err)
	return create.NewCommand(cakId, clientId, algorithm, publicKey)
}

func TestCreateClientAssertionKeyService_HappyPath(t *testing.T) {
	t.Parallel()

	cakId := uuid.New().String()
	clientId := uuid.New().String()

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

	svc := create.NewService(repo, bus)

	cmd := newCmd(t, cakId, clientId, "RS256", validPubKeyHex)
	require.NoError(t, svc.Create(context.Background(), cmd))

	require.NotNil(t, savedCAK)
	p := savedCAK.ToPrimitives()
	assert.Equal(t, cakId, p.Id)
	assert.Equal(t, clientId, p.ClientId)
	assert.NotEmpty(t, p.PublicKey, "saved key must have a non-empty public key")

	require.Len(t, publishedEvts, 1)
	evt, ok := publishedEvts[0].(securityEvents.ClientAssertionKeyCreated)
	require.True(t, ok, "expected ClientAssertionKeyCreated event")
	assert.Equal(t, clientId, evt.ClientId)
	assert.Equal(t, cakId, evt.AggregateID())
}

func TestCreateClientAssertionKeyService_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cakId       string
		clientId    string
		publicKey   string
		setupRepo   func(repo *cakmocks.MockClientAssertionKeyRepository)
		wantErrCode app_errors.ErrorCode
	}{
		{
			name:        "rejects empty public key",
			cakId:       uuid.New().String(),
			clientId:    uuid.New().String(),
			publicKey:   "",
			wantErrCode: client_assertion_keys.ErrClientAssertionKeyNotCreated,
		},
		{
			name:      "wraps repository save error",
			cakId:     uuid.New().String(),
			clientId:  uuid.New().String(),
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
			svc := create.NewService(repo, bus)

			cmd := newCmd(t, tc.cakId, tc.clientId, "RS256", tc.publicKey)
			err := svc.Create(context.Background(), cmd)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
