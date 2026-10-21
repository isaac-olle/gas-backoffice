package auth

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrLoginRequestNotFound     app_errors.ErrorCode = "SESSION_NOT_FOUND"
	ErrLoginRequestNotCreated   app_errors.ErrorCode = "SESSION_NOT_CREATED"
	ErrLoginRequestNotSaved     app_errors.ErrorCode = "SESSION_NOT_SAVED"
	ErrInitiateLoginFailed      app_errors.ErrorCode = "INITIATE_LOGIN_FAILED"
	ErrCannotExchangeToken      app_errors.ErrorCode = "CANNOT_EXCHANGE_TOKEN"
	ErrCannotCreateLoginRequest app_errors.ErrorCode = "CANNOT_CREATE_LOGIN_REQUEST"
	ErrInvalidIdpReferenceName  app_errors.ErrorCode = "INVALID_IDP_REFERENCE_NAME"
	ErrCannotCreateLoginSession app_errors.ErrorCode = "CANNOT_CREATE_LOGIN_SESSION"
)
