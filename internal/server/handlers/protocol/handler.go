package protocol

import (
	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
)

func unsupportedGrantTypeError(grantType string) error {
	return app_errors.BadRequestError(session.ErrUnsupportedGrantType).
		WithDescription("unsupported grant_type: " + grantType)
}
