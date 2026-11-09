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

func TestRequestID(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "generates an id when no header is provided", run: testRequestIDGeneratesWhenMissing},
		{name: "preserves the provided header", run: testRequestIDUsesExistingHeader},
		{name: "sets the gin context value", run: testRequestIDSetsGinContext},
		{name: "sets the request context value", run: testRequestIDSetsRequestContext},
		{name: "returns the id in the response header", run: testRequestIDSetsResponseHeader},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func testRequestIDGeneratesWhenMissing(t *testing.T) {
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

func testRequestIDUsesExistingHeader(t *testing.T) {
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

func testRequestIDSetsGinContext(t *testing.T) {
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

func testRequestIDSetsRequestContext(t *testing.T) {
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

func testRequestIDSetsResponseHeader(t *testing.T) {
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
