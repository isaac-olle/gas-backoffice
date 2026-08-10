package vos

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var ErrInvalidFlow app_errors.ErrorCode = "INVALID_FLOW"

type Flow string

const (
	AuthorizationCodeFlow Flow = "authorization_code"
	ClientCredentialsFlow Flow = "client_credentials"
)

func NewFlowFromString(raw string) (Flow, error) {
	switch Flow(raw) {
	case AuthorizationCodeFlow, ClientCredentialsFlow:
		return Flow(raw), nil
	default:
		return "", app_errors.BadRequestError(ErrInvalidFlow).
			WithDescription("invalid flow: " + raw)
	}
}

func (f Flow) String() string { return string(f) }
