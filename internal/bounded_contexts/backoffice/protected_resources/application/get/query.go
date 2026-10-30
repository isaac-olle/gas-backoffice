package get

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const QueryType query.Type = "query.backoffice.protected_resource.get"

type Query struct {
	protectedResourceId ids.ProtectedResourceId
	author              *vo.UserAuthor
}

func NewQuery(protectedResourceId ids.ProtectedResourceId, author *vo.UserAuthor) *Query {
	return &Query{protectedResourceId: protectedResourceId, author: author}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler {
	return &QueryHandler{service: service}
}

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.Get)
}
