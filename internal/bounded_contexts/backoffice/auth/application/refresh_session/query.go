package refresh_session

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

const QueryType query.Type = "query.backoffice.auth.refresh_session"

type Query struct{ refreshToken string }

func NewQuery(refreshToken string) *Query { return &Query{refreshToken: refreshToken} }
func (q *Query) Type() query.Type         { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler { return &QueryHandler{service: service} }

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.Refresh)
}
