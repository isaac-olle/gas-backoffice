package list

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const QueryType query.Type = "query.backoffice.subscription.list"

type Query struct {
	all             bool
	subscriptionIds []ids.SubscriptionId
}

func NewQuery(all bool, subscriptionIds []ids.SubscriptionId) *Query {
	return &Query{all: all, subscriptionIds: subscriptionIds}
}

func (q *Query) Type() query.Type { return QueryType }

type Service struct{ repository subscriptions.Repository }

func NewService(repository subscriptions.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, q *Query) (any, error) {
	result, err := s.repository.List(ctx, q.subscriptionIds, q.all)
	if err != nil {
		return nil, app_errors.InternalError(subscriptions.ErrSubscriptionsNotListed).WithCause(err, false)
	}
	return result, nil
}

type QueryHandler struct{ service *Service }

func NewQueryHandler(service *Service) *QueryHandler { return &QueryHandler{service: service} }

func (h *QueryHandler) Handle(ctx context.Context, qry query.Query) (any, error) {
	return query.Process[*Query, context.Context](ctx, qry, h.service.List)
}
