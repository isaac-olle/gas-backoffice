package list

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
)

const QueryType query.Type = "query.backoffice.client.list"

// Query lists clients visible to the caller's allowed subscriptions.
type Query struct {
	author *authors.UserAuthor
}

func NewQuery(author *authors.UserAuthor) *Query {
	return &Query{author: author}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler {
	return &QueryHandler{service: service}
}

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.List)
}
