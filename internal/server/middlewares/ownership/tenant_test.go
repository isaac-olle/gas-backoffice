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
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/actors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func init() { gin.SetMode(gin.TestMode) }

func TestOwnershipUser_QuerySucceeds_SetsOwnershipAndCallsNext(t *testing.T) {
	t.Parallel()
	adminId, _ := ids2.NewAdminIdFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
	expectedOwnership := userDomain.NewOwnership(adminId)
	bus := &querymocks.MockQueryBus{}
	bus.On("Dispatch", mock.Anything, mock.Anything).Return(expectedOwnership, nil)

	var capturedCtx context.Context
	r := gin.New()
	r.Use(func(c *gin.Context) {
		actor, _ := actors.NewAdminActorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(vo2.WithActor(c.Request.Context(), actor))
		userId, _ := ids2.NewUserIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithUserId(c.Request.Context(), userId))
		c.Next()
	})
	r.Use(ownership.User(bus))
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

func TestOwnershipUser_QueryFails_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	bus.On("Dispatch", mock.Anything, mock.Anything).Return(nil, errors.New("forbidden"))
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		actor, _ := actors.NewAdminActorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(vo2.WithActor(c.Request.Context(), actor))
		userId, _ := ids2.NewUserIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithUserId(c.Request.Context(), userId))
		c.Next()
	})
	r.Use(ownership.User(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipUser_MissingActor_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		userId, _ := ids2.NewUserIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithUserId(c.Request.Context(), userId))
		c.Next()
	})
	r.Use(ownership.User(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipUser_MissingUserId_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		actor, _ := actors.NewAdminActorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(vo2.WithActor(c.Request.Context(), actor))
		c.Next()
	})
	r.Use(ownership.User(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}
