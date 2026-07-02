package clients

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
)

var (
	ErrClientNotFound   app_errors.ErrorCode = "CLIENT_NOT_FOUND"
	ErrClientNotCreated app_errors.ErrorCode = "CLIENT_NOT_CREATED"
	ErrClientNotSaved   app_errors.ErrorCode = "CLIENT_NOT_SAVED"
	ErrInvalidOwnership app_errors.ErrorCode = "CLIENT_INVALID_OWNERSHIP"
)
