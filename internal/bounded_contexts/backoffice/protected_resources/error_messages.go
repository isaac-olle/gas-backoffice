package protected_resources

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrProtectedResourceNotFound   app_errors.ErrorCode = "PROTECTED_RESOURCE_NOT_FOUND"
	ErrProtectedResourceNotCreated app_errors.ErrorCode = "PROTECTED_RESOURCE_NOT_CREATED"
	ErrProtectedResourceNotSaved   app_errors.ErrorCode = "PROTECTED_RESOURCE_NOT_SAVED"
	ErrProtectedResourceNotUpdated app_errors.ErrorCode = "PROTECTED_RESOURCE_NOT_UPDATED"
	ErrInvalidOwnership            app_errors.ErrorCode = "PROTECTED_RESOURCE_INVALID_OWNERSHIP"
	ErrInvalidUserId               app_errors.ErrorCode = "PROTECTED_RESOURCE_INVALID_USER_ID"
	ErrSystemResource              app_errors.ErrorCode = "PROTECTED_RESOURCE_IS_SYSTEM_RESOURCE"
)
