package vos

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
)

var ErrInvalidFlow app_errors.ErrorCode = "INVALID_FLOW"

const (
	AuthorizationCodeFlow = "authorization_code"
	ClientCredentialsFlow = "client_credentials"
)

type Flow struct{ value string }

func NewAuthorizationCodeFlow() Flow {
	return Flow{value: AuthorizationCodeFlow}
}

func NewClientCredentialsFlow() Flow {
	return Flow{value: ClientCredentialsFlow}
}

func DetermineFlow(responseType string) (Flow, error) {
	switch responseType {
	case "code":
		return NewAuthorizationCodeFlow(), nil
	case "token", "id_token":
		return NewClientCredentialsFlow(), nil
	default:
		return Flow{}, app_errors.BadRequestError(ErrInvalidFlow).
			WithDescription("cannot determine flow from response_type: " + responseType)
	}
}

func NewFlowFromString(raw string) (Flow, error) {
	switch raw {
	case AuthorizationCodeFlow, ClientCredentialsFlow:
		return Flow{value: raw}, nil
	default:
		return Flow{}, app_errors.BadRequestError(ErrInvalidFlow).
			WithDescription("invalid flow: " + raw)
	}
}

func (f Flow) String() string { return f.value }
