package build_authorization_response

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const QueryType query.Type = "query.protocol.session.build_authorization_response"

type Query struct {
	sessionID ids.SessionId
}

func NewQuery(sessionID ids.SessionId) *Query {
	return &Query{sessionID: sessionID}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler { return &QueryHandler{service: service} }

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.BuildAuthorizationResponse)
}
