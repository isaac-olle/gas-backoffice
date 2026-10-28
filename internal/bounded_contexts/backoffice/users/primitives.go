package users

import "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"

type UserPrimitives struct {
	Id              string                      `json:"id"`
	Name            string                      `json:"name"`
	Email           string                      `json:"email"`
	SystemRole      string                      `json:"system_role"`
	IsSystem        bool                        `json:"is_system"`
	Version         int                         `json:"version"`
	AuthorId        string                      `json:"author_id"`
	AuthorRole      string                      `json:"author_role"`
	SubscriptionIds []string                    `json:"subscription_ids"`
	Memberships     []vo.SubscriptionMembership `json:"memberships"`
}

func (u *User) ToPrimitives() UserPrimitives {
	ids := make([]string, 0, u.subscriptionIds.Size())
	for _, id := range u.subscriptionIds.Values() {
		ids = append(ids, id.String())
	}
	return UserPrimitives{Id: u.id.String(), Name: u.name.String(), Email: u.email.String(), SystemRole: u.systemRole.String(), IsSystem: u.isSystem, Version: u.version, AuthorId: u.authorId.String(), AuthorRole: u.authorRole.String(), SubscriptionIds: ids, Memberships: u.memberships}
}
