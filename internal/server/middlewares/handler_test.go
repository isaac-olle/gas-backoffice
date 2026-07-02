package middlewares_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/isaac-olle/gas/internal/server/middlewares"
)

// stubHandler records which methods were called and can be configured to return errors.
type stubHandler struct {
	executeErr    error
	successErr    error
	executeCalled bool
	successCalled bool
	defaultCalled bool
}

func (s *stubHandler) Execute(_ *gin.Context) (any, error) {
	s.executeCalled = true
	if s.executeErr != nil {
		return nil, s.executeErr
	}
	return "ok", nil
}

func (s *stubHandler) Success(_ *gin.Context, _ any) error {
	s.successCalled = true
	return s.successErr
}

func (s *stubHandler) Default(_ *gin.Context) {
	s.defaultCalled = true
}

func TestHandlerMiddlewareRun(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		executeErr        error
		successErr        error
		wantAbort         bool
		wantNextCalled    bool
		wantDefaultCalled bool
	}{
		{
			name:           "execute_error_aborts_chain",
			executeErr:     errors.New("execute failed"),
			wantAbort:      true,
			wantNextCalled: false,
		},
		{
			name:           "success_continues_chain",
			wantAbort:      false,
			wantNextCalled: true,
		},
		{
			name:              "success_error_calls_default_but_does_not_abort",
			successErr:        errors.New("success failed"),
			wantAbort:         false,
			wantNextCalled:    true,
			wantDefaultCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &stubHandler{
				executeErr: tt.executeErr,
				successErr: tt.successErr,
			}

			nextCalled := false
			m := middlewares.NewHandlerMiddleware()

			r := gin.New()
			r.GET("/test", m.Run(h), func(c *gin.Context) {
				nextCalled = true
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantNextCalled, nextCalled, "next handler called")
			assert.Equal(t, tt.wantDefaultCalled, h.defaultCalled, "Default called")
			assert.True(t, h.executeCalled, "Execute should always be called")

			if tt.wantAbort {
				// when aborted, next is not called — verified above
				assert.False(t, nextCalled)
			}
		})
	}
}
