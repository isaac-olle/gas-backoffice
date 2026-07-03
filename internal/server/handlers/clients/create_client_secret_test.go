package clients_test

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

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/clients"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
)

func TestCreateClientSecretHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		body       io.Reader
		clientIdFn func(context.Context) (ids.ClientId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "invalid_body_missing_fields",
			body:       bytes.NewBufferString(`{}`),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "missing_client_id_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"secret_id": "secret-123", "plaintext_secret": "s3cr3t"}),
			clientIdFn: handlerstest.MissingClientIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "bus_error_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"secret_id": "secret-123", "plaintext_secret": "s3cr3t"}),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			busErr:     errors.New("bus error"),
			wantErr:    true,
		},
		{
			name:       "success_returns_200",
			body:       handlerstest.JsonBody(map[string]string{"secret_id": "secret-123", "plaintext_secret": "s3cr3t"}),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := clients.NewCreateClientSecretHandler(bus, tt.clientIdFn)

			_, c := handlerstest.NewCtx(http.MethodPost, tt.body)
			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			// Success is a no-op; default code is 200
			assert.Equal(t, tt.wantStatus, c.Writer.Status())
		})
	}
}

func TestCreateClientSecretHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := clients.NewCreateClientSecretHandler(bus, handlerstest.ValidClientIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
