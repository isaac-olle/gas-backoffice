package update_name_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas/internal/shared/events/backoffice"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/application/update_name"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T) (*update_name.Service, *prmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	repo := prmocks.NewMockRepository(t)
	bus := eventmocks.NewMockBus(t)
	finder := domain_services.NewProtectedResourceFinder(repo)
	return update_name.NewService(finder, repo, bus), repo, bus
}

func TestUpdateName(t *testing.T) {
	t.Parallel()

	t.Run("saves PR with updated name and publishes ProtectedResourceNameUpdated event", func(t *testing.T) {
		t.Parallel()
		svc, repo, bus := newService(t)
		pr := aPR(t)
		prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
		require.NoError(t, err)

		var saved *protected_resources.ProtectedResource
		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(p *protected_resources.ProtectedResource) bool {
				saved = p
				return true
			})).
			Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ProtectedResourceNameUpdatedType
			})).
			Return()

		require.NoError(t, svc.UpdateName(context.Background(), update_name.NewCommand(prId, "New Name")))
		require.NotNil(t, saved)
		assert.Equal(t, "New Name", saved.ToPrimitives().Name)
	})

	errCases := []struct {
		desc      string
		newName   string
		setupRepo func(repo *prmocks.MockRepository) ids.ProtectedResourceId
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns bad request error for empty name",
			newName:  "",
			wantCode: protected_resources.ErrProtectedResourceNotUpdated,
		},
		{
			desc:     "returns bad request error for reserved name",
			newName:  "Gas-Backoffice",
			wantCode: protected_resources.ErrProtectedResourceNotUpdated,
		},
		{
			desc:    "returns not found error when resource is not in repo",
			newName: "New Name",
			setupRepo: func(repo *prmocks.MockRepository) ids.ProtectedResourceId {
				prId := ids.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, fmt.Errorf("not found"))
				return prId
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc:    "returns internal error when save fails",
			newName: "New Name",
			setupRepo: func(repo *prmocks.MockRepository) ids.ProtectedResourceId {
				pr := aPR(t)
				prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
				require.NoError(t, err)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*protected_resources.ProtectedResource")).Return(fmt.Errorf("db error"))
				return prId
			},
			wantCode: protected_resources.ErrProtectedResourceNotSaved,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := prmocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			finder := domain_services.NewProtectedResourceFinder(repo)
			var prId ids.ProtectedResourceId
			if tc.setupRepo != nil {
				prId = tc.setupRepo(repo)
			} else {
				prId = ids.NewProtectedResourceId()
			}
			err := update_name.NewService(finder, repo, bus).UpdateName(context.Background(), update_name.NewCommand(prId, tc.newName))
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---------------------------------------------------------------------------
// assertErrCode — shared helper for the update_name_test package.
// ---------------------------------------------------------------------------

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var e *app_errors.AppError
	require.ErrorAs(t, err, &e)
	assert.Equal(t, code, e.Code)
}
