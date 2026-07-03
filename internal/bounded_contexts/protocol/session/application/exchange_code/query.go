package exchange_code

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

const QueryType query.Type = "query.protocol.session.exchange_code"

type Query struct {
	code         string
	clientID     string
	redirectURI  string
	codeVerifier string
}

func NewQuery(code, clientID, redirectURI, codeVerifier string) *Query {
	return &Query{
		code:         code,
		clientID:     clientID,
		redirectURI:  redirectURI,
		codeVerifier: codeVerifier,
	}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler { return &QueryHandler{service: service} }

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.ExchangeCode)
}
