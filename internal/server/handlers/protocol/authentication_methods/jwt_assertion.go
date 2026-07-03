package authentication_methods

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"

	verifyAssertion "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/verify_assertion"
)

const jwtBearerAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"

var ErrInvalidAssertion app_errors.ErrorCode = "INVALID_CLIENT_ASSERTION"

type JWTAssertionMethod struct {
	queryBus query.Bus
}

func NewJWTAssertionMethod(queryBus query.Bus) *JWTAssertionMethod {
	return &JWTAssertionMethod{queryBus: queryBus}
}

func (m *JWTAssertionMethod) CanProcess(c *gin.Context) (bool, error) {
	matches := c.PostForm("client_assertion_type") == jwtBearerAssertionType
	if matches && c.GetBool(authenticatedKey) {
		return false, app_errors.BadRequestError(ErrAlreadyAuthenticated)
	}
	return matches, nil
}

func (m *JWTAssertionMethod) Authenticate(c *gin.Context) (string, error) {
	assertion := c.PostForm("client_assertion")
	if assertion == "" {
		return "", app_errors.UnauthorizedError(ErrInvalidAssertion).
			WithDescription("client_assertion is required when client_assertion_type is jwt-bearer")
	}

	result, err := m.queryBus.Dispatch(c.Request.Context(), verifyAssertion.NewQuery(assertion))
	if err != nil {
		return "", err
	}

	claims, ok := result.(*verifyAssertion.ClientClaims)
	if !ok || claims.ClientID == "" {
		return "", app_errors.UnauthorizedError(ErrInvalidAssertion).
			WithDescription("invalid client assertion claims")
	}
	return claims.ClientID, nil
}
