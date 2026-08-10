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

func validAssertionKeyIdExtractor() func(context.Context) (ids.ClientAssertionKeyId, bool) {
	keyId, _ := ids.NewClientAssertionKeyIdFromString(assertionKeyID)
	return func(_ context.Context) (ids.ClientAssertionKeyId, bool) { return keyId, true }
}

func missingAssertionKeyIdExtractor() func(context.Context) (ids.ClientAssertionKeyId, bool) {
	return func(_ context.Context) (ids.ClientAssertionKeyId, bool) { return ids.ClientAssertionKeyId{}, false }
}

func TestCreateClientAssertionKeyHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		body           io.Reader
		clientIdFn     func(context.Context) (ids.ClientId, bool)
		assertionKeyFn func(context.Context) (ids.ClientAssertionKeyId, bool)
		busErr         error
		wantErr        bool
		wantStatus     int
	}{
		{
			name:           "invalid_body_missing_fields",
			body:           bytes.NewBufferString(`{}`),
			clientIdFn:     handlerstest.ValidClientIdExtractor(),
			assertionKeyFn: validAssertionKeyIdExtractor(),
			wantErr:        true,
		},
		{
			name:           "missing_client_id_returns_error",
			body:           handlerstest.JsonBody(map[string]string{"public_key": "somekey"}),
			clientIdFn:     handlerstest.MissingClientIdExtractor(),
			assertionKeyFn: validAssertionKeyIdExtractor(),
			wantErr:        true,
		},
		{
			name:           "missing_assertion_key_id_returns_error",
			body:           handlerstest.JsonBody(map[string]string{"public_key": "somekey"}),
			clientIdFn:     handlerstest.ValidClientIdExtractor(),
			assertionKeyFn: missingAssertionKeyIdExtractor(),
			wantErr:        true,
		},
		{
			name:           "bus_error_returns_error",
			body:           handlerstest.JsonBody(map[string]string{"public_key": "somekey"}),
			clientIdFn:     handlerstest.ValidClientIdExtractor(),
			assertionKeyFn: validAssertionKeyIdExtractor(),
			busErr:         errors.New("bus error"),
			wantErr:        true,
		},
		{
			name:           "success_returns_201",
			body:           handlerstest.JsonBody(map[string]string{"public_key": "somekey"}),
			clientIdFn:     handlerstest.ValidClientIdExtractor(),
			assertionKeyFn: validAssertionKeyIdExtractor(),
			wantErr:        false,
			wantStatus:     http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			h := security.NewCreateClientAssertionKeyHandler(bus, tt.clientIdFn, tt.assertionKeyFn)

			_, c := handlerstest.NewCtx(http.MethodPut, tt.body)
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
	bus := handlerstest.NewCommandExecutor(nil)
	h := security.NewCreateClientAssertionKeyHandler(bus, handlerstest.ValidClientIdExtractor(), validAssertionKeyIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPut, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
