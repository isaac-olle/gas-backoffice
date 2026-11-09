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

	t.Run("service_error_returns_error", func(t *testing.T) {
		t.Parallel()
		h := protocolHandler.NewJwksHandler(&stubJwksService{err: errors.New("db")})
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/oidc/jwks", nil)
		_, err := h.Execute(c)
		assert.Error(t, err)
	})

	t.Run("success_returns_key_set", func(t *testing.T) {
		t.Parallel()
		h := protocolHandler.NewJwksHandler(&stubJwksService{keys: []map[string]any{{"kid": "k1", "kty": "RSA"}}})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/oidc/jwks", nil)
		result, err := h.Execute(c)
		require.NoError(t, err)
		require.NoError(t, h.Success(c, result))
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"keys":[{"kid":"k1","kty":"RSA"}]}`, w.Body.String())
	})
}
