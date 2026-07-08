package scopes_test

import (
	"context"
	"errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
)

func TestCreateScopesFromBodyHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		ownershipFn    func(context.Context) (users.Ownership, bool)
		prIdFn         func(context.Context) (ids2.ProtectedResourceId, bool)
		userIdFn       func(context.Context) (ids2.UserId, bool)
		scopeIdFn      func(context.Context) (ids2.ScopeId, bool)
		body           interface{}
		busErr         error
		wantErr        bool
		wantStatus     int
		wantDispatches int
	}{
		{
			name:           "empty_scopes_no_dispatches_success",
			body:           map[string][]string{"scopes": {}},
			ownershipFn:    handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:         handlerstest.ValidPrIdExtractor(),
			userIdFn:       handlerstest.ValidUserIdExtractor(),
			scopeIdFn:      handlerstest.NewScopeIdExtractor(),
			wantErr:        false,
			wantStatus:     http.StatusOK,
			wantDispatches: 0,
		},
		{
			name:           "single_scope_one_dispatch",
			body:           map[string][]string{"scopes": {"read"}},
			ownershipFn:    handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:         handlerstest.ValidPrIdExtractor(),
			userIdFn:       handlerstest.ValidUserIdExtractor(),
			scopeIdFn:      handlerstest.NewScopeIdExtractor(),
			wantErr:        false,
			wantStatus:     http.StatusOK,
			wantDispatches: 1,
		},
		{
			name:           "multi_scope_multiple_dispatches",
			body:           map[string][]string{"scopes": {"read", "write", "delete"}},
			ownershipFn:    handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:         handlerstest.ValidPrIdExtractor(),
			userIdFn:       handlerstest.ValidUserIdExtractor(),
			scopeIdFn:      handlerstest.NewScopeIdExtractor(),
			wantErr:        false,
			wantStatus:     http.StatusOK,
			wantDispatches: 3,
		},
		{
			name:        "bus_error_returns_error",
			body:        map[string][]string{"scopes": {"read"}},
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			scopeIdFn:   handlerstest.NewScopeIdExtractor(),
			busErr:      errors.New("bus error"),
			wantErr:     true,
		},
		{
			name:        "missing_user_ownership_when_scopes_non_empty",
			body:        map[string][]string{"scopes": {"read"}},
			ownershipFn: handlerstest.MissingUserOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			scopeIdFn:   handlerstest.NewScopeIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "missing_pr_id_when_scopes_non_empty",
			body:        map[string][]string{"scopes": {"read"}},
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:      handlerstest.MissingPrIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			scopeIdFn:   handlerstest.NewScopeIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "missing_user_id_when_scopes_non_empty",
			body:        map[string][]string{"scopes": {"read"}},
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			userIdFn:    handlerstest.MissingUserIdExtractor(),
			scopeIdFn:   handlerstest.NewScopeIdExtractor(),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := scopeHandlers.NewCreateScopesFromBodyHandler(bus, tt.ownershipFn, tt.prIdFn, tt.userIdFn, tt.scopeIdFn)

			_, c := handlerstest.NewCtx(http.MethodPost, handlerstest.JsonBody(tt.body))
			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, tt.wantStatus, c.Writer.Status())
			bus.AssertNumberOfCalls(t, "Dispatch", tt.wantDispatches)
		})
	}
}

func TestCreateScopesFromBodyHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := scopeHandlers.NewCreateScopesFromBodyHandler(bus, handlerstest.ValidUserOwnershipExtractor(), handlerstest.ValidPrIdExtractor(), handlerstest.ValidUserIdExtractor(), handlerstest.ValidScopeIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}

func TestCreateScopesFromBodyHandler_InvalidBody(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := scopeHandlers.NewCreateScopesFromBodyHandler(
		bus,
		handlerstest.ValidUserOwnershipExtractor(),
		handlerstest.ValidPrIdExtractor(),
		handlerstest.ValidUserIdExtractor(),
		handlerstest.NewScopeIdExtractor(),
	)
	_, c := handlerstest.NewCtx(http.MethodPost, strings.NewReader("not valid json"))
	_, err := h.Execute(c)
	assert.Error(t, err)
}
