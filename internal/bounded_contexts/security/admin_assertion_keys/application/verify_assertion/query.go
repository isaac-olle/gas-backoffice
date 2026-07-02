package verify_assertion

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
)

const QueryType query.Type = "query.security.admin_assertion_key.verify_assertion"

type Query struct {
	assertion string
	grantType string
}

func NewQuery(assertion, grantType string) *Query {
	return &Query{assertion: assertion, grantType: grantType}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler { return &QueryHandler{service: service} }

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.VerifyAssertion)
}
