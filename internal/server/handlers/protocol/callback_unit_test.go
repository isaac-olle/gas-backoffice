package protocol_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	protocolHandler "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol"
)

type stubCreateSessionService struct {
	gotCode, gotState, gotIdpId string
	redirect                    string
	err                         error
}

func (s *stubCreateSessionService) CreateSession(_ context.Context, code, state, idpId string) (string, error) {
	s.gotCode, s.gotState, s.gotIdpId = code, state, idpId
	return s.redirect, s.err
}

func newCallbackContext(w *httptest.ResponseRecorder, target string) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, target, nil)
	c.Params = gin.Params{{Key: "idp_id", Value: "google"}}
	return c
}

func TestCallbackHandler(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		target       string
		svcErr       error
		svcRedirect  string
		wantErr      bool
		wantGotCode  string
		wantGotEmpty bool
	}{
		{
			name:         "missing_code_returns_error_without_calling_create_session",
			target:       "/v1/oidc/callback/google?state=s1",
			wantErr:      true,
			wantGotEmpty: true,
		},
		{
			name:         "missing_state_returns_error_without_calling_create_session",
			target:       "/v1/oidc/callback/google?code=c1",
			wantErr:      true,
			wantGotEmpty: true,
		},
		{
			name:    "create_session_error_is_propagated",
			target:  "/v1/oidc/callback/google?code=c1&state=s1",
			svcErr:  assert.AnError,
			wantErr: true,
		},
		{
			name:        "success_redirects_to_the_result_from_create_session",
			target:      "/v1/oidc/callback/google?code=c1&state=s1",
			svcRedirect: "https://client.example.com/cb?code=gas-code&state=s1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := &stubCreateSessionService{err: tt.svcErr, redirect: tt.svcRedirect}
			h := protocolHandler.NewCallbackHandler(svc)
			w := httptest.NewRecorder()
			c := newCallbackContext(w, tt.target)

			result, err := h.Execute(c)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantGotEmpty {
					assert.Empty(t, svc.gotCode)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "c1", svc.gotCode)
			assert.Equal(t, "s1", svc.gotState)
			assert.Equal(t, "google", svc.gotIdpId)

			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusFound, w.Code)
			assert.Equal(t, "https://client.example.com/cb?code=gas-code&state=s1", w.Header().Get("Location"))
		})
	}
}

func TestCallbackHandler_Default(t *testing.T) {
	t.Parallel()

	h := protocolHandler.NewCallbackHandler(&stubCreateSessionService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	h.Default(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
