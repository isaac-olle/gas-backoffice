package middlewares

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func NewGinLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Output:    os.Stdout,
		SkipPaths: []string{"/health"},
		Formatter: ginLogFormatter,
	})
}

func ginLogFormatter(p gin.LogFormatterParams) string {
	if location, ok := p.Keys["redirect"].(string); ok && location != "" {
		return fmt.Sprintf("[GIN] %s %s %d %s → %s\n", p.Method, p.Path, p.StatusCode, p.Latency, location)
	}
	return fmt.Sprintf("[GIN] %s %s %d %s\n", p.Method, p.Path, p.StatusCode, p.Latency)
}
