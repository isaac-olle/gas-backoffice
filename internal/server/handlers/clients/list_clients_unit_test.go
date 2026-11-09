package clients_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aggregate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

func TestListClientsHandler(t *testing.T) {
	t.Parallel()
	item, err := aggregate.FromPrimitives(aggregate.ClientPrimitives{Id: handlerstest.ClientID, Name: "My App", SubscriptionId: handlerstest.SubscriptionID})
	require.NoError(t, err)

	tests := []struct {
		name      string
		hasAuthor bool
		busResult any
		busErr    error
		wantErr   bool
		wantLen   int
	}{
		{name: "missing author returns unauthorized", wantErr: true},
		{name: "bus_error_returns_error", hasAuthor: true, busErr: errors.New("db"), wantErr: true},
		{name: "empty_list_returns_200", hasAuthor: true, busResult: []*aggregate.Client{}, wantLen: 0},
		{name: "success_returns_200", hasAuthor: true, busResult: []*aggregate.Client{item}, wantLen: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(tt.busResult, tt.busErr)
			h := clients.NewListClientsHandler(bus)

			_, c := handlerstest.NewCtx(http.MethodGet, nil)
			if tt.hasAuthor {
				author, _ := handlerstest.ValidAuthorExtractor()(context.Background())
				c.Request = c.Request.WithContext(vo.WithAuthor(c.Request.Context(), author))
			}
			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			primitives, ok := result.([]aggregate.ClientPrimitives)
			require.True(t, ok)
			assert.Len(t, primitives, tt.wantLen)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusOK, c.Writer.Status())
		})
	}
}

func TestListClientsHandler_Default(t *testing.T) {
	t.Parallel()
	h := clients.NewListClientsHandler(handlerstest.NewQueryBus(nil, nil))
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}

func TestListClientsHandler_MissingAuthor(t *testing.T) {
	t.Parallel()
	h := clients.NewListClientsHandler(handlerstest.NewQueryBus(nil, nil))
	_, c := handlerstest.NewCtx(http.MethodGet, nil)
	_, err := h.Execute(c)
	assert.Error(t, err)
}
