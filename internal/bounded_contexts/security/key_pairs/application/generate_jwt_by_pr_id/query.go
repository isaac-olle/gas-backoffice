package generate_jwt_by_pr_id

import (
	"context"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

const QueryType query.Type = "query.security.key_pair.generate_jwt_by_pr_id"

// Query asks security to sign a JWT for a protected resource. Only claims
// the caller alone can know (sub, aud, jti, and any domain-specific extra
// claim) are accepted here; iss, iat, exp, kid and alg are always decided by
// security itself, since they depend on the key pair it resolves internally.
type Query struct {
	protectedResourceId string
	expiresIn           time.Duration
	spec                key_pairs.TokenSpec
}

func NewQuery(protectedResourceId, sub string, expiresIn time.Duration) *Query {
	return &Query{
		protectedResourceId: protectedResourceId,
		expiresIn:           expiresIn,
		spec: key_pairs.TokenSpec{
			Sub:         sub,
			ExtraClaims: make(map[string]any),
		},
	}
}

// WithAudience sets the "aud" claim.
func (q *Query) WithAudience(aud string) *Query {
	q.spec.Aud = aud
	return q
}

// WithJTI sets the "jti" claim. The caller picks the value because it's the
// one that will need to reference this specific token later (e.g. for
// revocation or audit correlation).
func (q *Query) WithJTI(jti string) *Query {
	q.spec.JTI = jti
	return q
}

// WithExtraClaim adds a domain-specific claim (e.g. "scope", "email", "role")
// that security has no business interpreting.
func (q *Query) WithExtraClaim(name string, value any) *Query {
	q.spec.ExtraClaims[name] = value
	return q
}

// WithX5c asks security to embed the key pair's certificate chain in the
// header. The caller only expresses intent — it has no access to the key
// material, so it can't build the value itself.
func (q *Query) WithX5c() *Query {
	q.spec.IncludeX5c = true
	return q
}

// WithJWK asks security to embed the key pair's public key as a JWK in the
// header, for the same reason as WithX5c.
func (q *Query) WithJWK() *Query {
	q.spec.IncludeJWK = true
	return q
}

func (q *Query) Type() query.Type { return QueryType }

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler { return &QueryHandler{service: service} }

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.GenerateJWT)
}
