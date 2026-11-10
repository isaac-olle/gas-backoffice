package protocol_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	protocolHandler "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol"
)

type stubJwksService struct {
	keys []map[string]any
	err  error
}

func (s *stubJwksService) Get(context.Context) ([]map[string]any, error) { return s.keys, s.err }

func TestJwksHandler(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		keys     []map[string]any
		err      error
		wantErr  bool
		wantJSON string
	}{
		{
			name:    "service_error_returns_error",
			err:     errors.New("db"),
			wantErr: true,
		},
		{
			name:     "success_returns_key_set",
			keys:     []map[string]any{{"kid": "k1", "kty": "RSA"}},
			wantJSON: `{"keys":[{"kid":"k1","kty":"RSA"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := protocolHandler.NewJwksHandler(&stubJwksService{keys: tt.keys, err: tt.err})
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/oidc/jwks", nil)
			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, tt.wantJSON, w.Body.String())
		})
	}
}
