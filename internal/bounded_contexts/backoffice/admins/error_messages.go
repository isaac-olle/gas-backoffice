package admins

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrAdminNotFound               app_errors.ErrorCode = "ADMIN_NOT_FOUND"
	ErrAdminNotCreated             app_errors.ErrorCode = "ADMIN_NOT_CREATED"
	ErrAdminNotSaved               app_errors.ErrorCode = "ADMIN_NOT_SAVED"
	ErrAdminNameNotChanged         app_errors.ErrorCode = "ADMIN_NAME_NOT_CHANGED"
	ErrInvalidOwnership            app_errors.ErrorCode = "ADMIN_INVALID_OWNERSHIP"
	ErrInvalidEmail                app_errors.ErrorCode = "ADMIN_INVALID_EMAIL"
	ErrAdminConcurrentModification app_errors.ErrorCode = "ADMIN_CONCURRENT_MODIFICATION"
)
