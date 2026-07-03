package tenants_test

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
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/tenants"
)

func TestChangeTenantNameHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		body       io.Reader
		tenantIdFn func(context.Context) (ids.TenantId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "missing_tenant_id_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"name": "New Name"}),
			tenantIdFn: handlerstest.MissingTenantIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "invalid_body_returns_error",
			body:       bytes.NewBufferString(`{}`),
			tenantIdFn: handlerstest.ValidTenantIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "bus_error_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"name": "New Name"}),
			tenantIdFn: handlerstest.ValidTenantIdExtractor(),
			busErr:     errors.New("bus error"),
			wantErr:    true,
		},
		{
			name:       "success_returns_204",
			body:       handlerstest.JsonBody(map[string]string{"name": "New Name"}),
			tenantIdFn: handlerstest.ValidTenantIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := tenants.NewChangeTenantNameHandler(bus, tt.tenantIdFn)

			_, c := handlerstest.NewCtx(http.MethodPatch, tt.body)
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

func TestChangeTenantNameHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := tenants.NewChangeTenantNameHandler(bus, handlerstest.ValidTenantIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPatch, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
