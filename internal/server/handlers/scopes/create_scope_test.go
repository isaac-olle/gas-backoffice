package scopes_test

import (
	"bytes"
	"context"
	"errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	scopeHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/scopes"
)

func TestCreateScopeHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		body       io.Reader
		prIdFn     func(context.Context) (ids2.ProtectedResourceId, bool)
		scopeIdFn  func(context.Context) (ids2.ScopeId, bool)
		userIdFn   func(context.Context) (ids2.UserId, bool)
		queryErr   error
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:      "missing_pr_id_returns_error",
			body:      handlerstest.JsonBody(map[string]string{"name": "read"}),
			prIdFn:    handlerstest.MissingPrIdExtractor(),
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			userIdFn:  handlerstest.ValidUserIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "missing_scope_id_returns_error",
			body:      handlerstest.JsonBody(map[string]string{"name": "read"}),
			prIdFn:    handlerstest.ValidPrIdExtractor(),
			scopeIdFn: handlerstest.MissingScopeIdExtractor(),
			userIdFn:  handlerstest.ValidUserIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "missing_user_id_returns_error",
			body:      handlerstest.JsonBody(map[string]string{"name": "read"}),
			prIdFn:    handlerstest.ValidPrIdExtractor(),
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			userIdFn:  handlerstest.MissingUserIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "invalid_body_returns_error",
			body:      bytes.NewBufferString(`{}`),
			prIdFn:    handlerstest.ValidPrIdExtractor(),
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			userIdFn:  handlerstest.ValidUserIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "protected_resource_ownership_check_error_returns_error",
			body:      handlerstest.JsonBody(map[string]string{"name": "read"}),
			prIdFn:    handlerstest.ValidPrIdExtractor(),
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			userIdFn:  handlerstest.ValidUserIdExtractor(),
			queryErr:  errors.New("protected resource does not belong to the provided user"),
			wantErr:   true,
		},
		{
			name:      "bus_error_returns_error",
			body:      handlerstest.JsonBody(map[string]string{"name": "read"}),
			prIdFn:    handlerstest.ValidPrIdExtractor(),
			scopeIdFn: handlerstest.ValidScopeIdExtractor(),
			userIdFn:  handlerstest.ValidUserIdExtractor(),
			busErr:    errors.New("bus error"),
			wantErr:   true,
		},
		{
			name:       "success_returns_201",
			body:       handlerstest.JsonBody(map[string]string{"name": "read"}),
			prIdFn:     handlerstest.ValidPrIdExtractor(),
			scopeIdFn:  handlerstest.ValidScopeIdExtractor(),
			userIdFn:   handlerstest.ValidUserIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			qryBus := handlerstest.NewQueryBus(nil, tt.queryErr)
			h := scopeHandlers.NewCreateScopeHandler(bus, qryBus, tt.prIdFn, tt.scopeIdFn, tt.userIdFn, handlerstest.ValidAuthorExtractor())

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
	bus := handlerstest.NewCommandExecutor(nil)
	qryBus := handlerstest.NewQueryBus(nil, nil)
	h := scopeHandlers.NewCreateScopeHandler(bus, qryBus, handlerstest.ValidPrIdExtractor(), handlerstest.ValidScopeIdExtractor(), handlerstest.ValidUserIdExtractor(), handlerstest.ValidAuthorExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
