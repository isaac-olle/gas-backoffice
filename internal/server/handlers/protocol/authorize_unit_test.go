package protocol_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authorize"
	handlerprotocol "github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type stubAuthorizeService struct {
	url string
	err error
}

func (s *stubAuthorizeService) Authorize(_ context.Context, _ authorize.Request) (string, error) {
	return s.url, s.err
}

func newAuthorizeHandler(url string, err error) *handlerprotocol.AuthorizeHandler {
	return handlerprotocol.NewAuthorizeHandler(&stubAuthorizeService{url: url, err: err})
}

func TestAuthorizeHandler(t *testing.T) {
	t.Parallel()

	validQuery := "response_type=code&client_id=abc&redirect_uri=https://example.com/cb&scope=openid&state=xyz"
	googleURL := "https://accounts.google.com/o/oauth2/v2/auth?state=xyz&client_id=client-id"

	tests := []struct {
		name         string
		query        string
		serviceErr   error
		urlResult    string
		wantErr      bool
		wantStatus   int
		wantLocation string
	}{
		{
			name:    "missing_required_params",
			query:   "",
			wantErr: true,
		},
		{
			name:       "authorize_service_error",
			query:      validQuery,
			serviceErr: assert.AnError,
			wantErr:    true,
		},
		{
			name:         "success",
			query:        validQuery,
			urlResult:    googleURL,
			wantStatus:   http.StatusFound,
			wantLocation: googleURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := newAuthorizeHandler(tt.urlResult, tt.serviceErr)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/?"+tt.query, nil)

			result, err := h.Execute(c)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "missing_required_params" {
					var appErr *app_errors.AppError
					require.ErrorAs(t, err, &appErr)
					assert.Equal(t, app_errors.ErrInvalidRequestBody, appErr.Code)
				}
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantLocation, w.Header().Get("Location"))
		})
	}
}

func TestAuthorizeHandler_Default(t *testing.T) {
	t.Parallel()

	h := newAuthorizeHandler("", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	h.Default(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
