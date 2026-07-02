package query

import "context"

type Handler interface {
	Handle(context.Context, Query) (any, error)
}
