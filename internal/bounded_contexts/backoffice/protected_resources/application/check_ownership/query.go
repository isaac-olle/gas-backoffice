package check_ownership

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

const QueryType query.Type = "query.backoffice.protected_resource.check_ownership"

// Query checks that the protected resource identified by protectedResourceId
// belongs to one of the caller's allowed subscriptions. It mutates no state;
// it exists so middleware and handlers in other modules can verify ownership
// without importing the protected_resources domain directly.
type Query struct {
	protectedResourceId  string
	allowedSubscriptions []string
	allowAll             bool
}

func NewQuery(protectedResourceId string) *Query {
	return &Query{protectedResourceId: protectedResourceId}
}

func NewQueryForSubscriptions(protectedResourceId string, allowedSubscriptions []string, allowAll bool) *Query {
	return &Query{protectedResourceId: protectedResourceId, allowedSubscriptions: allowedSubscriptions, allowAll: allowAll}
}

// NewQueryForResource resolves the resource ownership without authorizing the
// caller. Authorization is performed by SubscriptionMiddleware.
func NewQueryForResource(protectedResourceId string) *Query {
	return &Query{protectedResourceId: protectedResourceId, allowAll: true}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler {
	return &QueryHandler{service: service}
}

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.Check)
}
