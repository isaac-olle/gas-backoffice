package admins

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrAdminNotFound    app_errors.ErrorCode = "ADMIN_NOT_FOUND"
	ErrAdminNotCreated  app_errors.ErrorCode = "ADMIN_NOT_CREATED"
	ErrAdminNotSaved    app_errors.ErrorCode = "ADMIN_NOT_SAVED"
	ErrInvalidOwnership app_errors.ErrorCode = "ADMIN_INVALID_OWNERSHIP"
	ErrSystemResource   app_errors.ErrorCode = "ADMIN_IS_SYSTEM_RESOURCE"
)
