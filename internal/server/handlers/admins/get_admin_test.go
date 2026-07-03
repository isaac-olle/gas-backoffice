package admins_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adminAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/admins"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
)

func TestGetAdminHandler(t *testing.T) {
	t.Parallel()
	validAdmin, _ := adminAggregate.FromPrimitives(adminAggregate.AdminPrimitives{
		Id:   handlerstest.AdminID,
		Name: "Test Admin",
	})

	tests := []struct {
		name       string
		adminIdFn  func(context.Context) (ids.AdminId, bool)
		busResult  any
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:      "missing_admin_id_returns_error",
			adminIdFn: handlerstest.MissingAdminIdExtractor(),
			wantErr:   true,
		},
		{
			name:      "bus_error_returns_error",
			adminIdFn: handlerstest.ValidAdminIdExtractor(),
			busErr:    errors.New("not found"),
			wantErr:   true,
		},
		{
			name:       "success_returns_200",
			adminIdFn:  handlerstest.ValidAdminIdExtractor(),
			busResult:  validAdmin,
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(tt.busResult, tt.busErr)
			h := admins.NewGetAdminHandler(bus, tt.adminIdFn)

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

func TestGetAdminHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewQueryBus(nil, nil)
	h := admins.NewGetAdminHandler(bus, handlerstest.ValidAdminIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
