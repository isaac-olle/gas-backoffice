package tenants_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tenantAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/tenants"
)

func TestGetTenantHandler(t *testing.T) {
	t.Parallel()
	validTenant, _ := tenantAggregate.FromPrimitives(tenantAggregate.TenantPrimitives{
		Id:      handlerstest.TenantID,
		Name:    "Acme",
		AdminId: handlerstest.AdminID,
	})

	tests := []struct {
		name       string
		tenantIdFn func(context.Context) (ids.TenantId, bool)
		busResult  any
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "missing_tenant_id_returns_error",
			tenantIdFn: handlerstest.MissingTenantIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "bus_error_returns_error",
			tenantIdFn: handlerstest.ValidTenantIdExtractor(),
			busErr:     errors.New("not found"),
			wantErr:    true,
		},
		{
			name:       "success_returns_200",
			tenantIdFn: handlerstest.ValidTenantIdExtractor(),
			busResult:  validTenant,
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(tt.busResult, tt.busErr)
			h := tenants.NewGetTenantHandler(bus, tt.tenantIdFn)

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

func TestGetTenantHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewQueryBus(nil, nil)
	h := tenants.NewGetTenantHandler(bus, handlerstest.ValidTenantIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
