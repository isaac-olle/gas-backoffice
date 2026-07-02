package hard_delete_test

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
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/application/hard_delete"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T) (*hard_delete.Service, *prmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	repo := prmocks.NewMockRepository(t)
	bus := eventmocks.NewMockBus(t)
	finder := domain_services.NewProtectedResourceFinder(repo)
	return hard_delete.NewService(finder, repo, bus), repo, bus
}

func TestHardDelete(t *testing.T) {
	t.Parallel()

	t.Run("hard deletes the correct resource and publishes ProtectedResourceDeleted event", func(t *testing.T) {
		t.Parallel()
		svc, repo, bus := newService(t)
		pr := aPR(t)
		primitives := pr.ToPrimitives()
		prId, err := ids.NewProtectedResourceIdFromString(primitives.Id)
		require.NoError(t, err)

		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
		repo.EXPECT().HardDelete(mock.Anything, prId).Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ProtectedResourceDeletedType
			})).
			Return()

		require.NoError(t, svc.HardDelete(context.Background(), hard_delete.NewCommand(prId)))
		assert.Equal(t, primitives.Id, prId.String())
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *prmocks.MockRepository) ids.ProtectedResourceId
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns internal error when resource search fails",
			setupRepo: func(repo *prmocks.MockRepository) ids.ProtectedResourceId {
				prId := ids.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, fmt.Errorf("db error"))
				return prId
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns not found error when resource is nil",
			setupRepo: func(repo *prmocks.MockRepository) ids.ProtectedResourceId {
				prId := ids.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, nil)
				return prId
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns internal error when hard delete fails",
			setupRepo: func(repo *prmocks.MockRepository) ids.ProtectedResourceId {
				pr := aPR(t)
				prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
				require.NoError(t, err)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				repo.EXPECT().HardDelete(mock.Anything, prId).Return(fmt.Errorf("db error"))
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
			prId := tc.setupRepo(repo)
			err := hard_delete.NewService(finder, repo, bus).HardDelete(context.Background(), hard_delete.NewCommand(prId))
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---------------------------------------------------------------------------
// assertErrCode — shared helper for the hard_delete_test package.
// ---------------------------------------------------------------------------

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var e *app_errors.AppError
	require.ErrorAs(t, err, &e)
	assert.Equal(t, code, e.Code)
}
