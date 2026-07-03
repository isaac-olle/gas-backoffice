package protected_resources_test

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
	prHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protected_resources"
)

func TestUpdateProtectedResourceNameHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		body       io.Reader
		prIdFn     func(context.Context) (ids.ProtectedResourceId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:    "invalid_body_returns_error",
			body:    bytes.NewBufferString(`{}`),
			prIdFn:  handlerstest.ValidPrIdExtractor(),
			wantErr: true,
		},
		{
			name:    "missing_pr_id_returns_error",
			body:    handlerstest.JsonBody(map[string]string{"name": "New Name"}),
			prIdFn:  handlerstest.MissingPrIdExtractor(),
			wantErr: true,
		},
		{
			name:    "bus_error_returns_error",
			body:    handlerstest.JsonBody(map[string]string{"name": "New Name"}),
			prIdFn:  handlerstest.ValidPrIdExtractor(),
			busErr:  errors.New("bus error"),
			wantErr: true,
		},
		{
			name:       "success_returns_204",
			body:       handlerstest.JsonBody(map[string]string{"name": "New Name"}),
			prIdFn:     handlerstest.ValidPrIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := prHandlers.NewUpdateProtectedResourceNameHandler(bus, tt.prIdFn)

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

func TestUpdateProtectedResourceNameHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := prHandlers.NewUpdateProtectedResourceNameHandler(bus, handlerstest.ValidPrIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPatch, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
