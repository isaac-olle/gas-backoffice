package _event_handlers_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas/internal/shared/events/backoffice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients"
	event_handlers "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/application/_event_handlers"
	clientmocks "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/mocks"
)

func TestRemoveScopeOnScopeClientRemoved(t *testing.T) {
	t.Parallel()

	t.Run("removes scope id from client and publishes event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		scopeId := ids2.NewScopeId()
		c, clientId := aClientWithScope(t, scopeId)
		evt := backofficeevents.NewScopeClientRemoved(scopeId.String(), clientId.String())

		var saved *clients.Client
		repo.EXPECT().SearchById(mock.Anything, clientId).Return(c, nil)
		repo.EXPECT().
			Save(mock.Anything, mock.MatchedBy(func(cl *clients.Client) bool {
				saved = cl
				return true
			})).
			Return(nil)
		bus.EXPECT().
			Publish(mock.Anything, mock.MatchedBy(func(evts []event.Event) bool {
				return len(evts) == 0
			}))

		handler := event_handlers.NewRemoveScopeOnScopeClientRemovedEventHandler(repo, bus)
		err := handler.Handle(context.Background(), evt)
		require.NoError(t, err)
		assert.NotContains(t, saved.ToPrimitives().AllowedScopes, scopeId.String())
	})

	errCases := []struct {
		desc      string
		evt       backofficeevents.ScopeClientRemoved
		setupRepo func(repo *clientmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "rejects invalid client id in event",
			evt:      backofficeevents.NewScopeClientRemoved(ids2.NewScopeId().String(), "not-a-uuid"),
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "returns not found when repository errors",
			evt:  backofficeevents.NewScopeClientRemoved(ids2.NewScopeId().String(), ids2.NewClientId().String()),
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("not found"))
			},
			wantCode: clients.ErrClientNotFound,
		},
		{
			desc: "wraps repository save failure",
			evt:  backofficeevents.NewScopeClientRemoved(ids2.NewScopeId().String(), ids2.NewClientId().String()),
			setupRepo: func(repo *clientmocks.MockRepository) {
				scopeId := ids2.NewScopeId()
				prims := clients.ClientPrimitives{
					Id:            ids2.NewClientId().String(),
					Name:          "valid-client",
					AdminId:       ids2.NewAdminId().String(),
					TenantId:      ids2.NewTenantId().String(),
					AllowedScopes: []string{scopeId.String()},
				}
				c, _ := clients.FromPrimitives(prims)
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
				repo.EXPECT().Save(mock.Anything, mock.AnythingOfType("*clients.Client")).Return(fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotSaved,
		},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo, bus := newMocks(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}
			handler := event_handlers.NewRemoveScopeOnScopeClientRemovedEventHandler(repo, bus)
			err := handler.Handle(context.Background(), tc.evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}
