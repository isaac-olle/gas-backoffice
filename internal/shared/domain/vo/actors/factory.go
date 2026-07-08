package actors

import (
	"fmt"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func NewActorForRole(role vo2.Role, sub string) (vo2.Actor, error) {
	switch {
	case role.IsAdmin():
		adminId, err := ids2.NewAdminIdFromString(sub)
		if err != nil {
			return nil, err
		}
		return &AdminActor{Id: adminId}, nil

	case role.IsUser():
		userId, err := ids2.NewUserIdFromString(sub)
		if err != nil {
			return nil, err
		}
		return &UserActor{Id: userId}, nil

	case role.IsClient():
		clientId, err := ids2.NewClientIdFromString(sub)
		if err != nil {
			return nil, err
		}
		return &ClientActor{Id: clientId}, nil

	default:
		return nil, fmt.Errorf("unknown role: %s", role)
	}
}
