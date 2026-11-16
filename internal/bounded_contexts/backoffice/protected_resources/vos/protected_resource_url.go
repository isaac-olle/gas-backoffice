package vos

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

type ProtectedResourceURL struct{ vo.Url }

func NewProtectedResourceURL(raw string) (ProtectedResourceURL, error) {
	if raw == "" {
		return ProtectedResourceURL{}, nil
	}
	value, err := vo.NewUrl(raw, vo.NewFragmentNotAllowedOption())
	if err != nil {
		return ProtectedResourceURL{}, app_errors.BadRequestError("INVALID_PROTECTED_RESOURCE_URL").
			WithDescription("protected resource URL must be an absolute HTTP or HTTPS URL").
			WithCause(err, false)
	}
	return ProtectedResourceURL{Url: value}, nil
}
