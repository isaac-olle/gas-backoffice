package middlewares_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type stubPresenter struct {
	called bool
	err    *app_errors.AppError
}

func (s *stubPresenter) Present(_ *gin.Context, err *app_errors.AppError) {
	s.called = true
	s.err = err
}

type stubManager struct {
	called bool
}

func (s *stubManager) Manage(_ context.Context, _ *app_errors.AppError) {
	s.called = true
}

func TestHttpError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		handlerErr     error
		alreadyWritten bool
		wantPresent    bool
		wantManage     bool
	}{
		{
			name:        "no_error_skips_presenter",
			wantPresent: false,
			wantManage:  false,
		},
		{
			name:        "app_error_calls_manager_and_presenter",
			handlerErr:  app_errors.BadRequestError("TEST_ERROR"),
			wantPresent: true,
			wantManage:  true,
		},
		{
			name:        "non_app_error_wraps_as_internal_and_calls_presenter",
			handlerErr:  errors.New("raw error"),
			wantPresent: true,
			wantManage:  true,
		},
		{
			name:           "already_written_response_skips_presenter",
			handlerErr:     app_errors.BadRequestError("TEST_ERROR"),
			alreadyWritten: true,
			wantPresent:    false,
			wantManage:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			presenter := &stubPresenter{}
			manager := &stubManager{}

			r := gin.New()
			r.GET("/test",
				middlewares.HttpError(presenter, manager),
				func(c *gin.Context) {
					if tt.alreadyWritten {
						c.Data(http.StatusOK, "text/plain", []byte("ok"))
					}
					if tt.handlerErr != nil {
						_ = c.Error(tt.handlerErr)
					}
				},
			)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantPresent, presenter.called)
			assert.Equal(t, tt.wantManage, manager.called)
		})
	}
}
