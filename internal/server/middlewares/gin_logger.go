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
	ts := p.TimeStamp.Format("15:04:05")

	line := fmt.Sprintf("[GIN] %s %s %s %s %d %s", ts, p.ClientIP, p.Method, p.Path, p.StatusCode, p.Latency)
	if location, ok := p.Keys["redirect"].(string); ok && location != "" {
		line += " → " + location
	}

	if reqBody, ok := p.Keys["request_body"].(string); ok && reqBody != "" {
		line += fmt.Sprintf("\n  request_body:  %s", reqBody)
	}
	if respBody, ok := p.Keys["response_body"].(string); ok && respBody != "" {
		line += fmt.Sprintf("\n  response_body: %s", respBody)
	}

	return line + "\n"
}
