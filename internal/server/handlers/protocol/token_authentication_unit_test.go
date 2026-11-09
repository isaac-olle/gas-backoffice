package protocol_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
)

type stubAuthClientService struct {
	result authenticate_client.Result
	err    error
}

func (s *stubAuthClientService) Authenticate(_ context.Context, _ authenticate_client.Credentials) (authenticate_client.Result, error) {
	return s.result, s.err
}

func TestTokenAuthentication_ConfidentialClient(t *testing.T) {
	t.Parallel()

	nextCalled := false
	hm := middlewares.NewHandlerMiddleware()
	svc := &stubAuthClientService{result: authenticate_client.Result{ClientID: "conf-client"}}

	r := gin.New()
	r.POST("/", protocol.TokenAuthentication(svc), hm.Run(&noopHandler{}), func(c *gin.Context) {
		nextCalled = true
		clientID, ok := protocol.AuthenticatedClientIDFromContext(c.Request.Context())
		assert.True(t, ok)
		assert.Equal(t, "conf-client", clientID)
		assert.False(t, protocol.IsPublicClientFromContext(c.Request.Context()))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	r.ServeHTTP(w, req)

	assert.True(t, nextCalled)
}

func TestTokenAuthentication_PublicClient(t *testing.T) {
	t.Parallel()

	nextCalled := false
	hm := middlewares.NewHandlerMiddleware()
	svc := &stubAuthClientService{result: authenticate_client.Result{ClientID: "pub-client", IsPublic: true}}

	r := gin.New()
	r.POST("/", protocol.TokenAuthentication(svc), hm.Run(&noopHandler{}), func(c *gin.Context) {
		nextCalled = true
		clientID, ok := protocol.AuthenticatedClientIDFromContext(c.Request.Context())
		assert.True(t, ok)
		assert.Equal(t, "pub-client", clientID)
		assert.True(t, protocol.IsPublicClientFromContext(c.Request.Context()))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	r.ServeHTTP(w, req)

	assert.True(t, nextCalled)
}

func TestTokenAuthentication_AuthFailure(t *testing.T) {
	t.Parallel()

	nextCalled := false
	svc := &stubAuthClientService{err: assert.AnError}

	r := gin.New()
	r.POST("/", protocol.TokenAuthentication(svc), func(c *gin.Context) {
		nextCalled = true
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
	_ = w
}

type noopHandler struct{}

func (h *noopHandler) Execute(_ *gin.Context) (any, error) { return nil, nil }
func (h *noopHandler) Success(c *gin.Context, _ any) error { c.Next(); return nil }
func (h *noopHandler) Default(_ *gin.Context)              {}
