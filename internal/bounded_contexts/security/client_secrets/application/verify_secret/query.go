package verify_secret

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

const QueryType query.Type = "query.security.client_secret.verify_secret"

type Query struct {
	clientID    string
	plainSecret string
}

func NewQuery(clientID, plainSecret string) *Query {
	return &Query{clientID: clientID, plainSecret: plainSecret}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler { return &QueryHandler{service: service} }

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.VerifySecret)
}
