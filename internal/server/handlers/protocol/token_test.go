package protocol_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
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

func TestTokenHandler(t *testing.T) {
	t.Parallel()

	tokenResp := &session.TokenResponse{Token: "access-token", Scope: "openid"}

	tests := []struct {
		name       string
		form       url.Values
		qryResult  any
		qryErr     error
		wantErr    bool
		wantStatus int
	}{
		{
			name:    "missing_grant_type",
			form:    url.Values{},
			wantErr: true,
		},
		{
			name:    "unsupported_grant_type",
			form:    url.Values{"grant_type": {"implicit"}},
			wantErr: true,
		},
		{
			name:    "auth_code_missing_code",
			form:    url.Values{"grant_type": {"authorization_code"}, "redirect_uri": {"https://example.com/cb"}},
			wantErr: true,
		},
		{
			name:    "auth_code_query_bus_error",
			form:    url.Values{"grant_type": {"authorization_code"}, "code": {"xyz"}, "redirect_uri": {"https://example.com/cb"}},
			qryErr:  assert.AnError,
			wantErr: true,
		},
		{
			name:       "auth_code_success",
			form:       url.Values{"grant_type": {"authorization_code"}, "code": {"xyz"}, "redirect_uri": {"https://example.com/cb"}},
			qryResult:  tokenResp,
			wantStatus: http.StatusOK,
		},
		{
			name:    "client_credentials_missing_scope",
			form:    url.Values{"grant_type": {"client_credentials"}},
			wantErr: true,
		},
		{
			name:    "client_credentials_query_bus_error",
			form:    url.Values{"grant_type": {"client_credentials"}, "scope": {"openid"}},
			qryErr:  assert.AnError,
			wantErr: true,
		},
		{
			name:       "client_credentials_success",
			form:       url.Values{"grant_type": {"client_credentials"}, "scope": {"openid"}},
			qryResult:  tokenResp,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := protocol.NewTokenHandler(handlerstest.NewQueryBus(tt.qryResult, tt.qryErr))
			w, c := tokenCtx(tt.form)

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

	h := protocol.NewTokenHandler(handlerstest.NewQueryBus(nil, nil))
	_, c := tokenCtx(url.Values{})

	h.Default(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
