package generate_jwt_by_pr_id

import (
	"context"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

const QueryType query.Type = "query.security.key_pair.generate_jwt_by_pr_id"

type Query struct {
	protectedResourceId string
	claims              map[string]any
	expiresIn           time.Duration
}

func NewQuery(protectedResourceId string, claims map[string]any, expiresIn time.Duration) *Query {
	return &Query{protectedResourceId: protectedResourceId, claims: claims, expiresIn: expiresIn}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler { return &QueryHandler{service: service} }

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.GenerateJWT)
}
