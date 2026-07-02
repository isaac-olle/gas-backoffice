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

	"github.com/isaac-olle/gas/internal/server/middlewares/ownership"
	querymocks "github.com/isaac-olle/gas/internal/shared/buses/query/mocks"
	vo2 "github.com/isaac-olle/gas/internal/shared/domain/vo"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/actors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

func fullScopeSetup(c *gin.Context) {
	actor, _ := actors.NewAdminActorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
	c.Request = c.Request.WithContext(vo2.WithActor(c.Request.Context(), actor))
	tenantId, _ := ids2.NewTenantIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	c.Request = c.Request.WithContext(ids2.WithTenantId(c.Request.Context(), tenantId))
	prId, _ := ids2.NewProtectedResourceIdFromString("d1b2c3d4-e5f6-7890-abcd-ef1234567890")
	c.Request = c.Request.WithContext(ids2.WithProtectedResourceId(c.Request.Context(), prId))
	scopeId, _ := ids2.NewScopeIdFromString("e1b2c3d4-e5f6-7890-abcd-ef1234567890")
	c.Request = c.Request.WithContext(ids2.WithScopeId(c.Request.Context(), scopeId))
}

func TestOwnershipScope_QuerySucceeds_SetsOwnershipAndCallsNext(t *testing.T) {
	t.Parallel()
	expectedOwnership := "scope-ownership"
	bus := &querymocks.MockQueryBus{}
	bus.On("Dispatch", mock.Anything, mock.Anything).Return(expectedOwnership, nil)

	var capturedCtx context.Context
	r := gin.New()
	r.Use(func(c *gin.Context) { fullScopeSetup(c); c.Next() })
	r.Use(ownership.Scope(bus))
	r.GET("/test", func(c *gin.Context) {
		capturedCtx = c.Request.Context()
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	result, ok := vo2.OwnershipFromContext[string](capturedCtx)
	assert.True(t, ok)
	assert.Equal(t, expectedOwnership, result)
}

func TestOwnershipScope_QueryFails_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	bus.On("Dispatch", mock.Anything, mock.Anything).Return(nil, errors.New("forbidden"))
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) { fullScopeSetup(c); c.Next() })
	r.Use(ownership.Scope(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipScope_MissingActor_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		tenantId, _ := ids2.NewTenantIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithTenantId(c.Request.Context(), tenantId))
		prId, _ := ids2.NewProtectedResourceIdFromString("d1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithProtectedResourceId(c.Request.Context(), prId))
		scopeId, _ := ids2.NewScopeIdFromString("e1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithScopeId(c.Request.Context(), scopeId))
		c.Next()
	})
	r.Use(ownership.Scope(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipScope_MissingScopeId_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		actor, _ := actors.NewAdminActorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(vo2.WithActor(c.Request.Context(), actor))
		tenantId, _ := ids2.NewTenantIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithTenantId(c.Request.Context(), tenantId))
		prId, _ := ids2.NewProtectedResourceIdFromString("d1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithProtectedResourceId(c.Request.Context(), prId))
		c.Next()
	})
	r.Use(ownership.Scope(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipScope_MissingProtectedResourceId_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		actor, _ := actors.NewAdminActorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(vo2.WithActor(c.Request.Context(), actor))
		tenantId, _ := ids2.NewTenantIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithTenantId(c.Request.Context(), tenantId))
		scopeId, _ := ids2.NewScopeIdFromString("e1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithScopeId(c.Request.Context(), scopeId))
		c.Next()
	})
	r.Use(ownership.Scope(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipScope_MissingTenantId_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		actor, _ := actors.NewAdminActorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(vo2.WithActor(c.Request.Context(), actor))
		prId, _ := ids2.NewProtectedResourceIdFromString("d1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithProtectedResourceId(c.Request.Context(), prId))
		scopeId, _ := ids2.NewScopeIdFromString("e1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithScopeId(c.Request.Context(), scopeId))
		c.Next()
	})
	r.Use(ownership.Scope(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}
