package subscriptions

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Repository interface {
	Save(context.Context, *Subscription) error
	SearchById(context.Context, ids.SubscriptionId) (*Subscription, error)
	List(context.Context, []ids.SubscriptionId, bool) ([]Summary, error)
}

type Summary struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
