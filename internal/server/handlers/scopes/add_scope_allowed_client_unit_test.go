package scopes_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func TestAddScopeAllowedClientHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		scopeIdFn  func(context.Context) (ids.ScopeId, bool)
		clientIdFn func(context.Context) (ids.ClientId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "missing_scope_id_returns_error",
			scopeIdFn:  handlerstest.MissingScopeIdExtractor(),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "missing_client_id_returns_error",
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			clientIdFn: handlerstest.MissingClientIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "bus_error_returns_error",
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			busErr:     errors.New("bus error"),
			wantErr:    true,
		},
		{
			name:       "success_returns_204",
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			h := scopeHandlers.NewAddScopeAllowedClientHandler(bus, tt.scopeIdFn, tt.clientIdFn)

			_, c := handlerstest.NewCtx(http.MethodPost, nil)
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

func TestAddScopeAllowedClientHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandExecutor(nil)
	h := scopeHandlers.NewAddScopeAllowedClientHandler(bus, handlerstest.ValidScopeIdExtractor(), handlerstest.ValidClientIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
