package create_test

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
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("saves PR with correct data and publishes ProtectedResourceCreated event", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		duplicateNameChecker := domain_services.NewDuplicateNameChecker(repo)
		cmd := aCommand(t, withCmdName("My-Resource"))

		var saved *protected_resources.ProtectedResource
		repo.EXPECT().
			Matching(mock.Anything, mock.Anything).
			Return(nil, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(pr *protected_resources.ProtectedResource) bool {
				saved = pr
				return true
			})).
			Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 1 && evts[0].Type() == backofficeEvents.ProtectedResourceCreatedType
			})).
			Return()

		require.NoError(t, create.NewService(duplicateNameChecker, repo, bus).Create(context.Background(), cmd))
		require.NotNil(t, saved)
		assert.Equal(t, "My-Resource", saved.ToPrimitives().Name)
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *create.Command
		setupRepo func(repo *prmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns bad request error for empty name",
			cmd:      func(t *testing.T) *create.Command { return aCommand(t, withCmdName("")) },
			wantCode: protected_resources.ErrProtectedResourceNotCreated,
		},
		{
			desc: "returns conflict error when name already exists",
			cmd:  func(t *testing.T) *create.Command { return aCommand(t) },
			setupRepo: func(repo *prmocks.MockRepository) {
				existing, _ := protected_resources.NewProtectedResource(
					ids.NewProtectedResourceId(), "My Resource", protected_resources.NewOwnership(ids.NewUserId(), ids.NewUserId()), false,
				)
				repo.EXPECT().
					Matching(mock.Anything, mock.Anything).
					Return([]*protected_resources.ProtectedResource{existing}, nil)
			},
			wantCode: protected_resources.ErrProtectedResourceDuplicateName,
		},
		{
			desc: "returns internal error when save fails",
			cmd:  func(t *testing.T) *create.Command { return aCommand(t) },
			setupRepo: func(repo *prmocks.MockRepository) {
				repo.EXPECT().
					Matching(mock.Anything, mock.Anything).
					Return(nil, nil)
				repo.EXPECT().
					Save(mock.Anything, mock.AnythingOfType("*protected_resources.ProtectedResource")).
					Return(fmt.Errorf("db error"))
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
			duplicateNameChecker := domain_services.NewDuplicateNameChecker(repo)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			err := create.NewService(duplicateNameChecker, repo, bus).Create(context.Background(), tc.cmd(t))
			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// ---------------------------------------------------------------------------
// assertErrCode — shared helper for the create_test package.
// ---------------------------------------------------------------------------

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var e *app_errors.AppError
	require.ErrorAs(t, err, &e)
	assert.Equal(t, code, e.Code)
}
