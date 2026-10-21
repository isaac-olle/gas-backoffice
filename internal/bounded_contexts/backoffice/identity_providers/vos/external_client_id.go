package vos

import "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"

var ErrInvalidExternalClientID app_errors.ErrorCode = "INVALID_EXTERNAL_CLIENT_ID"

type ExternalClientID struct{ value string }

func NewExternalClientID(raw string) (ExternalClientID, error) {
	if raw == "" {
		return ExternalClientID{}, app_errors.BadRequestError(ErrInvalidExternalClientID).
			WithDescription("client_id cannot be empty")
	}
	return ExternalClientID{value: raw}, nil
}

func (e ExternalClientID) String() string { return e.value }
