package check_ownership

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const QueryType query.Type = "query.backoffice.user.check_ownership"

type Query struct {
	userId  string
	adminId ids.AdminId
}

func NewQuery(userId string, adminId ids.AdminId) *Query {
	return &Query{userId: userId, adminId: adminId}
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
