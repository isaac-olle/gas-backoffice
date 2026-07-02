package authentication_methods

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/query"

	verifySecret "github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets/application/verify_secret"
)

type BodyCredentialsMethod struct {
	queryBus query.Bus
}

func NewBodyCredentialsMethod(queryBus query.Bus) *BodyCredentialsMethod {
	return &BodyCredentialsMethod{queryBus: queryBus}
}

func (m *BodyCredentialsMethod) CanProcess(c *gin.Context) (bool, error) {
	matches := strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded") &&
		c.PostForm("client_id") != "" &&
		c.PostForm("client_secret") != ""
	if matches && c.GetBool(authenticatedKey) {
		return false, app_errors.BadRequestError(ErrAlreadyAuthenticated)
	}
	return matches, nil
}

func (m *BodyCredentialsMethod) Authenticate(c *gin.Context) (string, error) {
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")
	if _, err := m.queryBus.Dispatch(c.Request.Context(), verifySecret.NewQuery(clientID, clientSecret)); err != nil {
		return "", err
	}
	return clientID, nil
}
