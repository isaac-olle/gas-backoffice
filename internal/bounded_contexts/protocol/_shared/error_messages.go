package shared

import "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"

var (
	ErrCannotGetIdentityProvider app_errors.ErrorCode = "ERROR_CANNOT_GET_IDENTITY_PROVIDER"
	ErrIdentityProviderNotFound  app_errors.ErrorCode = "IDENTITY_PROVIDER_NOT_FOUND"
)
