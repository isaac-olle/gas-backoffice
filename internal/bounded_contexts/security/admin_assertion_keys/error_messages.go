package admin_assertion_keys

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrAdminAssertionKeyNotFound    app_errors.ErrorCode = "ADMIN_ASSERTION_KEY_NOT_FOUND"
	ErrAdminAssertionKeyNotCreated  app_errors.ErrorCode = "ADMIN_ASSERTION_KEY_NOT_CREATED"
	ErrAdminAssertionKeyNotSaved    app_errors.ErrorCode = "ADMIN_ASSERTION_KEY_NOT_SAVED"
	ErrInvalidAdminAssertion        app_errors.ErrorCode = "INVALID_ADMIN_ASSERTION"
	ErrInvalidOwnership             app_errors.ErrorCode = "ADMIN_ASSERTION_KEY_INVALID_OWNERSHIP"
	ErrAdminAssertionKeyQueryFailed app_errors.ErrorCode = "ADMIN_ASSERTION_KEY_QUERY_FAILED"
)
