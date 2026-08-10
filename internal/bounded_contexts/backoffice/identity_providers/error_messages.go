package identity_providers

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrIdentityProviderNotFound      app_errors.ErrorCode = "IDENTITY_PROVIDER_NOT_FOUND"
	ErrIdentityProviderDisabled      app_errors.ErrorCode = "IDENTITY_PROVIDER_DISABLED"
	ErrIdentityProviderMisconfigured app_errors.ErrorCode = "IDENTITY_PROVIDER_MISCONFIGURED"
)
