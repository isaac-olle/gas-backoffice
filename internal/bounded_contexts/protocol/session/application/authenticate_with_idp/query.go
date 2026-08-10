package authenticate_with_idp

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const QueryType query.Type = "query.protocol.session.authenticate_with_idp"

type Query struct {
	sessionID ids.SessionId
	idpID     string
}

func NewQuery(sessionID ids.SessionId, idpID string) *Query {
	return &Query{sessionID: sessionID, idpID: idpID}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler {
	return &QueryHandler{service: service}
}

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.Authenticate)
}
