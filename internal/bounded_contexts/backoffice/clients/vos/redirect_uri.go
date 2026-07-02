package vos

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo"
)

var ErrInvalidRedirectUri app_errors.ErrorCode = "INVALID_REDIRECT_URI"

type RedirectUri struct{ vo.Url }

func NewRedirectUri(raw string) (RedirectUri, error) {
	u, err := vo.NewUrl(raw, vo.NewFragmentNotAllowedOption())
	if err != nil {
		return RedirectUri{}, app_errors.BadRequestError(ErrInvalidRedirectUri).
			WithDescription("invalid redirect uri").
			WithCause(err)
	}
	return RedirectUri{u}, nil
}

func (r RedirectUri) String() string { return r.Url.String() }
