package check_ownership

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

const QueryType query.Type = "query.backoffice.protected_resource.check_ownership"

// Query checks that the protected resource identified by protectedResourceId
// belongs to the caller's active subscription (read from context). It
// mutates no state; it exists purely so handlers in other modules can verify
// ownership of a referenced protected resource without importing the
// protected_resources domain directly.
type Query struct {
	protectedResourceId string
}

func NewQuery(protectedResourceId string) *Query {
	return &Query{protectedResourceId: protectedResourceId}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler {
	return &QueryHandler{service: service}
}

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.Check)
}
