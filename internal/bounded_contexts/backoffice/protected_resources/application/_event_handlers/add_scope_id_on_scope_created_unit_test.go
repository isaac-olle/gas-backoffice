package _event_handlers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	event_handlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/_event_handlers"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	prmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

func TestAddScopeIdOnScopeCreated(t *testing.T) {
	t.Parallel()

	newHandler := func(repo *prmocks.MockRepository) *event_handlers.AddScopeIdOnScopeCreatedEventHandler {
		return event_handlers.NewAddScopeIdOnScopeCreatedEventHandler(domain_services.NewProtectedResourceFinder(repo), repo)
	}

	t.Run("adds scope id to protected resource scope ids and saves it", func(t *testing.T) {
		t.Parallel()
		repo := prmocks.NewMockRepository(t)
		pr := aPRForEventHandlers(t, ids2.NewSubscriptionId())
		prId, err := ids2.NewProtectedResourceIdFromString(pr.ToPrimitives().Id)
		require.NoError(t, err)
		scopeId := ids2.NewScopeId()

		evt := newScopeCreated(scopeId.String(), prId.String())

		repo.EXPECT().SearchById(mock.Anything, prId).Return(pr, nil)
		repo.EXPECT().Save(mock.Anything, pr).Return(nil)

		require.NoError(t, newHandler(repo).Handle(context.Background(), evt))
		assert.Contains(t, pr.ToPrimitives().ScopeIds, scopeId.String())
	})

	errCases := []struct {
		desc      string
		evt       backofficeevents.ScopeCreated
		setupRepo func(repo *prmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid protected resource id in event",
			evt:      newScopeCreated(ids2.NewScopeId().String(), "not-a-uuid"),
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc:     "rejects invalid scope id in event",
			evt:      newScopeCreated("not-a-uuid", ids2.NewProtectedResourceId().String()),
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "fails when protected resource does not exist",
			evt:  newScopeCreated(ids2.NewScopeId().String(), ids2.NewProtectedResourceId().String()),
			setupRepo: func(repo *prmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantCode: protected_resources.ErrProtectedResourceNotFound,
		},
		{
			desc: "maps concurrent modification to conflict",
			evt:  newScopeCreated(ids2.NewScopeId().String(), ids2.NewProtectedResourceId().String()),
			setupRepo: func(repo *prmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(aPRForEventHandlers(t, ids2.NewSubscriptionId()), nil)
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(persistance.ErrConcurrentModification)
			},
			wantCode: protected_resources.ErrProtectedResourceConcurrentModification,
		},
		{
			desc: "wraps repository save failure",
			evt:  newScopeCreated(ids2.NewScopeId().String(), ids2.NewProtectedResourceId().String()),
			setupRepo: func(repo *prmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(aPRForEventHandlers(t, ids2.NewSubscriptionId()), nil)
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(fmt.Errorf("db error"))
			},
			wantCode: protected_resources.ErrProtectedResourceNotSaved,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := prmocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			err := newHandler(repo).Handle(context.Background(), tc.evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}
