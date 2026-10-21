package jwks

import "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"

var (
	ErrCannotGetSigningKeys app_errors.ErrorCode = "CANNOT_GET_SIGNING_KEYS"
	ErrInvalidSigningKey    app_errors.ErrorCode = "INVALID_SIGNING_KEY"
)
