package protected_resources_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
)

func TestSoftDeleteProtectedResourceHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		prIdFn     func(context.Context) (ids.ProtectedResourceId, bool)
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
			busErr:  errors.New("bus error"),
			wantErr: true,
		},
		{
			name:       "success_returns_204",
			prIdFn:     handlerstest.ValidPrIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			h := prHandlers.NewSoftDeleteProtectedResourceHandler(bus, tt.prIdFn)

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

func TestSoftDeleteProtectedResourceHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandExecutor(nil)
	h := prHandlers.NewSoftDeleteProtectedResourceHandler(bus, handlerstest.ValidPrIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodDelete, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
