package soft_delete_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/soft_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T) (*soft_delete.Service, *prmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	repo := prmocks.NewMockRepository(t)
	bus := eventmocks.NewMockBus(t)
	finder := domain_services.NewProtectedResourceFinder(repo)
	return soft_delete.NewService(finder, repo, bus), repo, bus
}

func TestSoftDelete(t *testing.T) {
	t.Parallel()

	t.Run("soft deletes the correct resource and publishes ProtectedResourceDeleted event", func(t *testing.T) {
		t.Parallel()
		svc, repo, bus := newService(t)
		userId := ids.NewUserId()
		pr := aPR(t, userId)
		prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
		require.NoError(t, err)

		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
		repo.EXPECT().SoftDelete(mock.Anything, prId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ProtectedResourceDeletedType
			})).
			Return()

		require.NoError(t, svc.SoftDelete(context.Background(), soft_delete.NewCommand(prId, ids.NewUserId(), userId)))
		assert.Equal(t, prId.String(), pr.ToPrimitives().Id)
	})

	t.Run("returns forbidden error when user is not the owner", func(t *testing.T) {
		t.Parallel()
		svc, repo, _ := newService(t)
		pr := aPR(t, ids.NewUserId())
		prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
		require.NoError(t, err)

		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)

		err = svc.SoftDelete(context.Background(), soft_delete.NewCommand(prId, ids.NewUserId(), ids.NewUserId()))
		assertErrCode(t, err, protected_resources.ErrInvalidOwnership)
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, ids.UserId)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns internal error when resource search fails",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, ids.UserId) {
				prId := ids.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, fmt.Errorf("db error"))
				return prId, ids.NewUserId()
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns not found error when resource is nil",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, ids.UserId) {
				prId := ids.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, nil)
				return prId, ids.NewUserId()
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns internal error when soft delete fails",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, ids.UserId) {
				userId := ids.NewUserId()
				pr := aPR(t, userId)
				prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
				require.NoError(t, err)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				repo.EXPECT().SoftDelete(mock.Anything, prId).Return(fmt.Errorf("db error"))
				return prId, userId
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
			prId, userId := tc.setupRepo(repo)
			err := soft_delete.NewService(finder, repo, bus).SoftDelete(context.Background(), soft_delete.NewCommand(prId, ids.NewUserId(), userId))
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---------------------------------------------------------------------------
// assertErrCode — shared helper for the soft_delete_test package.
// ---------------------------------------------------------------------------

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var e *app_errors.AppError
	require.ErrorAs(t, err, &e)
	assert.Equal(t, code, e.Code)
}
