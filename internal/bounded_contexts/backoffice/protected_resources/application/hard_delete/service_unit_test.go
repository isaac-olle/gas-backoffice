package hard_delete_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	eventmocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	backofficeEvents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/hard_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T) (*hard_delete.Service, *prmocks.MockRepository, *eventmocks.MockBus) {
	t.Helper()
	repo := prmocks.NewMockRepository(t)
	bus := eventmocks.NewMockBus(t)
	finder := domain_services.NewProtectedResourceFinder(repo)
	return hard_delete.NewService(finder, repo, bus, sharedDomainServices.NewSubscriptionAuthorizer()), repo, bus
}

func authorFor(t *testing.T, subscriptionId ids.SubscriptionId) *vo.UserAuthor {
	t.Helper()
	author, err := vo.NewUserAuthorWithMemberships(ids.NewUserId().String(), ids.NewUserId().String(), vo.SubscriptionUserRole,
		[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
	require.NoError(t, err)
	return author
}

func TestHardDelete(t *testing.T) {
	t.Parallel()

	t.Run("hard deletes the correct resource and publishes ProtectedResourceDeleted event", func(t *testing.T) {
		t.Parallel()
		svc, repo, bus := newService(t)
		subscriptionId := ids.NewSubscriptionId()
		pr := aPROwnedBy(t, subscriptionId)
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

		require.NoError(t, svc.HardDelete(context.Background(), hard_delete.NewCommand(prId, authorFor(t, subscriptionId))))
		assert.Equal(t, primitives.Id, prId.String())
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, *vo.UserAuthor)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns internal error when resource search fails",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, *vo.UserAuthor) {
				prId := ids.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, fmt.Errorf("db error"))
				return prId, authorFor(t, ids.NewSubscriptionId())
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns not found error when resource is nil",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, *vo.UserAuthor) {
				prId := ids.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, nil)
				return prId, authorFor(t, ids.NewSubscriptionId())
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns forbidden when caller does not have access to the resource's subscription",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, *vo.UserAuthor) {
				pr := aPR(t)
				prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
				require.NoError(t, err)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				return prId, authorFor(t, ids.NewSubscriptionId())
			},
			wantCode: protected_resources.ErrInvalidOwnership,
		},
		{
			desc: "returns forbidden when protected resource is system resource",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, *vo.UserAuthor) {
				prId := ids.NewProtectedResourceId()
				subscriptionId := ids.NewSubscriptionId()
				pr, err := protected_resources.FromPrimitives(protected_resources.ProtectedResourcePrimitives{
					Id:             prId.String(),
					Name:           "system-resource",
					SubscriptionId: subscriptionId.String(),
					IsSystem:       true,
				})
				require.NoError(t, err)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				return prId, authorFor(t, subscriptionId)
			},
			wantCode: protected_resources.ErrSystemResource,
		},
		{
			desc: "returns internal error when hard delete fails",
			setupRepo: func(repo *prmocks.MockRepository) (ids.ProtectedResourceId, *vo.UserAuthor) {
				subscriptionId := ids.NewSubscriptionId()
				pr := aPROwnedBy(t, subscriptionId)
				prId, err := ids.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
				require.NoError(t, err)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				repo.EXPECT().HardDelete(mock.Anything, prId).Return(fmt.Errorf("db error"))
				return prId, authorFor(t, subscriptionId)
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
			prId, author := tc.setupRepo(repo)
			err := hard_delete.NewService(finder, repo, bus, sharedDomainServices.NewSubscriptionAuthorizer()).HardDelete(context.Background(), hard_delete.NewCommand(prId, author))
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
