package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/isaac-olle/gas-backoffice/internal/shared/request_id"
)

const RequestIDHeader = request_id.Header

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(RequestIDHeader, requestID)
		c.Header(RequestIDHeader, requestID)
		c.Request = c.Request.WithContext(request_id.WithContext(c.Request.Context(), requestID))
		c.Next()
	}
}
