package protocol_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/server/handlers/protocol"
	"github.com/isaac-olle/gas/internal/server/handlers/protocol/authentication_methods"
	"github.com/isaac-olle/gas/internal/server/middlewares"
)

type mockAuthMethod struct {
	canProcess bool
	canErr     error
	clientID   string
	authErr    error
}

func (m *mockAuthMethod) CanProcess(_ *gin.Context) (bool, error) {
	return m.canProcess, m.canErr
}

func (m *mockAuthMethod) Authenticate(_ *gin.Context) (string, error) {
	return m.clientID, m.authErr
}

func newAuthCtx() (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", nil)
	return w, c
}

func TestTokenAuthenticationHandler_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     *mockAuthMethod
		wantErr    bool
		wantClient string
	}{
		{
			name:    "no_method_matches",
			method:  &mockAuthMethod{canProcess: false},
			wantErr: true,
		},
		{
			name:    "can_process_returns_error",
			method:  &mockAuthMethod{canProcess: true, canErr: assert.AnError},
			wantErr: true,
		},
		{
			name:    "authenticate_returns_error",
			method:  &mockAuthMethod{canProcess: true, authErr: assert.AnError},
			wantErr: true,
		},
		{
			name:       "success",
			method:     &mockAuthMethod{canProcess: true, clientID: "client-123"},
			wantClient: "client-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := protocol.NewTokenAuthenticationHandler(
				authentication_methods.NewClientAuthMethodFactory(tt.method),
			)
			_, c := newAuthCtx()

			_, err := h.Execute(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			clientID, ok := protocol.AuthenticatedClientIDFromContext(c.Request.Context())
			assert.True(t, ok)
			assert.Equal(t, tt.wantClient, clientID)
		})
	}
}

func TestTokenAuthenticationHandler_Success_CallsNext(t *testing.T) {
	t.Parallel()

	nextCalled := false

	hm := middlewares.NewHandlerMiddleware()
	h := protocol.NewTokenAuthenticationHandler(
		authentication_methods.NewClientAuthMethodFactory(
			&mockAuthMethod{canProcess: true, clientID: "client-123"},
		),
	)

	r := gin.New()
	r.POST("/", hm.Run(h), func(c *gin.Context) {
		nextCalled = true
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	r.ServeHTTP(w, req)

	assert.True(t, nextCalled)
}

func TestTokenAuthenticationHandler_Default(t *testing.T) {
	t.Parallel()

	h := protocol.NewTokenAuthenticationHandler(
		authentication_methods.NewClientAuthMethodFactory(),
	)
	w, c := newAuthCtx()

	h.Default(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
