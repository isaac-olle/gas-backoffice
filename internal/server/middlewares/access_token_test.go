package middlewares_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	querymocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/query/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
)

func TestAccessToken(t *testing.T) {
	t.Parallel()
	validAuthor, _ := authors.NewAdminAuthorFromString("b1b2c3d4-e5f6-7890-abcd-000000000001")

	tests := []struct {
		name       string
		authHeader string
		busResult  any
		busErr     error
		busCalled  bool
		wantNext   bool
	}{
		{
			name:       "no_auth_header_aborts",
			authHeader: "",
			busCalled:  false,
			wantNext:   false,
		},
		{
			name:       "non_bearer_header_aborts",
			authHeader: "Basic sometoken",
			busCalled:  false,
			wantNext:   false,
		},
		{
			name:       "bus_error_aborts",
			authHeader: "Bearer sometoken",
			busErr:     errors.New("invalid token"),
			busCalled:  true,
			wantNext:   false,
		},
		{
			name:       "valid_token_sets_actor_and_continues",
			authHeader: "Bearer validtoken",
			busResult:  vo.Author(validAuthor),
			busCalled:  true,
			wantNext:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			queryBus := &querymocks.MockQueryBus{}
			if tt.busCalled {
				queryBus.On("Dispatch", mock.Anything, mock.Anything).Return(tt.busResult, tt.busErr)
			}

			r := gin.New()
			r.GET("/test", middlewares.AccessToken(queryBus), func(c *gin.Context) {
				nextCalled = true
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantNext, nextCalled)
		})
	}
}
