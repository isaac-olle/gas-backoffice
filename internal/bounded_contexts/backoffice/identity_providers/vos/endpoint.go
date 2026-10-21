package vos

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

var ErrInvalidEndpoint app_errors.ErrorCode = "INVALID_ENDPOINT"

type Endpoint struct{ vo.Url }

func NewEndpoint(raw string) (Endpoint, error) {
	u, err := vo.NewUrl(raw)
	if err != nil {
		return Endpoint{}, app_errors.BadRequestError(ErrInvalidEndpoint).
			WithDescription("invalid endpoint URL").
			WithCause(err, true)
	}
	return Endpoint{u}, nil
}

func (e Endpoint) String() string { return e.Url.String() }
