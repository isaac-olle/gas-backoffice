package users

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrUserNotFound               app_errors.ErrorCode = "USER_NOT_FOUND"
	ErrUserNotCreated             app_errors.ErrorCode = "USER_NOT_CREATED"
	ErrUserNotSaved               app_errors.ErrorCode = "USER_NOT_SAVED"
	ErrUserNameNotChanged         app_errors.ErrorCode = "USER_NAME_NOT_CHANGED"
	ErrInvalidOwnership           app_errors.ErrorCode = "USER_INVALID_OWNERSHIP"
	ErrSystemResource             app_errors.ErrorCode = "USER_IS_SYSTEM_RESOURCE"
	ErrInvalidEmail               app_errors.ErrorCode = "USER_INVALID_EMAIL"
	ErrUsersNotListed             app_errors.ErrorCode = "TENANTS_NOT_LISTED"
	ErrUserConcurrentModification app_errors.ErrorCode = "USER_CONCURRENT_MODIFICATION"
)
