package user_assertion_keys

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrUserAssertionKeyNotFound    app_errors.ErrorCode = "USER_ASSERTION_KEY_NOT_FOUND"
	ErrUserAssertionKeyNotCreated  app_errors.ErrorCode = "USER_ASSERTION_KEY_NOT_CREATED"
	ErrUserAssertionKeyNotSaved    app_errors.ErrorCode = "USER_ASSERTION_KEY_NOT_SAVED"
	ErrInvalidUserAssertion        app_errors.ErrorCode = "INVALID_USER_ASSERTION"
	ErrInvalidOwnership            app_errors.ErrorCode = "USER_ASSERTION_KEY_INVALID_OWNERSHIP"
	ErrUserAssertionKeyQueryFailed app_errors.ErrorCode = "USER_ASSERTION_KEY_QUERY_FAILED"
)
