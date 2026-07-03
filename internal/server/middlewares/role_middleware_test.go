package middlewares_test

import (
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	actors2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/actors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
)

const testAdminID = "b1b2c3d4-e5f6-7890-abcd-ef1234567890"
const testTenantID = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

func withAdminActor(c *gin.Context) {
	actor, _ := actors2.NewAdminActorFromString(testAdminID)
	c.Request = c.Request.WithContext(vo2.WithActor(c.Request.Context(), actor))
}

func withTenantActor(c *gin.Context) {
	actor, _ := actors2.NewTenantActorFromString(testTenantID)
	c.Request = c.Request.WithContext(vo2.WithActor(c.Request.Context(), actor))
}

func TestRoleMiddleware(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		setupActor   func(c *gin.Context)
		allowedRoles []vo2.Role
		wantNext     bool
	}{
		{
			name:         "no_actor_aborts",
			setupActor:   nil,
			allowedRoles: []vo2.Role{vo2.AdminRole},
			wantNext:     false,
		},
		{
			name:         "correct_role_admin_passes",
			setupActor:   withAdminActor,
			allowedRoles: []vo2.Role{vo2.AdminRole},
			wantNext:     true,
		},
		{
			name:         "wrong_role_aborts",
			setupActor:   withAdminActor,
			allowedRoles: []vo2.Role{vo2.TenantRole},
			wantNext:     false,
		},
		{
			name:         "multiple_allowed_tenant_in_list_passes",
			setupActor:   withTenantActor,
			allowedRoles: []vo2.Role{vo2.AdminRole, vo2.TenantRole},
			wantNext:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false

			r := gin.New()
			r.GET("/test",
				func(c *gin.Context) {
					if tt.setupActor != nil {
						tt.setupActor(c)
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
