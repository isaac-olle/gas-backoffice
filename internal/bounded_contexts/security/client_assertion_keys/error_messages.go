package client_assertion_keys

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrClientAssertionKeyNotFound   app_errors.ErrorCode = "CLIENT_ASSERTION_KEY_NOT_FOUND"
	ErrClientAssertionKeyNotCreated app_errors.ErrorCode = "CLIENT_ASSERTION_KEY_NOT_CREATED"
	ErrClientAssertionKeyNotSaved   app_errors.ErrorCode = "CLIENT_ASSERTION_KEY_NOT_SAVED"
	ErrInvalidOwnership             app_errors.ErrorCode = "CLIENT_ASSERTION_KEY_INVALID_OWNERSHIP"
	ErrInvalidClientAssertion       app_errors.ErrorCode = "CLIENT_ASSERTION_INVALID"
)
