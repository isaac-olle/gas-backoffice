package ownership_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	userDomain "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/ownership"
	querymocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/query/mocks"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func seedTenantId(tenantId ids2.UserId) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(ids2.WithUserId(c.Request.Context(), tenantId))
		c.Next()
	}
}

func seedAdminAuthor(adminId ids2.UserId) gin.HandlerFunc {
	return func(c *gin.Context) {
		author := &authors.UserAuthor{Id: adminId, Role: vo2.AdminRole}
		c.Request = c.Request.WithContext(vo2.WithAuthor(c.Request.Context(), author))
		c.Next()
	}
}

func TestOwnershipTenant_QuerySucceeds_SetsOwnership(t *testing.T) {
	t.Parallel()
	tenantId := ids2.NewUserId()
	adminId := ids2.NewUserId()
	expectedOwnership := userDomain.NewOwnership(tenantId)

	bus := &querymocks.MockQueryBus{}
	bus.On("Dispatch", mock.Anything, mock.Anything).Return(expectedOwnership, nil)

	var capturedCtx context.Context
	r := gin.New()
	r.Use(seedTenantId(tenantId), seedAdminAuthor(adminId))
	r.Use(ownership.Tenant(bus))
	r.GET("/test", func(c *gin.Context) {
		capturedCtx = c.Request.Context()
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	result, ok := vo2.OwnershipFromContext[userDomain.Ownership](capturedCtx)
	assert.True(t, ok)
	assert.Equal(t, expectedOwnership, result)
}

func TestOwnershipTenant_QueryFails_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	bus.On("Dispatch", mock.Anything, mock.Anything).Return(nil, errors.New("forbidden"))
	nextCalled := false

	r := gin.New()
	r.Use(seedTenantId(ids2.NewUserId()), seedAdminAuthor(ids2.NewUserId()))
	r.Use(ownership.Tenant(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipTenant_MissingTenantId_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(seedAdminAuthor(ids2.NewUserId()))
	r.Use(ownership.Tenant(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipTenant_MissingAuthor_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(seedTenantId(ids2.NewUserId()))
	r.Use(ownership.Tenant(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}
