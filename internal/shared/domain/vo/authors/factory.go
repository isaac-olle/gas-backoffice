package authors

import (
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func NewAuthorForRole(role vo.Role, sub, adminId string, rawSubscriptionIds ...[]string) (vo.Author, error) {
	switch {
	case role.IsAdmin(), role.IsTenant():
		userId, err := ids.NewUserIdFromString(sub)
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
		return &UserAuthor{Id: userId, AdminId: adminUserId, Role: role, SubscriptionIds: subscriptions}, nil

	default:
		return nil, fmt.Errorf("unknown role: %s", role)
	}
}
