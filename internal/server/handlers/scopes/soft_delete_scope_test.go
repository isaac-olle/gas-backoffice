package scopes_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/server/handlers/handlerstest"
	scopeHandlers "github.com/isaac-olle/gas/internal/server/handlers/scopes"
)

func TestSoftDeleteScopeHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		scopeIdFn  func(context.Context) (ids.ScopeId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:      "missing_scope_id_returns_error",
			scopeIdFn: handlerstest.MissingScopeIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "bus_error_returns_error",
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			busErr:    errors.New("bus error"),
			wantErr:   true,
		},
		{
			name:       "success_returns_204",
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := scopeHandlers.NewSoftDeleteScopeHandler(bus, tt.scopeIdFn)

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

func TestSoftDeleteScopeHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := scopeHandlers.NewSoftDeleteScopeHandler(bus, handlerstest.ValidScopeIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodDelete, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
