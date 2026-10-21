package list

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const QueryType query.Type = "query.backoffice.protected_resource.list"

// Query lists the protected resources belonging to a single subscription
// (the one in the request path).
type Query struct {
	subscriptionId ids.SubscriptionId
}

func NewQuery(subscriptionId ids.SubscriptionId) *Query {
	return &Query{subscriptionId: subscriptionId}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler {
	return &QueryHandler{service: service}
}

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.List)
}
