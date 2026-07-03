package tenants

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrTenantNotFound       app_errors.ErrorCode = "TENANT_NOT_FOUND"
	ErrTenantNotCreated     app_errors.ErrorCode = "TENANT_NOT_CREATED"
	ErrTenantNotSaved       app_errors.ErrorCode = "TENANT_NOT_SAVED"
	ErrTenantNameNotChanged app_errors.ErrorCode = "TENANT_NAME_NOT_CHANGED"
	ErrInvalidOwnership     app_errors.ErrorCode = "TENANT_INVALID_OWNERSHIP"
	ErrSystemResource       app_errors.ErrorCode = "TENANT_IS_SYSTEM_RESOURCE"
)
