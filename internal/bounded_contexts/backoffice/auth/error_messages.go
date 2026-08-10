package auth

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrSessionNotFound   app_errors.ErrorCode = "SESSION_NOT_FOUND"
	ErrSessionNotCreated app_errors.ErrorCode = "SESSION_NOT_CREATED"
	ErrSessionNotSaved   app_errors.ErrorCode = "SESSION_NOT_SAVED"
)
