package session

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var (
	ErrInvalidSessionParameters     app_errors.ErrorCode = "INVALID_SESSION_PARAMETERS"
	ErrCannotSaveSession            app_errors.ErrorCode = "ERROR_CANNOT_SAVE_SESSION"
	ErrSessionNotFound              app_errors.ErrorCode = "SESSION_NOT_FOUND"
	ErrSessionConsumed              app_errors.ErrorCode = "SESSION_ALREADY_CONSUMED"
	ErrCannotGetClient              app_errors.ErrorCode = "ERROR_CANNOT_GET_CLIENT"
	ErrClientNotFound               app_errors.ErrorCode = "CLIENT_NOT_FOUND"
	ErrInvalidRedirectURI           app_errors.ErrorCode = "INVALID_REDIRECT_URI"
	ErrInvalidScope                 app_errors.ErrorCode = "INVALID_SCOPE"
	ErrGrantTypeNotAllowed          app_errors.ErrorCode = "GRANT_TYPE_NOT_ALLOWED"
	ErrPkceRequired                 app_errors.ErrorCode = "PKCE_REQUIRED"
	ErrCodeNotFound                 app_errors.ErrorCode = "CODE_NOT_FOUND"
	ErrInvalidCodeVerifier          app_errors.ErrorCode = "INVALID_CODE_VERIFIER"
	ErrInvalidClientSecret          app_errors.ErrorCode = "INVALID_CLIENT_SECRET"
	ErrClientSecretNotFound         app_errors.ErrorCode = "CLIENT_SECRET_NOT_FOUND"
	ErrUnsupportedGrantType         app_errors.ErrorCode = "UNSUPPORTED_GRANT_TYPE"
	ErrCannotSignToken              app_errors.ErrorCode = "CANNOT_SIGN_TOKEN"
	ErrAmbiguousScope               app_errors.ErrorCode = "AMBIGUOUS_SCOPE"
	ErrScopesFromMultipleResources  app_errors.ErrorCode = "SCOPES_FROM_MULTIPLE_RESOURCES"
	ErrClientNotPublic              app_errors.ErrorCode = "CLIENT_NOT_PUBLIC"
	ErrAuthenticatedUserNotResolved app_errors.ErrorCode = "AUTHENTICATED_USER_NOT_RESOLVED"
	ErrSessionNotConsumed           app_errors.ErrorCode = "SESSION_NOT_CONSUMED"
	ErrSessionNotAuthorizing        app_errors.ErrorCode = "SESSION_NOT_AUTHORIZING"
	ErrInvalidIdpState              app_errors.ErrorCode = "INVALID_IDP_STATE"
	ErrCannotAuthenticateWithIdp    app_errors.ErrorCode = "CANNOT_AUTHENTICATE_WITH_IDP"
)
