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

const secretID = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

func validSecretIdExtractor() func(context.Context) (ids.ClientSecretId, bool) {
	secretId, _ := ids.NewClientSecretIdFromString(secretID)
	return func(_ context.Context) (ids.ClientSecretId, bool) { return secretId, true }
}

func missingSecretIdExtractor() func(context.Context) (ids.ClientSecretId, bool) {
	return func(_ context.Context) (ids.ClientSecretId, bool) { return ids.ClientSecretId{}, false }
}

func TestCreateClientSecretHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		body       io.Reader
		clientIdFn func(context.Context) (ids.ClientId, bool)
		secretIdFn func(context.Context) (ids.ClientSecretId, bool)
		busErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "invalid_body_missing_fields",
			body:       bytes.NewBufferString(`{}`),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			secretIdFn: validSecretIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "missing_client_id_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"plaintext_secret": "s3cr3t"}),
			clientIdFn: handlerstest.MissingClientIdExtractor(),
			secretIdFn: validSecretIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "missing_secret_id_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"plaintext_secret": "s3cr3t"}),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			secretIdFn: missingSecretIdExtractor(),
			wantErr:    true,
		},
		{
			name:       "bus_error_returns_error",
			body:       handlerstest.JsonBody(map[string]string{"plaintext_secret": "s3cr3t"}),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			secretIdFn: validSecretIdExtractor(),
			busErr:     errors.New("bus error"),
			wantErr:    true,
		},
		{
			name:       "success_returns_201",
			body:       handlerstest.JsonBody(map[string]string{"plaintext_secret": "s3cr3t"}),
			clientIdFn: handlerstest.ValidClientIdExtractor(),
			secretIdFn: validSecretIdExtractor(),
			wantErr:    false,
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewCommandExecutor(tt.busErr)
			h := security.NewCreateClientSecretHandler(bus, tt.clientIdFn, tt.secretIdFn)

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

func TestCreateClientSecretHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewCommandExecutor(nil)
	h := security.NewCreateClientSecretHandler(bus, handlerstest.ValidClientIdExtractor(), validSecretIdExtractor())
	_, c := handlerstest.NewCtx(http.MethodPut, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
