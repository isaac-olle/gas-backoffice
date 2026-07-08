package protected_resources_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	protectedResourcesAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
)

func TestGetProtectedResourceHandler(t *testing.T) {
	t.Parallel()
	validPR, _ := protectedResourcesAggregate.FromPrimitives(protectedResourcesAggregate.ProtectedResourcePrimitives{
		Id:      handlerstest.PrID,
		Name:    "My Resource",
		UserId:  handlerstest.UserId,
		AdminId: handlerstest.AdminID,
	})

	tests := []struct {
		name       string
		prIdFn     func(context.Context) (ids.ProtectedResourceId, bool)
		busResult  any
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:    "missing_pr_id_returns_error",
			prIdFn:  handlerstest.MissingPrIdExtractor(),
			wantErr: true,
		},
		{
			name:    "bus_error_returns_error",
			prIdFn:  handlerstest.ValidPrIdExtractor(),
			busErr:  errors.New("not found"),
			wantErr: true,
		},
		{
			name:       "success_returns_200",
			prIdFn:     handlerstest.ValidPrIdExtractor(),
			busResult:  validPR,
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(tt.busResult, tt.busErr)
			h := prHandlers.NewGetProtectedResourceHandler(bus, tt.prIdFn)

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

func TestGetProtectedResourceHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewQueryBus(nil, nil)
	h := prHandlers.NewGetProtectedResourceHandler(bus, handlerstest.ValidPrIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
