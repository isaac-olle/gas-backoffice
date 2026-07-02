package _event_handlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs/application/_event_handlers"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs/application/create"
	kpmocks "github.com/isaac-olle/gas/internal/bounded_contexts/security/key_pairs/mocks"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas/internal/shared/buses/event/mocks"
	backofficeevents "github.com/isaac-olle/gas/internal/shared/events/backoffice"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestCreateKeyPairOnProtectedResourceCreated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		repoSaveErr error
		wantErr     bool
		wantErrCode app_errors.ErrorCode
		wantSaved   bool
	}{
		{
			name:      "creates key pair when protected resource is created",
			wantSaved: true,
		},
		{
			name:        "propagates repository error",
			repoSaveErr: errors.New("db failure"),
			wantErr:     true,
			wantErrCode: key_pairs.ErrKeyPairNotSaved,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var savedKP *key_pairs.KeyPair
			repo := kpmocks.NewMockRepository(t)
			if tc.repoSaveErr != nil {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(tc.repoSaveErr)
			} else {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Run(func(_ context.Context, kp *key_pairs.KeyPair) {
					savedKP = kp
				}).Return(nil)
			}

			bus := eventmocks.NewMockBus(t)
			if tc.repoSaveErr == nil {
				bus.EXPECT().Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
					return len(evts) > 0
				}))
			}

			svc := create.NewService(repo, bus)
			handler := _event_handlers.NewCreateKeyPairOnProtectedResourceCreatedEventHandler(svc)

			evt := backofficeevents.NewProtectedResourceCreated(
				uuid.New().String(), "my-resource", uuid.New().String(),
			)

			err := handler.Handle(context.Background(), evt)

			if tc.wantErr {
				assertErrCode(t, err, tc.wantErrCode)
				return
			}
			require.NoError(t, err)
			if tc.wantSaved {
				assert.NotNil(t, savedKP)
			}
		})
	}
}
