package clients_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
)

func TestSoftDeleteClientHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		clientIdFn func(context.Context) (ids.ClientId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "missing_client_id_returns_error",
			clientIdFn: handlerstest.MissingClientIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "bus_error_returns_error",
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			busErr:     errors.New("bus error"),
			wantErr:    true,
		},
		{
			name:       "success_returns_204",
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			h := clients.NewSoftDeleteClientHandler(bus, tt.clientIdFn, handlerstest.ValidUserIdExtractor(), handlerstest.ValidAuthorExtractor())

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

func TestSoftDeleteClientHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandExecutor(nil)
	h := clients.NewSoftDeleteClientHandler(bus, handlerstest.ValidClientIdExtractor(), handlerstest.ValidUserIdExtractor(), handlerstest.ValidAuthorExtractor())
	_, c := handlerstest.NewCtx(http.MethodDelete, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
