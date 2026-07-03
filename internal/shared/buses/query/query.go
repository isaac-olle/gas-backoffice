package query

import (
	"context"
	"fmt"
	errorDomain "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"reflect"
)

var (
	ErrInvalidContext errorDomain.ErrorCode = "INVALID_CONTEXT"
	ErrInvalidQuery   errorDomain.ErrorCode = "INVALID_QUERY"
)

type Type string

func (t Type) String() string { return string(t) }

type Query interface {
	Type() Type
}

func Process[T Query, C context.Context](ctx context.Context, qry Query, handler func(C, T) (any, error)) (any, error) {
	cCtx, ok := ctx.(C)
	if !ok {
		return nil, errorDomain.InternalError(ErrInvalidContext).WithDescription("unexpected context type")
	}
	parsedQuery, ok := qry.(T)
	if !ok {
		expectedType := reflect.TypeOf((*T)(nil)).Elem()
		return nil, errorDomain.InternalError(ErrInvalidQuery).
			WithDescription(fmt.Sprintf("unexpected query type. Expected: %s Got: %s", expectedType.String(), qry.Type().String()))
	}

	return handler(cCtx, parsedQuery)
}
