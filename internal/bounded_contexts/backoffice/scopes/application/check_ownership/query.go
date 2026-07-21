package check_ownership

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

const QueryType query.Type = "query.backoffice.scope.check_ownership"

type Query struct {
	scopeId             string
	protectedResourceId string
	userId              string
}

func NewQuery(scopeId, protectedResourceId, userId string) *Query {
	return &Query{scopeId: scopeId, protectedResourceId: protectedResourceId, userId: userId}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct {
	service *Service
}

func NewQueryHandler(service *Service) *QueryHandler {
	return &QueryHandler{service: service}
}

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.CheckOwnership)
}
