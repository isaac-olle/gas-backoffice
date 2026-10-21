package scopes_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	scopesAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
)

func TestGetScopeHandler(t *testing.T) {
	t.Parallel()
	validScope, _ := scopesAggregate.FromPrimitives(scopesAggregate.ScopePrimitives{
		Id:                  handlerstest.ScopeID,
		Name:                "read",
		ProtectedResourceId: handlerstest.PrID,
		AdminId:             handlerstest.AdminID,
		UserId:              handlerstest.UserId,
	})

	tests := []struct {
		name       string
		scopeIdFn  func(context.Context) (ids.ScopeId, bool)
		prIdFn     func(context.Context) (ids.ProtectedResourceId, bool)
		userIdFn   func(context.Context) (ids.UserId, bool)
		busResult  any
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:      "missing_scope_id_returns_error",
			scopeIdFn: handlerstest.MissingScopeIdExtractor(),
			prIdFn:    handlerstest.ValidPrIdExtractor(),
			userIdFn:  handlerstest.ValidUserIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "missing_pr_id_returns_error",
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			prIdFn:    handlerstest.MissingPrIdExtractor(),
			userIdFn:  handlerstest.ValidUserIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "missing_user_id_returns_error",
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			prIdFn:    handlerstest.ValidPrIdExtractor(),
			userIdFn:  handlerstest.MissingUserIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "bus_error_returns_error",
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			prIdFn:    handlerstest.ValidPrIdExtractor(),
			userIdFn:  handlerstest.ValidUserIdExtractor(),
			busErr:    errors.New("not found"),
			wantErr:   true,
		},
		{
			name:       "success_returns_200",
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			prIdFn:     handlerstest.ValidPrIdExtractor(),
			userIdFn:   handlerstest.ValidUserIdExtractor(),
			busResult:  validScope,
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(tt.busResult, tt.busErr)
			h := scopeHandlers.NewGetScopeHandler(bus, tt.scopeIdFn, tt.prIdFn, tt.userIdFn, handlerstest.ValidAuthorExtractor())

			_, c := handlerstest.NewCtx(http.MethodGet, nil)
			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, tt.wantStatus, c.Writer.Status())
		})
	}
}

func TestGetScopeHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewQueryBus(nil, nil)
	h := scopeHandlers.NewGetScopeHandler(bus, handlerstest.ValidScopeIdExtractor(), handlerstest.ValidPrIdExtractor(), handlerstest.ValidUserIdExtractor(), handlerstest.ValidAuthorExtractor())
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
