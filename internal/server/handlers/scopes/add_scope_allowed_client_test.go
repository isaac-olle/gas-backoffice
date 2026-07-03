package scopes_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
)

func TestAddScopeAllowedClientHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		body       io.Reader
		scopeIdFn  func(context.Context) (ids.ScopeId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:      "invalid_body_missing_client_id",
			body:      bytes.NewBufferString(`{}`),
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "missing_scope_id_returns_error",
			body:      handlerstest.JsonBody(map[string]string{"client_id": handlerstest.ClientID}),
			scopeIdFn: handlerstest.MissingScopeIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "bus_error_returns_error",
			body:      handlerstest.JsonBody(map[string]string{"client_id": handlerstest.ClientID}),
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			busErr:    errors.New("bus error"),
			wantErr:   true,
		},
		{
			name:       "success_returns_204",
			body:       handlerstest.JsonBody(map[string]string{"client_id": handlerstest.ClientID}),
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := scopeHandlers.NewAddScopeAllowedClientHandler(bus, tt.scopeIdFn)

			_, c := handlerstest.NewCtx(http.MethodPost, tt.body)
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
	bus := handlerstest.NewCommandBus(nil)
	h := scopeHandlers.NewAddScopeAllowedClientHandler(bus, handlerstest.ValidScopeIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
