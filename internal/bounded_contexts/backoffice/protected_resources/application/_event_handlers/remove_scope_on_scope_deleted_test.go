package _event_handlers_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event/mocks"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	eventhandlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/_event_handlers"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// makePRWithScope builds a ProtectedResource that already contains the given scope.
func makePRWithScope(t *testing.T, scopeId ids2.ScopeId) *protected_resources.ProtectedResource {
	t.Helper()
	pr := makePR(t)
	pr.AddScopeId(scopeId)
	return pr
}

func TestRemoveScopeOnScopeDeleted(t *testing.T) {
	t.Parallel()

	t.Run("saves PR with scope id removed and publishes empty events", func(t *testing.T) {
		t.Parallel()
		scopeId := ids2.NewScopeId()
		pr := makePRWithScope(t, scopeId)
		prId, err := ids2.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
		require.NoError(t, err)

		var savedPR *protected_resources.ProtectedResource
		repo := prmocks.NewMockRepository(t)
		bus := eventmocks.NewMockBus(t)
		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(p *protected_resources.ProtectedResource) bool {
				savedPR = p
				return true
			})).
			Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 0
			})).
			Return()

		handler := eventhandlers.NewRemoveScopeOnScopeDeletedEventHandler(repo, bus)
		evt := backofficeevents.NewScopeDeleted(scopeId.String(), prId.String())
		require.NoError(t, handler.Handle(context.Background(), evt))

		require.NotNil(t, savedPR)
		assert.NotContains(t, savedPR.ToPrimitives().ScopeIds, scopeId.String())
	})

	errCases := []struct {
		desc      string
		setupRepo func(repo *prmocks.MockRepository, scopeId ids2.ScopeId) string
		wantCode  app_errors.ErrorCode
	}{
		{
			desc: "returns bad request error for invalid protected resource id",
			setupRepo: func(_ *prmocks.MockRepository, _ ids2.ScopeId) string {
				return "bad-pr-id"
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns not found error when PR is not in repo",
			setupRepo: func(repo *prmocks.MockRepository, _ ids2.ScopeId) string {
				prId := ids2.NewProtectedResourceId()
				repo.EXPECT().SearchById(mock.Anything, prId).Return(nil, fmt.Errorf("not found"))
				return prId.String()
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "returns internal error when save fails",
			setupRepo: func(repo *prmocks.MockRepository, scopeId ids2.ScopeId) string {
				pr := makePRWithScope(t, scopeId)
				prId, err := ids2.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
				require.NoError(t, err)
				repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*protected_resources.ProtectedResource")).Return(fmt.Errorf("save error"))
				return prId.String()
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
			scopeId := ids2.NewScopeId()
			prIdStr := tc.setupRepo(repo, scopeId)
			handler := eventhandlers.NewRemoveScopeOnScopeDeletedEventHandler(repo, bus)
			evt := backofficeevents.NewScopeDeleted(scopeId.String(), prIdStr)
			err := handler.Handle(context.Background(), evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}
