package authentication_methods

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
)

const authenticatedKey = "request_authenticated"

var ErrAlreadyAuthenticated app_errors.ErrorCode = "ALREADY_AUTHENTICATED"

type AuthenticationMethod interface {
	CanProcess(c *gin.Context) (bool, error)
	Authenticate(c *gin.Context) (clientID string, err error)
}

func MarkAuthenticated(c *gin.Context) {
	c.Set(authenticatedKey, true)
}
