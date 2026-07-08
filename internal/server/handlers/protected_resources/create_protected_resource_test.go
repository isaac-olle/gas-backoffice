package protected_resources_test

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

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
)

func TestCreateProtectedResourceHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		body        io.Reader
		ownershipFn func(context.Context) (users.Ownership, bool)
		prIdFn      func(context.Context) (ids2.ProtectedResourceId, bool)
		userIdFn    func(context.Context) (ids2.UserId, bool)
		busErr      error
		wantErr     bool
		wantStatus  int
	}{
		{
			name:        "missing_user_ownership_returns_error",
			body:        handlerstest.JsonBody(map[string]string{"name": "My PR"}),
			ownershipFn: handlerstest.MissingUserOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "missing_pr_id_returns_error",
			body:        handlerstest.JsonBody(map[string]string{"name": "My PR"}),
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:      handlerstest.MissingPrIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "missing_user_id_returns_error",
			body:        handlerstest.JsonBody(map[string]string{"name": "My PR"}),
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			userIdFn:    handlerstest.MissingUserIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "invalid_body_returns_error",
			body:        bytes.NewBufferString(`{}`),
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			wantErr:     true,
		},
		{
			name:        "bus_error_returns_error",
			body:        handlerstest.JsonBody(map[string]string{"name": "My PR"}),
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			busErr:      errors.New("bus error"),
			wantErr:     true,
		},
		{
			name:        "success_returns_201",
			body:        handlerstest.JsonBody(map[string]string{"name": "My PR"}),
			ownershipFn: handlerstest.ValidUserOwnershipExtractor(),
			prIdFn:      handlerstest.ValidPrIdExtractor(),
			userIdFn:    handlerstest.ValidUserIdExtractor(),
			wantErr:     false,
			wantStatus:  http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := prHandlers.NewCreateProtectedResourceHandler(bus, tt.ownershipFn, tt.prIdFn, tt.userIdFn)

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

func TestCreateProtectedResourceHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := prHandlers.NewCreateProtectedResourceHandler(bus, handlerstest.ValidUserOwnershipExtractor(), handlerstest.ValidPrIdExtractor(), handlerstest.ValidUserIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
