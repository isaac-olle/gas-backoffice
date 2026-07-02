package generate_jwt

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
	"time"
)

const QueryType query.Type = "query.security.key_pair.generate_jwt"

type Query struct {
	keyPairId string
	claims    map[string]any
	expiresIn time.Duration
}

func NewQuery(keyPairId string, claims map[string]any, expiresIn time.Duration) *Query {
	return &Query{keyPairId: keyPairId, claims: claims, expiresIn: expiresIn}
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler { return &QueryHandler{service: service} }

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.GenerateJWT)
}
