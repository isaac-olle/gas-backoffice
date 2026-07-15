package security_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/security"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const assertionKeyID = "f1b2c3d4-e5f6-7890-abcd-ef1234567890"

func TestCreateClientAssertionKeyHandler(t *testing.T) {
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
			name:       "invalid_assertion_key_id_not_uuid",
			body:       handlerstest.JsonBody(map[string]string{"assertion_key_id": "not-a-uuid", "public_key": "somekey"}),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "missing_client_id_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"assertion_key_id": assertionKeyID, "public_key": "somekey"}),
			clientIdFn: handlerstest.MissingClientIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "bus_error_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"assertion_key_id": assertionKeyID, "public_key": "somekey"}),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			busErr:     errors.New("bus error"),
			wantErr:    true,
		},
		{
			name:       "success_is_noop",
			body:       handlerstest.JsonBody(map[string]string{"assertion_key_id": assertionKeyID, "public_key": "somekey"}),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandBus(tt.busErr)
			h := security.NewCreateClientAssertionKeyHandler(bus, tt.clientIdFn)

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

func TestCreateClientAssertionKeyHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandBus(nil)
	h := security.NewCreateClientAssertionKeyHandler(bus, handlerstest.ValidClientIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
