package authentication_methods

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/query"

	verifySecret "github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets/application/verify_secret"
)

var ErrInvalidBasicAuth app_errors.ErrorCode = "INVALID_BASIC_AUTH"

type BasicAuthMethod struct {
	queryBus query.Bus
}

func NewBasicAuthMethod(queryBus query.Bus) *BasicAuthMethod {
	return &BasicAuthMethod{queryBus: queryBus}
}

func (m *BasicAuthMethod) CanProcess(c *gin.Context) (bool, error) {
	matches := strings.HasPrefix(c.GetHeader("Authorization"), "Basic ")
	if matches && c.GetBool(authenticatedKey) {
		return false, app_errors.BadRequestError(ErrAlreadyAuthenticated)
	}
	return matches, nil
}

func (m *BasicAuthMethod) Authenticate(c *gin.Context) (string, error) {
	encoded := strings.TrimPrefix(c.GetHeader("Authorization"), "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", app_errors.UnauthorizedError(ErrInvalidBasicAuth).
			WithDescription("invalid Basic Auth encoding")
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", app_errors.UnauthorizedError(ErrInvalidBasicAuth).
			WithDescription("invalid Basic Auth format: expected client_id:client_secret")
	}

	clientID, secret := parts[0], parts[1]
	if _, err = m.queryBus.Dispatch(c.Request.Context(), verifySecret.NewQuery(clientID, secret)); err != nil {
		return "", err
	}
	return clientID, nil
}
