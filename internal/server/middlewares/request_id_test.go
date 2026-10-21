package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/shared/request_id"
)

func TestRequestID_GeneratesIDWhenNotProvided(t *testing.T) {
	t.Parallel()
	r := gin.New()
	r.Use(middlewares.RequestID())
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	requestID := w.Header().Get(middlewares.RequestIDHeader)
	assert.NotEmpty(t, requestID)
}

func TestRequestID_UsesExistingHeaderWhenProvided(t *testing.T) {
	t.Parallel()
	r := gin.New()
	r.Use(middlewares.RequestID())
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(middlewares.RequestIDHeader, "my-custom-request-id")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "my-custom-request-id", w.Header().Get(middlewares.RequestIDHeader))
}

func TestRequestID_SetsIDInGinContext(t *testing.T) {
	t.Parallel()
	var capturedID string
	r := gin.New()
	r.Use(middlewares.RequestID())
	r.GET("/test", func(c *gin.Context) {
		val, exists := c.Get(middlewares.RequestIDHeader)
		require.True(t, exists)
		capturedID = val.(string)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(middlewares.RequestIDHeader, "test-id-123")
	r.ServeHTTP(w, req)

	assert.Equal(t, "test-id-123", capturedID)
}

func TestRequestID_SetsIDInRequestContext(t *testing.T) {
	t.Parallel()
	var capturedID string
	r := gin.New()
	r.Use(middlewares.RequestID())
	r.GET("/test", func(c *gin.Context) {
		capturedID = request_id.FromContext(c.Request.Context())
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(middlewares.RequestIDHeader, "context-id-456")
	r.ServeHTTP(w, req)

	assert.Equal(t, "context-id-456", capturedID)
}

func TestRequestID_SetsIDAsResponseHeader(t *testing.T) {
	t.Parallel()
	r := gin.New()
	r.Use(middlewares.RequestID())
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(middlewares.RequestIDHeader, "response-header-check")
	r.ServeHTTP(w, req)

	assert.Equal(t, "response-header-check", w.Header().Get(middlewares.RequestIDHeader))
}
