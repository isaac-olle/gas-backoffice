package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// maxLoggedBodyBytes caps how much of a request/response body gets kept in
// memory for logging, so a large upload/download doesn't blow up RAM.
const maxLoggedBodyBytes = 8 * 1024

const redacted = "***"

// sensitiveKeyParts matches field names case-insensitively by substring, so
// it also catches "client_secret", "plaintext_secret", "secret_hash",
// "access_token", "refresh_token", "code_verifier", etc.
var sensitiveKeyParts = []string{"password", "secret", "token", "code_verifier", "authorization"}

func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, part := range sensitiveKeyParts {
		if strings.Contains(lower, part) {
			return true
		}
	}
	return false
}

// redactBody masks sensitive field values in a JSON body before it's logged.
// Non-JSON bodies (empty, form-encoded, etc.) are returned unchanged — this
// only ever touches the copy kept for logging, never the body the handler
// actually binds.
func redactBody(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return raw
	}

	var parsed any
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return raw
	}

	redactValue(parsed)

	out, err := json.Marshal(parsed)
	if err != nil {
		return raw
	}
	return string(out)
}

func redactValue(v any) {
	switch val := v.(type) {
	case map[string]any:
		for key, nested := range val {
			if isSensitiveKey(key) {
				val[key] = redacted
				continue
			}
			redactValue(nested)
		}
	case []any:
		for _, item := range val {
			redactValue(item)
		}
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	if w.buf.Len() < maxLoggedBodyBytes {
		w.buf.Write(b[:min(len(b), maxLoggedBodyBytes-w.buf.Len())])
	}
	return w.ResponseWriter.Write(b)
}

// BodyLogger captures the request and response bodies so ginLogFormatter can
// print them. Local-dev debugging only: request bodies can contain secrets
// (client_secret, passwords, tokens) and will be printed in plain text to
// stdout — never enable this in a deployed environment.
func BodyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxLoggedBodyBytes))
			if err == nil {
				c.Set("request_body", redactBody(string(body)))
			}
			// Restore the body (possibly truncated) for the actual handler to bind.
			rest, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(body), bytes.NewReader(rest)))
		}

		writer := &bodyWriter{ResponseWriter: c.Writer, buf: &bytes.Buffer{}}
		c.Writer = writer

		c.Next()

		c.Set("response_body", redactBody(writer.buf.String()))
	}
}
