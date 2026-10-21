package protocol_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/token"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol"
)

func formBody(values url.Values) *strings.Reader {
	return strings.NewReader(values.Encode())
}

func tokenCtx(form url.Values) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", formBody(form))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return w, c
}

func seedAuth(c *gin.Context, clientID string, isPublic bool) {
	c.Request = c.Request.WithContext(
		protocol.WithAuthenticatedClient(c.Request.Context(), clientID, isPublic),
	)
}

type stubTokenService struct {
	resp *session.TokenResponse
	err  error
}

func (s *stubTokenService) Execute(_ context.Context, _ string, _ token.Request) (*session.TokenResponse, error) {
	return s.resp, s.err
}

func TestTokenHandler(t *testing.T) {
	t.Parallel()

	tokenResp := &session.TokenResponse{AccessToken: "access-token", Scope: "openid"}

	tests := []struct {
		name       string
		form       url.Values
		isPublic   bool
		serviceErr error
		resp       *session.TokenResponse
		wantErr    bool
		wantStatus int
	}{
		{
			name:    "missing_grant_type",
			form:    url.Values{},
			wantErr: true,
		},
		{
			name:       "unsupported_grant_type",
			form:       url.Values{"grant_type": {"implicit"}},
			serviceErr: assert.AnError,
			wantErr:    true,
		},
		{
			name:    "auth_code_missing_code",
			form:    url.Values{"grant_type": {"authorization_code"}, "redirect_uri": {"https://example.com/cb"}},
			wantErr: true,
		},
		{
			name:       "auth_code_service_error",
			form:       url.Values{"grant_type": {"authorization_code"}, "code": {"xyz"}, "redirect_uri": {"https://example.com/cb"}},
			serviceErr: assert.AnError,
			wantErr:    true,
		},
		{
			name:       "auth_code_success",
			form:       url.Values{"grant_type": {"authorization_code"}, "code": {"xyz"}, "redirect_uri": {"https://example.com/cb"}},
			resp:       tokenResp,
			wantStatus: http.StatusOK,
		},
		{
			name:    "client_credentials_missing_scope",
			form:    url.Values{"grant_type": {"client_credentials"}},
			wantErr: true,
		},
		{
			name:       "client_credentials_service_error",
			form:       url.Values{"grant_type": {"client_credentials"}, "scope": {"openid"}},
			serviceErr: assert.AnError,
			wantErr:    true,
		},
		{
			name:       "client_credentials_success",
			form:       url.Values{"grant_type": {"client_credentials"}, "scope": {"openid"}},
			resp:       tokenResp,
			wantStatus: http.StatusOK,
		},
		{
			name:       "public_client_cannot_use_client_credentials",
			form:       url.Values{"grant_type": {"client_credentials"}, "scope": {"openid"}},
			isPublic:   true,
			serviceErr: assert.AnError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := protocol.NewTokenHandler(&stubTokenService{resp: tt.resp, err: tt.serviceErr})
			w, c := tokenCtx(tt.form)
			seedAuth(c, "test-client", tt.isPublic)

			result, err := h.Execute(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestTokenHandler_Default(t *testing.T) {
	t.Parallel()

	h := protocol.NewTokenHandler(&stubTokenService{})
	_, c := tokenCtx(url.Values{})

	h.Default(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
