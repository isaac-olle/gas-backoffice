package scopes

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrScopeNotFound    app_errors.ErrorCode = "SCOPE_NOT_FOUND"
	ErrScopeNotCreated  app_errors.ErrorCode = "SCOPE_NOT_CREATED"
	ErrScopeNotSaved    app_errors.ErrorCode = "SCOPE_NOT_SAVED"
	ErrScopeNotUpdated  app_errors.ErrorCode = "SCOPE_NOT_UPDATED"
	ErrInvalidOwnership app_errors.ErrorCode = "SCOPE_INVALID_OWNERSHIP"
	ErrSystemResource   app_errors.ErrorCode = "SCOPE_IS_SYSTEM_RESOURCE"
)
