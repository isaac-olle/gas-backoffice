package subscription

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// Resolver resolves the subscription relevant to the current request. It
// does not authorize the caller; SubscriptionMiddleware does that.
type Resolver interface {
	Resolve(context.Context, *gin.Context) (ids.SubscriptionId, error)
}

type PathResolver struct{ param string }

func NewPathResolver(param string) *PathResolver { return &PathResolver{param: param} }

func (r *PathResolver) Resolve(_ context.Context, c *gin.Context) (ids.SubscriptionId, error) {
	return ids.NewSubscriptionIdFromString(c.Param(r.param))
}
