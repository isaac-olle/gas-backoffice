package authors

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// SubscriptionAccessFromContext returns the subscriptions carried by the
// verified author token and whether the author has global access.
func SubscriptionAccessFromContext(ctx context.Context) ([]ids.SubscriptionId, bool, bool) {
	author, ok := vo.AuthorFromContext[*UserAuthor](ctx)
	if !ok || author == nil {
		return nil, false, false
	}
	return author.SubscriptionIds, author.Role.IsSystemAdministrator(), true
}
