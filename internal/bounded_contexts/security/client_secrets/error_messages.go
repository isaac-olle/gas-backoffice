package client_secrets

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrClientSecretNotFound   app_errors.ErrorCode = "CLIENT_SECRET_NOT_FOUND"
	ErrClientSecretNotCreated app_errors.ErrorCode = "CLIENT_SECRET_NOT_CREATED"
	ErrClientSecretNotSaved   app_errors.ErrorCode = "CLIENT_SECRET_NOT_SAVED"
	ErrInvalidOwnership       app_errors.ErrorCode = "CLIENT_SECRET_INVALID_OWNERSHIP"
	ErrInvalidClientSecret    app_errors.ErrorCode = "CLIENT_SECRET_INVALID"
)
