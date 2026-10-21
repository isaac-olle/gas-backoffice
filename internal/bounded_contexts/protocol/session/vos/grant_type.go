package vos

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var ErrInvalidGrantType app_errors.ErrorCode = "INVALID_GRANT_TYPE"

type GrantType string

const (
	AuthorizationCodeGrantType GrantType = "authorization_code"
	ClientCredentialsGrantType GrantType = "client_credentials"
)

func NewGrantTypeFromString(raw string) (GrantType, error) {
	switch GrantType(raw) {
	case AuthorizationCodeGrantType, ClientCredentialsGrantType:
		return GrantType(raw), nil
	default:
		return "", app_errors.BadRequestError(ErrInvalidGrantType).
			WithDescription("invalid grant type: " + raw)
	}
}

func (g GrantType) String() string { return string(g) }
