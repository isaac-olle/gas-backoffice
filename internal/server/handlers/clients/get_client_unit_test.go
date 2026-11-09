package clients_test

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	clientAggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
)

func TestGetClientHandler(t *testing.T) {
	t.Parallel()
	validClient, _ := clientAggregate.FromPrimitives(clientAggregate.ClientPrimitives{
		Id:             handlerstest.ClientID,
		Name:           "My App",
		SubscriptionId: handlerstest.SubscriptionID,
	})

	tests := []struct {
		name       string
		clientIdFn func(context.Context) (ids.ClientId, bool)
		busResult  any
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
			busErr:     errors.New("not found"),
			wantErr:    true,
		},
		{
			name:       "success_returns_200",
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			busResult:  validClient,
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(tt.busResult, tt.busErr)
			h := clients.NewGetClientHandler(bus, tt.clientIdFn, handlerstest.ValidAuthorExtractor())

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

func TestGetClientHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewQueryBus(nil, nil)
	h := clients.NewGetClientHandler(bus, handlerstest.ValidClientIdExtractor(), handlerstest.ValidAuthorExtractor())
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
