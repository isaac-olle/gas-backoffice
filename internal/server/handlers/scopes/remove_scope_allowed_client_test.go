package scopes_test

import (
	"context"
	"errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
)

func TestRemoveScopeAllowedClientHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		scopeIdFn  func(context.Context) (ids2.ScopeId, bool)
		clientIdFn func(context.Context) (ids2.ClientId, bool)
		prIdFn     func(context.Context) (ids2.ProtectedResourceId, bool)
		userIdFn   func(context.Context) (ids2.UserId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "missing_scope_id_returns_error",
			scopeIdFn:  handlerstest.MissingScopeIdExtractor(),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			prIdFn:     handlerstest.ValidPrIdExtractor(),
			userIdFn:   handlerstest.ValidUserIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "missing_client_id_returns_error",
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			clientIdFn: handlerstest.MissingClientIdExtractor(),
			prIdFn:     handlerstest.ValidPrIdExtractor(),
			userIdFn:   handlerstest.ValidUserIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "missing_pr_id_returns_error",
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			prIdFn:     handlerstest.MissingPrIdExtractor(),
			userIdFn:   handlerstest.ValidUserIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "missing_user_id_returns_error",
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			prIdFn:     handlerstest.ValidPrIdExtractor(),
			userIdFn:   handlerstest.MissingUserIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "bus_error_returns_error",
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			prIdFn:     handlerstest.ValidPrIdExtractor(),
			userIdFn:   handlerstest.ValidUserIdExtractor(),
			busErr:     errors.New("bus error"),
			wantErr:    true,
		},
		{
			name:       "success_returns_204",
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			prIdFn:     handlerstest.ValidPrIdExtractor(),
			userIdFn:   handlerstest.ValidUserIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			h := scopeHandlers.NewRemoveScopeAllowedClientHandler(bus, tt.scopeIdFn, tt.clientIdFn, tt.prIdFn, tt.userIdFn, handlerstest.ValidAuthorExtractor())

			_, c := handlerstest.NewCtx(http.MethodDelete, nil)
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

func TestRemoveScopeAllowedClientHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandExecutor(nil)
	h := scopeHandlers.NewRemoveScopeAllowedClientHandler(bus, handlerstest.ValidScopeIdExtractor(), handlerstest.ValidClientIdExtractor(), handlerstest.ValidPrIdExtractor(), handlerstest.ValidUserIdExtractor(), handlerstest.ValidAuthorExtractor())
	_, c := handlerstest.NewCtx(http.MethodDelete, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
