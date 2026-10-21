package vos

import "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"

var ErrInvalidName app_errors.ErrorCode = "INVALID_IDENTITY_PROVIDER_NAME"

type Name struct{ value string }

func NewName(raw string) (Name, error) {
	if raw == "" {
		return Name{}, app_errors.BadRequestError(ErrInvalidName).
			WithDescription("name cannot be empty")
	}
	return Name{value: raw}, nil
}

func (n Name) String() string { return n.value }
