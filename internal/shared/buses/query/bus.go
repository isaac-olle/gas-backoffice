package query

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

var ErrInvalidResponse app_errors.ErrorCode = "INVALID_RESPONSE"

type Bus interface {
	Dispatch(context.Context, Query) (any, error)
	Register(Type, Handler)
}

func DispatchTyped[R any](ctx context.Context, qb Bus, qry Query) (R, error) {
	res, err := qb.Dispatch(ctx, qry)
	if err != nil {
		return *new(R), err
	}
	typedRes, ok := res.(R)
	if !ok {
		return *new(R), app_errors.InternalError(ErrInvalidResponse).
			WithDescription("invalid response type from typed dispatch")
	}
	return typedRes, nil
}
