package middlewares

import "github.com/gin-gonic/gin"

// Handler interface that all handlers must implement
type Handler interface {
	Execute(c *gin.Context) (any, error)
	Success(c *gin.Context, result any) error
	Default(c *gin.Context)
}

// HandlerMiddleware wraps handlers to convert them to gin.HandlerFunc
type HandlerMiddleware struct{}

func NewHandlerMiddleware() *HandlerMiddleware {
	return &HandlerMiddleware{}
}

func (e *HandlerMiddleware) Run(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.Execute(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		if err = h.Success(c, result); err != nil {
			h.Default(c)
		}
	}
}
