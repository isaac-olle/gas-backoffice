package clients

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrClientDuplicateName                     app_errors.ErrorCode = "CLIENT_DUPLICATE_NAME"
	ErrClientNotFound                          app_errors.ErrorCode = "CLIENT_NOT_FOUND"
	ErrClientNotCreated                        app_errors.ErrorCode = "CLIENT_NOT_CREATED"
	ErrClientNotSaved                          app_errors.ErrorCode = "CLIENT_NOT_SAVED"
	ErrInvalidOwnership                        app_errors.ErrorCode = "CLIENT_INVALID_OWNERSHIP"
	ErrPublicClientCannotHaveClientCredentials app_errors.ErrorCode = "PUBLIC_CLIENT_CANNOT_HAVE_CLIENT_CREDENTIALS"
	ErrSystemResource                          app_errors.ErrorCode = "CLIENT_IS_SYSTEM_RESOURCE"
	ErrClientsNotListed                        app_errors.ErrorCode = "CLIENTS_NOT_LISTED"
	ErrClientConcurrentModification            app_errors.ErrorCode = "CLIENT_CONCURRENT_MODIFICATION"
)
