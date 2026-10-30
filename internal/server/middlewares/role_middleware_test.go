package middlewares_test

import (
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
)

const testAdminID = "b1b2c3d4-e5f6-7890-abcd-ef1234567890"
const testUserId = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

func withAdminAuthor(c *gin.Context) {
	author, _ := vo2.NewUserAuthorFromString(testAdminID, testAdminID, vo2.SystemAdministratorRole)
	c.Request = c.Request.WithContext(vo2.WithAuthor(c.Request.Context(), author))
}

func withTenantAuthor(c *gin.Context) {
	author, _ := vo2.NewUserAuthorFromString(testUserId, testAdminID, vo2.SubscriptionUserRole)
	c.Request = c.Request.WithContext(vo2.WithAuthor(c.Request.Context(), author))
}

func TestRoleMiddleware(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		setupAuthor  func(c *gin.Context)
		allowedRoles []vo2.Role
		wantNext     bool
	}{
		{
			name:         "no_actor_aborts",
			setupAuthor:  nil,
			allowedRoles: []vo2.Role{vo2.SystemAdministratorRole},
			wantNext:     false,
		},
		{
			name:         "correct_role_admin_passes",
			setupAuthor:  withAdminAuthor,
			allowedRoles: []vo2.Role{vo2.SystemAdministratorRole},
			wantNext:     true,
		},
		{
			name:         "wrong_role_aborts",
			setupAuthor:  withAdminAuthor,
			allowedRoles: []vo2.Role{vo2.SubscriptionUserRole},
			wantNext:     false,
		},
		{
			name:         "multiple_allowed_user_in_list_passes",
			setupAuthor:  withTenantAuthor,
			allowedRoles: []vo2.Role{vo2.SystemAdministratorRole, vo2.SubscriptionUserRole},
			wantNext:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false

			r := gin.New()
			r.GET("/test",
				func(c *gin.Context) {
					if tt.setupAuthor != nil {
						tt.setupAuthor(c)
					}
					c.Next()
				},
				middlewares.RoleMiddleware(tt.allowedRoles...),
				func(c *gin.Context) {
					nextCalled = true
					c.Status(http.StatusOK)
				},
			)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantNext, nextCalled)
		})
	}
}
