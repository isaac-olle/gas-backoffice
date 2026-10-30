package vo

import (
	"context"
	"sort"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type authorContextKey struct{}

func WithAuthor(ctx context.Context, author *UserAuthor) context.Context {
	return context.WithValue(ctx, authorContextKey{}, author)
}

func AuthorFromContext(ctx context.Context) (*UserAuthor, bool) {
	author, ok := ctx.Value(authorContextKey{}).(*UserAuthor)
	return author, ok
}

// UserAuthor is the authenticated caller. AdminId is the id of the admin that
// owns the caller (for an admin it is the admin's own id).
type UserAuthor struct {
	Id          ids.UserId
	AdminId     ids.UserId
	Role        Role
	Memberships map[ids.SubscriptionId]SubscriptionRole
}

type AuthorPrimitives struct {
	Id   string
	Role string
}

func (a *UserAuthor) MembershipFor(subscriptionId ids.SubscriptionId) (SubscriptionRole, bool) {
	if a == nil {
		return "", false
	}
	role, ok := a.Memberships[subscriptionId]
	return role, ok
}

// SubscriptionIDs returns the subscriptions available to the author in a
// stable order. Memberships is the source of truth; the slice is only an
// adapter for consumers that need to pass a list around.
func (a *UserAuthor) SubscriptionIDs() []ids.SubscriptionId {
	if a == nil {
		return nil
	}
	result := make([]ids.SubscriptionId, 0, len(a.Memberships))
	for subscriptionID := range a.Memberships {
		result = append(result, subscriptionID)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].String() < result[j].String() })
	return result
}

func NewUserAuthorFromString(id, adminId string, role Role, rawSubscriptionIds ...[]string) (*UserAuthor, error) {
	userId, err := ids.NewUserIdFromString(id)
	if err != nil {
		return nil, err
	}
	adminUserId, err := ids.NewUserIdFromString(adminId)
	if err != nil {
		return nil, err
	}
	var raws []string
	if len(rawSubscriptionIds) > 0 {
		raws = rawSubscriptionIds[0]
	}
	subscriptions := make([]ids.SubscriptionId, 0, len(raws))
	for _, raw := range raws {
		sid, e := ids.NewSubscriptionIdFromString(raw)
		if e != nil {
			return nil, e
		}
		subscriptions = append(subscriptions, sid)
	}
	memberships := make(map[ids.SubscriptionId]SubscriptionRole, len(subscriptions))
	for _, subscriptionID := range subscriptions {
		memberships[subscriptionID] = UserSubscriptionRole
	}
	return &UserAuthor{Id: userId, AdminId: adminUserId, Role: role, Memberships: memberships}, nil
}

func NewUserAuthorWithMemberships(id, adminId string, role Role, memberships []SubscriptionMembership) (*UserAuthor, error) {
	author, err := NewUserAuthorFromString(id, adminId, role)
	if err != nil {
		return nil, err
	}
	author.Memberships = make(map[ids.SubscriptionId]SubscriptionRole, len(memberships))
	for _, membership := range memberships {
		author.Memberships[membership.SubscriptionId] = membership.SubscriptionRole
	}
	return author, nil
}

func (a *UserAuthor) EqualsRole(other Role) bool {
	return a != nil && a.Role.Equals(other)
}

func (a *UserAuthor) ToPrimitives() *AuthorPrimitives {
	return &AuthorPrimitives{Id: a.Id.String(), Role: a.Role.String()}
}
