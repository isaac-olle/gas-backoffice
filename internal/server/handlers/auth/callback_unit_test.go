package auth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/create_session"
	authHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/auth"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type stubQueryBus struct {
	result     any
	err        error
	dispatched bool
}

func (s *stubQueryBus) Dispatch(_ context.Context, _ query.Query) (any, error) {
	s.dispatched = true
	return s.result, s.err
}

func (s *stubQueryBus) Register(query.Type, query.Handler) {}

func TestCallbackHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		body         string
		dispatchErr  error
		wantErr      bool
		wantDispatch bool
	}{
		{
			name:    "invalid_json",
			body:    `{`,
			wantErr: true,
		},
		{
			name:    "missing_code",
			body:    `{"state":"xyz"}`,
			wantErr: true,
		},
		{
			name:    "missing_state",
			body:    `{"code":"google-code"}`,
			wantErr: true,
		},
		{
			name:         "dispatch_error",
			body:         `{"code":"google-code","state":"xyz"}`,
			dispatchErr:  errors.New("create session failed"),
			wantErr:      true,
			wantDispatch: true,
		},
		{
			name:         "success",
			body:         `{"idp_reference_name":"google","code":"google-code","state":"xyz"}`,
			wantDispatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bus := &stubQueryBus{result: create_session.Result{AccessToken: "the-token"}, err: tt.dispatchErr}
			h := authHandlers.NewCallbackHandler(bus)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPost, "/v1/backoffice/auth/callback", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			result, err := h.Execute(c)
			assert.Equal(t, tt.wantDispatch, bus.dispatched)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Contains(t, w.Header().Get("Set-Cookie"), authenticators.AccessTokenCookieName+"=the-token")
		})
	}
}

func TestCallbackHandler_Default(t *testing.T) {
	t.Parallel()

	h := authHandlers.NewCallbackHandler(&stubQueryBus{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/v1/backoffice/auth/callback", nil)

	h.Default(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}
