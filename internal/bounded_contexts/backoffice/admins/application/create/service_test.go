package create_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/application/create"
	adminmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

func newService(t *testing.T) (*create.Service, *adminmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	repo := adminmocks.NewMockRepository(t)
	bus := eventmocks.NewMockBus(t)
	return create.NewService(repo, bus), repo, bus
}

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("saves admin as a system resource and publishes AdminCreated event", func(t *testing.T) {
		t.Parallel()
		svc, repo, bus := newService(t)

		var saved *admins.Admin
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(a *admins.Admin) bool {
				saved = a
				return true
			})).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.AdminCreatedType
			}))

		require.NoError(t, svc.Create(context.Background(), aCommand(t, withEmail("admin@example.com"))))
		assert.Equal(t, "admin@example.com", saved.ToPrimitives().Email)
		assert.True(t, saved.IsSystemResource())
	})

	errCases := []struct {
		desc      string
		cmd       *create.Command
		setupRepo func(repo *adminmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "empty name returns ErrAdminNotCreated",
			cmd:      aCommand(t, withName("")),
			wantCode: admins.ErrAdminNotCreated,
		},
		{
			desc:     "malformed email returns ErrAdminNotCreated",
			cmd:      aCommand(t, withEmail("not-an-email")),
			wantCode: admins.ErrAdminNotCreated,
		},
		{
			desc: "repository error returns ErrAdminNotSaved",
			cmd:  aCommand(t),
			setupRepo: func(repo *adminmocks.MockRepository) {
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*admins.Admin")).
					Return(fmt.Errorf("db error"))
			},
			wantCode: admins.ErrAdminNotSaved,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			svc, repo, bus := newService(t)
			_ = bus
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			err := svc.Create(context.Background(), tc.cmd)

			require.Error(t, err)
			var appErr *app_errors.AppError
			require.ErrorAs(t, err, &appErr)
			assert.Equal(t, tc.wantCode, appErr.Code)
		})
	}
}
