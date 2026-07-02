package tenants_test

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

	"github.com/isaac-olle/gas/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas/internal/server/handlers/tenants"
)

func TestCreateTenantHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		body       io.Reader
		adminIdFn  func(context.Context) (ids2.AdminId, bool)
		tenantIdFn func(context.Context) (ids2.TenantId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "invalid_body_returns_error",
			body:       bytes.NewBufferString(`{}`),
			adminIdFn:  handlerstest.ValidAdminIdExtractor(),
			tenantIdFn: handlerstest.ValidTenantIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "missing_admin_id_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"name": "My Tenant"}),
			adminIdFn:  handlerstest.MissingAdminIdExtractor(),
			tenantIdFn: handlerstest.ValidTenantIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "missing_tenant_id_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"name": "My Tenant"}),
			adminIdFn:  handlerstest.ValidAdminIdExtractor(),
			tenantIdFn: handlerstest.MissingTenantIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "bus_error_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"name": "My Tenant"}),
			adminIdFn:  handlerstest.ValidAdminIdExtractor(),
			tenantIdFn: handlerstest.ValidTenantIdExtractor(),
			busErr:     errors.New("bus error"),
			wantErr:    true,
		},
		{
			name:       "success_returns_201",
			body:       handlerstest.JsonBody(map[string]string{"name": "My Tenant"}),
			adminIdFn:  handlerstest.ValidAdminIdExtractor(),
			tenantIdFn: handlerstest.ValidTenantIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := tenants.NewCreateTenantHandler(bus, tt.adminIdFn, tt.tenantIdFn)

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

func TestCreateTenantHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := tenants.NewCreateTenantHandler(bus, handlerstest.ValidAdminIdExtractor(), handlerstest.ValidTenantIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
