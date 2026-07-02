package scopes_test

import (
	"bytes"
	"context"
	"errors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas/internal/server/handlers/handlerstest"
	scopeHandlers "github.com/isaac-olle/gas/internal/server/handlers/scopes"
)

func TestCreateScopeHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		body        io.Reader
		ownershipFn func(context.Context) (protected_resources.Ownership, bool)
		prIdFn      func(context.Context) (ids2.ProtectedResourceId, bool)
		scopeIdFn   func(context.Context) (ids2.ScopeId, bool)
		busErr      error
		wantErr     bool
		wantStatus  int
	}{
		{
			name:        "missing_pr_ownership_returns_error",
			body:        handlerstest.JsonBody(map[string]string{"name": "read"}),
			ownershipFn: handlerstest.MissingPrOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			scopeIdFn:   handlerstest.ValidScopeIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "missing_pr_id_returns_error",
			body:        handlerstest.JsonBody(map[string]string{"name": "read"}),
			ownershipFn: handlerstest.ValidPrOwnershipExtractor(),
			prIdFn:      handlerstest.MissingPrIdExtractor(),
			scopeIdFn:   handlerstest.ValidScopeIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "missing_scope_id_returns_error",
			body:        handlerstest.JsonBody(map[string]string{"name": "read"}),
			ownershipFn: handlerstest.ValidPrOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			scopeIdFn:   handlerstest.MissingScopeIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "invalid_body_returns_error",
			body:        bytes.NewBufferString(`{}`),
			ownershipFn: handlerstest.ValidPrOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			scopeIdFn:   handlerstest.ValidScopeIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "bus_error_returns_error",
			body:        handlerstest.JsonBody(map[string]string{"name": "read"}),
			ownershipFn: handlerstest.ValidPrOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			scopeIdFn:   handlerstest.ValidScopeIdExtractor(),
			busErr:      errors.New("bus error"),
			wantErr:     true,
		},
		{
			name:        "success_returns_201",
			body:        handlerstest.JsonBody(map[string]string{"name": "read"}),
			ownershipFn: handlerstest.ValidPrOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			scopeIdFn:   handlerstest.ValidScopeIdExtractor(),
			wantErr:     false,
			wantStatus:  http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := scopeHandlers.NewCreateScopeHandler(bus, tt.ownershipFn, tt.prIdFn, tt.scopeIdFn)

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

func TestCreateScopeHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := scopeHandlers.NewCreateScopeHandler(bus, handlerstest.ValidPrOwnershipExtractor(), handlerstest.ValidPrIdExtractor(), handlerstest.ValidScopeIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
