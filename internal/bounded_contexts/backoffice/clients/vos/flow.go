package vos

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var ErrInvalidFlow app_errors.ErrorCode = "INVALID_FLOW"

const (
	AuthorizationCodeFlow = "authorization_code"
	PKCEFlow              = "pkce"
	ClientCredentialsFlow = "client_credentials"
)

type Flow struct {
	value string
}

func NewFlow(raw string) (Flow, error) {
	switch raw {
	case AuthorizationCodeFlow, PKCEFlow, ClientCredentialsFlow:
		return Flow{value: raw}, nil
	default:
		return Flow{}, app_errors.BadRequestError(ErrInvalidFlow).
			WithDescription("invalid flow, must be one of: authorization_code, pkce, client_credentials")
	}
}

func (f Flow) String() string { return f.value }
