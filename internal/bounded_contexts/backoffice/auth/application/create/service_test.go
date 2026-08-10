package create_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/create"
	authmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"

	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("saves session with correct data and publishes SessionCreated event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		userId := ids2.NewUserId()
		cmd := aCommand(t, withUserId(userId))

		var saved *auth.Session
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(s *auth.Session) bool {
				saved = s
				return true
			})).
			Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.SessionCreatedType
			}))

		require.NoError(t, create.NewService(repo, bus).Create(context.Background(), cmd))

		p := saved.ToPrimitives()
		assert.Equal(t, userId.String(), p.UserId)
	})

	t.Run("wraps repository save failure", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		repo.EXPECT().
			Save(mock.Anything, mock.AnythingOfType("*auth.Session")).
			Return(fmt.Errorf("db error"))

		err := create.NewService(repo, bus).Create(context.Background(), aCommand(t))
		assertErrCode(t, err, auth.ErrSessionNotSaved)
	})
}

// ── helpers ───────────────────────────────────────────────────────────────────

func newMocks(t *testing.T) (*authmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	return authmocks.NewMockRepository(t), eventmocks.NewMockBus(t)
}

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
