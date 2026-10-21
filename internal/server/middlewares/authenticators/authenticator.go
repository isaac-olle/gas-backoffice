package authenticators

import "github.com/gin-gonic/gin"

// Authenticator extracts the raw access token from a request. Each
// implementation reads it from a different transport; verifying the token is
// not its concern.
type Authenticator interface {
	Token(c *gin.Context) (string, error)
}
