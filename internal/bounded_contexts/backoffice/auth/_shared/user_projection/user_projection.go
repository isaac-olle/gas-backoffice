package user_projection

import "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"

type UserProjection struct {
	UserId          string
	Email           string
	SystemRole      string
	AdminId         string
	SubscriptionIds []string
	Memberships     []vo.SubscriptionMembership
}
