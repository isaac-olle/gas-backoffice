package update_name_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/update_name"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T) (*update_name.Service, *prmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	repo := prmocks.NewMockRepository(t)
	bus := eventmocks.NewMockBus(t)
	finder := domain_services.NewProtectedResourceFinder(repo)
	duplicateNameChecker := domain_services.NewDuplicateNameChecker(repo)
	return update_name.NewService(finder, duplicateNameChecker, repo, bus), repo, bus
}

func TestUpdateName(t *testing.T) {
	t.Parallel()

	t.Run("saves PR with updated name and publishes ProtectedResourceNameUpdated event", func(t *testing.T) {
		t.Parallel()
		svc, repo, bus := newService(t)
		userId := ids.NewUserId()
		pr := aPR(t, userId)
		prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
		require.NoError(t, err)

		var saved *protected_resources.ProtectedResource
		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
		repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, nil)
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

		require.NoError(t, svc.UpdateName(context.Background(), update_name.NewCommand(prId, "New Name", ids.NewUserId().String(), userId.String())))
		require.NotNil(t, saved)
		assert.Equal(t, "New Name", saved.ToPrimitives().Name)
	})

	t.Run("returns forbidden error when user is not the owner", func(t *testing.T) {
		t.Parallel()
		svc, repo, _ := newService(t)
		pr := aPR(t, ids.NewUserId())
		prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
		require.NoError(t, err)

		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)

		err = svc.UpdateName(context.Background(), update_name.NewCommand(prId, "New Name", ids.NewUserId().String(), ids.NewUserId().String()))
		assertErrCode(t, err, protected_resources.ErrInvalidOwnership)
	})

	errCases := []struct {
		desc      string
		newName   string
		setupRepo func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, string)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns bad request error for empty name",
			newName:  "",
			wantCode: protected_resources.ErrProtectedResourceNotUpdated,
		},
		{
			desc:    "returns not found error when resource is not in repo",
			newName: "New Name",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, string) {
				prId := ids.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, fmt.Errorf("not found"))
				return prId, ids.NewUserId().String()
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc:    "returns conflict error when name already exists",
			newName: "New Name",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, string) {
				userId := ids.NewUserId()
				pr := aPR(t, userId)
				prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
				require.NoError(t, err)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return([]*protected_resources.ProtectedResource{pr}, nil)
				return prId, userId.String()
			},
			wantCode: protected_resources.ErrProtectedResourceDuplicateName,
		},
		{
			desc:    "returns internal error when save fails",
			newName: "New Name",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, string) {
				userId := ids.NewUserId()
				pr := aPR(t, userId)
				prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
				require.NoError(t, err)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*protected_resources.ProtectedResource")).Return(fmt.Errorf("db error"))
				return prId, userId.String()
			},
			wantCode: protected_resources.ErrProtectedResourceNotSaved,
		},
		{
			desc:    "returns conflict when repository reports concurrent modification",
			newName: "New Name",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, string) {
				userId := ids.NewUserId()
				pr := aPR(t, userId)
				prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
				require.NoError(t, err)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				repo.EXPECT().Matching(mock.Anything, mock.Anything).Return(nil, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*protected_resources.ProtectedResource")).Return(persistance.ErrConcurrentModification)
				return prId, userId.String()
			},
			wantCode: protected_resources.ErrProtectedResourceConcurrentModification,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := prmocks.NewMockRepository(t)
			bus := eventmocks.NewMockBus(t)
			finder := domain_services.NewProtectedResourceFinder(repo)
			duplicateNameChecker := domain_services.NewDuplicateNameChecker(repo)
			var prId ids.ProtectedResourceId
			var userId string
			if tc.setupRepo != nil {
				prId, userId = tc.setupRepo(repo)
			} else {
				prId = ids.NewProtectedResourceId()
				userId = ids.NewUserId().String()
			}
			err := update_name.NewService(finder, duplicateNameChecker, repo, bus).UpdateName(context.Background(), update_name.NewCommand(prId, tc.newName, ids.NewUserId().String(), userId))
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
