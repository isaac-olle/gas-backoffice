package ownership_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	userDomain "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/ownership"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func seedUserAuthor(userId ids2.UserId) gin.HandlerFunc {
	return func(c *gin.Context) {
		author := &authors.UserAuthor{Id: userId, Role: vo2.TenantRole}
		c.Request = c.Request.WithContext(vo2.WithAuthor(c.Request.Context(), author))
		c.Next()
	}
}

func TestOwnershipUser_SetsOwnershipFromAuthor(t *testing.T) {
	t.Parallel()
	userId := ids2.NewUserId()
	expectedOwnership := userDomain.NewOwnership(userId)

	var capturedCtx context.Context
	r := gin.New()
	r.Use(seedUserAuthor(userId))
	r.Use(ownership.User())
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

func TestOwnershipUser_MissingAuthor_Aborts(t *testing.T) {
	t.Parallel()
	nextCalled := false

	r := gin.New()
	r.Use(ownership.User())
	r.GET("/test", func(c *gin.Context) { nextCalled = true; c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}
