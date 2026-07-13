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

	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/ownership"
	querymocks "github.com/isaac-olle/gas-backoffice/internal/shared/buses/query/mocks"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func TestOwnershipClient_QuerySucceeds_SetsOwnershipAndCallsNext(t *testing.T) {
	t.Parallel()
	expectedOwnership := "client-ownership"
	bus := &querymocks.MockQueryBus{}
	bus.On("Dispatch", mock.Anything, mock.Anything).Return(expectedOwnership, nil)

	var capturedCtx context.Context
	r := gin.New()
	r.Use(func(c *gin.Context) {
		author, _ := authors.NewAdminAuthorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(vo2.WithAuthor(c.Request.Context(), author))
		userId, _ := ids2.NewUserIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithUserId(c.Request.Context(), userId))
		clientId, _ := ids2.NewClientIdFromString("c1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithClientId(c.Request.Context(), clientId))
		c.Next()
	})
	r.Use(ownership.Client(bus))
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

func TestOwnershipClient_QueryFails_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	bus.On("Dispatch", mock.Anything, mock.Anything).Return(nil, errors.New("forbidden"))
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		author, _ := authors.NewAdminAuthorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(vo2.WithAuthor(c.Request.Context(), author))
		userId, _ := ids2.NewUserIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithUserId(c.Request.Context(), userId))
		clientId, _ := ids2.NewClientIdFromString("c1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithClientId(c.Request.Context(), clientId))
		c.Next()
	})
	r.Use(ownership.Client(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipClient_MissingActor_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		userId, _ := ids2.NewUserIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithUserId(c.Request.Context(), userId))
		clientId, _ := ids2.NewClientIdFromString("c1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithClientId(c.Request.Context(), clientId))
		c.Next()
	})
	r.Use(ownership.Client(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipClient_MissingClientId_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		author, _ := authors.NewAdminAuthorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(vo2.WithAuthor(c.Request.Context(), author))
		userId, _ := ids2.NewUserIdFromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithUserId(c.Request.Context(), userId))
		c.Next()
	})
	r.Use(ownership.Client(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestOwnershipClient_MissingUserId_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		author, _ := authors.NewAdminAuthorFromString("b1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(vo2.WithAuthor(c.Request.Context(), author))
		clientId, _ := ids2.NewClientIdFromString("c1b2c3d4-e5f6-7890-abcd-ef1234567890")
		c.Request = c.Request.WithContext(ids2.WithClientId(c.Request.Context(), clientId))
		c.Next()
	})
	r.Use(ownership.Client(bus))
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}
