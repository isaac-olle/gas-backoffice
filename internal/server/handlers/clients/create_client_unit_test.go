package clients_test

import (
	"bytes"
	"context"
	"errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
)

func TestCreateClientHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		body             io.Reader
		clientIdFn       func(context.Context) (ids2.ClientId, bool)
		subscriptionIdFn func(context.Context) (ids2.SubscriptionId, bool)
		busErr           error
		wantErr          bool
		wantStatus       int
	}{
		{
			name:       "invalid_body_returns_error",
			body:       bytes.NewBufferString(`{}`),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			wantErr:    true,
		},
		{
			name:             "missing_client_id_returns_error",
			body:             handlerstest.JsonBody(validCreateClientBody()),
			clientIdFn:       handlerstest.MissingClientIdExtractor(),
			subscriptionIdFn: handlerstest.ValidSubscriptionIdExtractor(),
			wantErr:          true,
		},
		{
			name:             "missing_subscription_id_returns_error",
			body:             handlerstest.JsonBody(validCreateClientBody()),
			clientIdFn:       handlerstest.ValidClientIdExtractor(),
			subscriptionIdFn: handlerstest.MissingSubscriptionIdExtractor(),
			wantErr:          true,
		},
		{
			name:             "bus_error_returns_error",
			body:             handlerstest.JsonBody(validCreateClientBody()),
			clientIdFn:       handlerstest.ValidClientIdExtractor(),
			subscriptionIdFn: handlerstest.ValidSubscriptionIdExtractor(),
			busErr:           errors.New("bus error"),
			wantErr:          true,
		},
		{
			name:             "success_returns_201",
			body:             handlerstest.JsonBody(validCreateClientBody()),
			clientIdFn:       handlerstest.ValidClientIdExtractor(),
			subscriptionIdFn: handlerstest.ValidSubscriptionIdExtractor(),
			wantErr:          false,
			wantStatus:       http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			if tt.subscriptionIdFn == nil {
				tt.subscriptionIdFn = handlerstest.ValidSubscriptionIdExtractor()
			}
			h := clients.NewCreateClientHandler(bus, tt.clientIdFn, tt.subscriptionIdFn)

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

func TestCreateClientHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandExecutor(nil)
	h := clients.NewCreateClientHandler(bus, handlerstest.ValidClientIdExtractor(), handlerstest.ValidSubscriptionIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}

func validCreateClientBody() map[string]any {
	return map[string]any{
		"name":                "My Client",
		"allowed_grant_types": []string{"authorization_code"},
		"redirect_uris":       []string{"https://example.com/callback"},
		"authentication":      map[string]any{"client_secret": map[string]string{"id": "11111111-1111-4111-8111-111111111111", "secret": "secret"}},
	}
}
