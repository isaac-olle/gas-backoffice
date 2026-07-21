package _event_handlers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	event_handlers "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/_event_handlers"
	clientmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

func validOpenIDScopeIdFn() func() (ids2.ScopeId, error) {
	id := ids2.NewScopeId()
	return func() (ids2.ScopeId, error) { return id, nil }
}

func invalidOpenIDScopeIdFn() func() (ids2.ScopeId, error) {
	return func() (ids2.ScopeId, error) {
		return ids2.ScopeId{}, fmt.Errorf("invalid scope id")
	}
}

func TestAddOpenIDScopeOnClientCreated(t *testing.T) {
	t.Parallel()

	t.Run("adds openid scope to client and publishes event", func(t *testing.T) {
		t.Parallel()
		repo, bus := newMocks(t)
		c, clientId := aClient(t)
		scopeIdFn := validOpenIDScopeIdFn()
		scopeId, _ := scopeIdFn()
		evt := backofficeevents.NewClientCreated(clientId.String(), "test-client", ids2.NewAdminId().String(), ids2.NewUserId().String())

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

		handler := event_handlers.NewAddOpenIDScopeOnClientCreatedEventHandler(repo, bus, scopeIdFn)
		err := handler.Handle(context.Background(), evt)
		require.NoError(t, err)
		assert.Contains(t, saved.ToPrimitives().AllowedScopes, scopeId.String())
	})

	errCases := []struct {
		desc      string
		evt       backofficeevents.ClientCreated
		scopeIdFn func() (ids2.ScopeId, error)
		setupRepo func(repo *clientmocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:      "rejects invalid client id in event",
			evt:       backofficeevents.NewClientCreated("not-a-uuid", "test-client", ids2.NewAdminId().String(), ids2.NewUserId().String()),
			scopeIdFn: validOpenIDScopeIdFn(),
			wantCode:  clients.ErrClientNotFound,
		},
		{
			desc:      "returns error when repository search fails",
			evt:       backofficeevents.NewClientCreated(ids2.NewClientId().String(), "test-client", ids2.NewAdminId().String(), ids2.NewUserId().String()),
			scopeIdFn: validOpenIDScopeIdFn(),
			setupRepo: func(repo *clientmocks.MockRepository) {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))
			},
			wantCode: clients.ErrClientNotSaved,
		},
		{
			desc:      "returns error when openid scope id is invalid",
			evt:       backofficeevents.NewClientCreated(ids2.NewClientId().String(), "test-client", ids2.NewAdminId().String(), ids2.NewUserId().String()),
			scopeIdFn: invalidOpenIDScopeIdFn(),
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.NewClient(ids2.NewClientId(), "valid-client", clients.NewOwnership(ids2.NewAdminId(), ids2.NewUserId()), nil, nil)
				c.PullDomainEvents()
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(c, nil)
			},
			wantCode: "OPENID_SCOPE_ID_INVALID",
		},
		{
			desc:      "wraps repository save failure",
			evt:       backofficeevents.NewClientCreated(ids2.NewClientId().String(), "test-client", ids2.NewAdminId().String(), ids2.NewUserId().String()),
			scopeIdFn: validOpenIDScopeIdFn(),
			setupRepo: func(repo *clientmocks.MockRepository) {
				c, _ := clients.NewClient(ids2.NewClientId(), "valid-client", clients.NewOwnership(ids2.NewAdminId(), ids2.NewUserId()), nil, nil)
				c.PullDomainEvents()
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
			handler := event_handlers.NewAddOpenIDScopeOnClientCreatedEventHandler(repo, bus, tc.scopeIdFn)
			err := handler.Handle(context.Background(), tc.evt)
			assertErrCode(t, err, tc.wantCode)
		})
	}
}
