package tenants

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrUserNotFound               app_errors.ErrorCode = "TENANT_NOT_FOUND"
	ErrUserNotCreated             app_errors.ErrorCode = "TENANT_NOT_CREATED"
	ErrUserNotSaved               app_errors.ErrorCode = "TENANT_NOT_SAVED"
	ErrUserNameNotChanged         app_errors.ErrorCode = "TENANT_NAME_NOT_CHANGED"
	ErrInvalidOwnership           app_errors.ErrorCode = "TENANT_INVALID_OWNERSHIP"
	ErrSystemResource             app_errors.ErrorCode = "TENANT_IS_SYSTEM_RESOURCE"
	ErrInvalidEmail               app_errors.ErrorCode = "TENANT_INVALID_EMAIL"
	ErrUsersNotListed             app_errors.ErrorCode = "TENANTS_NOT_LISTED"
	ErrUserConcurrentModification app_errors.ErrorCode = "TENANT_CONCURRENT_MODIFICATION"
)
