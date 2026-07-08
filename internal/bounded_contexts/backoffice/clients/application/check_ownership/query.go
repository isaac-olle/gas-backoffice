package check_ownership

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const QueryType query.Type = "query.backoffice.client.check_ownership"

type Query struct {
	clientId string
	adminId  ids2.AdminId
	userId   ids2.UserId
}

func NewQuery(clientId string, adminId ids2.AdminId, userId ids2.UserId) *Query {
	return &Query{clientId: clientId, adminId: adminId, userId: userId}
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
