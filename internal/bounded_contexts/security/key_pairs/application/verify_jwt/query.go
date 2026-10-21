package verify_jwt

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

const QueryType query.Type = "query.security.key_pair.verify_jwt"

type Query struct {
	token            string
	expectedClaims   map[string]any
	claimsToReturn   map[string]any
	verifyExpiration bool
}

func NewQuery(token string) *Query {
	return &Query{
		token:          token,
		expectedClaims: make(map[string]any),
	}
}

func (q *Query) WithExpectedClaim(name string, value any) *Query {
	q.expectedClaims[name] = value
	return q
}

func (q *Query) WithVerifyExpiration() *Query {
	q.verifyExpiration = true
	return q
}

// WithClaimsToReturn declares the shape of what the caller needs. A nil
// value means "return this claim as-is"; a map[string]any value means
// "descend into this claim (which must itself be an object) and apply the
// same rule recursively".
//
// Ús (pla):
//
//	WithClaimsToReturn(map[string]any{
//		"sub":   nil,
//		"role":  nil,
//		"scope": nil,
//	})
//
// Ús (amb un claim niat, per exemple "address"):
//
//	WithClaimsToReturn(map[string]any{
//		"sub": nil,
//		"address": map[string]any{
//			"city": nil,
//		},
//	})
func (q *Query) WithClaimsToReturn(template map[string]any) *Query {
	q.claimsToReturn = template
	return q
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct {
	service *Service
}

func NewQueryHandler(service *Service) *QueryHandler {
	return &QueryHandler{service: service}
}

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.VerifyJWT)
}
