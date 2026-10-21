package vos

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var ErrInvalidGrantType app_errors.ErrorCode = "INVALID_GRANT_TYPE"

const (
	AuthorizationCodeGrantType = "authorization_code"
	ClientCredentialsGrantType = "client_credentials"
)

type GrantType struct {
	value string
}

func NewGrantType(raw string) (GrantType, error) {
	switch raw {
	case AuthorizationCodeGrantType, ClientCredentialsGrantType:
		return GrantType{value: raw}, nil
	default:
		return GrantType{}, app_errors.BadRequestError(ErrInvalidGrantType).
			WithDescription("invalid grant type, must be one of: authorization_code, client_credentials")
	}
}

func (g GrantType) String() string { return g.value }
