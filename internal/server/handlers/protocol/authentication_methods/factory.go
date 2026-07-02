package authentication_methods

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
)

var ErrNoAuthenticationMethod app_errors.ErrorCode = "NO_AUTHENTICATION_METHOD"

type ClientAuthMethodFactory struct {
	methods []AuthenticationMethod
}

func NewClientAuthMethodFactory(methods ...AuthenticationMethod) *ClientAuthMethodFactory {
	return &ClientAuthMethodFactory{methods: methods}
}

func (f *ClientAuthMethodFactory) Process(c *gin.Context) (string, error) {
	for _, m := range f.methods {
		ok, err := m.CanProcess(c)
		if err != nil {
			return "", err
		}
		if !ok {
			continue
		}
		clientID, err := m.Authenticate(c)
		if err != nil {
			return "", err
		}
		MarkAuthenticated(c)
		return clientID, nil
	}
	return "", app_errors.UnauthorizedError(ErrNoAuthenticationMethod).
		WithDescription("no supported authentication method found in request")
}
