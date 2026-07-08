package actors

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type UserActor struct {
	Id      ids2.UserId
	AdminId ids2.AdminId
}

func NewUserActorFromString(id string) (*UserActor, error) {
	userId, err := ids2.NewUserIdFromString(id)
	if err != nil {
		return nil, err
	}
	return &UserActor{Id: userId}, nil
}

func (a *UserActor) EqualsRole(other vo2.Role) bool {
	return other.Equals(vo2.UserRole)
}

func (a *UserActor) ToPrimitives() *vo2.ActorPrimitives {
	return &vo2.ActorPrimitives{
		Id:   a.Id.String(),
		Role: vo2.UserRole.String(),
	}
}

func (a *UserActor) Ownership() users.Ownership {
	return users.NewOwnership(a.AdminId)
}
