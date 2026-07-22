package check_tenant_access

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const QueryType query.Type = "query.backoffice.user.check_tenant_access"

type Query struct {
	tenantId ids.UserId
	author   *authors.UserAuthor
}

func NewQuery(tenantId ids.UserId, author *authors.UserAuthor) *Query {
	return &Query{tenantId: tenantId, author: author}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct {
	service *Service
}

func NewQueryHandler(service *Service) *QueryHandler {
	return &QueryHandler{service: service}
}

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.CheckTenantAccess)
}
