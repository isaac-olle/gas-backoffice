package admins_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/admins"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
)

func TestTokenHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		body    io.Reader
		wantErr bool
	}{
		{
			name:    "invalid_body_missing_required_fields",
			body:    bytes.NewBufferString(`{}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := handlerstest.NewQueryBus(nil, nil)
			h := admins.NewTokenHandler(bus, "key-id", "https://issuer.example.com")

			_, c := handlerstest.NewCtx(http.MethodPost, tt.body)
			// TokenHandler uses ShouldBind (form), set form content type
			c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			c.Request.Body = io.NopCloser(tt.body)

			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
		})
	}
}

func TestTokenHandler_Success(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewQueryBus(nil, nil)
	h := admins.NewTokenHandler(bus, "key-id", "https://issuer.example.com")
	w, c := handlerstest.NewCtx(http.MethodPost, nil)
	result := map[string]any{"access_token": "tok", "token_type": "Bearer", "expires_in": 86400}
	require.NoError(t, h.Success(c, result))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTokenHandler_Default(t *testing.T) {
	t.Parallel()
	bus := handlerstest.NewQueryBus(nil, nil)
	h := admins.NewTokenHandler(bus, "key-id", "https://issuer.example.com")
	_, c := handlerstest.NewCtx(http.MethodPost, nil)
	h.Default(c)
	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
