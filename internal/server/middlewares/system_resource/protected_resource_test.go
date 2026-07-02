package system_resource_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/isaac-olle/gas/internal/server/middlewares/system_resource"
	querymocks "github.com/isaac-olle/gas/internal/shared/buses/query/mocks"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

func withProtectedResourceId(c *gin.Context) {
	id, _ := ids2.NewProtectedResourceIdFromString("d1b2c3d4-e5f6-7890-abcd-ef1234567890")
	c.Request = c.Request.WithContext(ids2.WithProtectedResourceId(c.Request.Context(), id))
}

func TestSystemResourceProtectedResource_Success_CallsNext(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	bus.On("Dispatch", mock.Anything, mock.Anything).Return(nil, nil)
	nextCalled := false

	r := gin.New()
	r.GET("/test", func(c *gin.Context) { withProtectedResourceId(c); c.Next() }, system_resource.ProtectedResource(bus), func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.True(t, nextCalled)
}

func TestSystemResourceProtectedResource_MissingId_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	nextCalled := false

	r := gin.New()
	r.GET("/test", system_resource.ProtectedResource(bus), func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}

func TestSystemResourceProtectedResource_BusError_Aborts(t *testing.T) {
	t.Parallel()
	bus := &querymocks.MockQueryBus{}
	bus.On("Dispatch", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
	nextCalled := false

	r := gin.New()
	r.GET("/test", func(c *gin.Context) { withProtectedResourceId(c); c.Next() }, system_resource.ProtectedResource(bus), func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.False(t, nextCalled)
}
