package authors

import (
	"context"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// UserAuthor is the authenticated caller. AdminId is the id of the admin that
// owns the caller (for an admin it is the admin's own id); it is an ids.UserId
// because an admin is just a User.
type UserAuthor struct {
	Id              ids2.UserId
	AdminId         ids2.UserId
	Role            vo2.Role
	SubscriptionIds []ids2.SubscriptionId
	Memberships     []vo2.SubscriptionMembership
}

func (a *UserAuthor) MembershipFor(subscriptionId ids2.SubscriptionId) (vo2.SubscriptionMembership, bool) {
	if a == nil {
		return vo2.SubscriptionMembership{}, false
	}
	for _, membership := range a.Memberships {
		if membership.SubscriptionId == subscriptionId {
			return membership, true
		}
	}
	return vo2.SubscriptionMembership{}, false
}

func SubscriptionIdsFromContext(ctx context.Context) []ids2.SubscriptionId {
	author, ok := vo2.AuthorFromContext[*UserAuthor](ctx)
	if !ok {
		return nil
	}
	return author.SubscriptionIds
}

func (a *UserAuthor) SoleSubscriptionId() (ids2.SubscriptionId, bool) {
	if a == nil || len(a.SubscriptionIds) != 1 {
		return ids2.SubscriptionId{}, false
	}
	return a.SubscriptionIds[0], true
}

func NewUserAuthorFromString(id, adminId string, role vo2.Role, rawSubscriptionIds ...[]string) (*UserAuthor, error) {
	userId, err := ids2.NewUserIdFromString(id)
	if err != nil {
		return nil, err
	}
	adminUserId, err := ids2.NewUserIdFromString(adminId)
	if err != nil {
		return nil, err
	}
	var raws []string
	if len(rawSubscriptionIds) > 0 {
		raws = rawSubscriptionIds[0]
	}
	subscriptions := make([]ids2.SubscriptionId, 0, len(raws))
	for _, raw := range raws {
		sid, e := ids2.NewSubscriptionIdFromString(raw)
		if e != nil {
			return nil, e
		}
		subscriptions = append(subscriptions, sid)
	}
	return &UserAuthor{Id: userId, AdminId: adminUserId, Role: role, SubscriptionIds: subscriptions}, nil
}

func NewUserAuthorWithMemberships(id, adminId string, role vo2.Role, memberships []vo2.SubscriptionMembership) (*UserAuthor, error) {
	author, err := NewUserAuthorFromString(id, adminId, role)
	if err != nil {
		return nil, err
	}
	author.Memberships = memberships
	author.SubscriptionIds = make([]ids2.SubscriptionId, 0, len(memberships))
	for _, membership := range memberships {
		author.SubscriptionIds = append(author.SubscriptionIds, membership.SubscriptionId)
	}
	return author, nil
}

func (a *UserAuthor) EqualsRole(other vo2.Role) bool {
	return a.Role.Equals(other)
}

func (a *UserAuthor) ToPrimitives() *vo2.AuthorPrimitives {
	return &vo2.AuthorPrimitives{
		Id:   a.Id.String(),
		Role: a.Role.String(),
	}
}
