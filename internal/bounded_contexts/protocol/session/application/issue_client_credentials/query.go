package issue_client_credentials

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
)

const QueryType query.Type = "query.protocol.session.issue_client_credentials"

type Query struct {
	clientID string
	scope    string
}

func NewQuery(clientID, scope string) *Query {
	return &Query{clientID: clientID, scope: scope}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler { return &QueryHandler{service: service} }

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.IssueClientCredentials)
}
