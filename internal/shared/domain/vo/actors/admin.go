package actors

import (
	vo2 "github.com/isaac-olle/gas/internal/shared/domain/vo"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

type AdminActor struct {
	Id ids.AdminId
}

func NewAdminActorFromString(id string) (*AdminActor, error) {
	adminId, err := ids.NewAdminIdFromString(id)
	if err != nil {
		return nil, err
	}
	return &AdminActor{Id: adminId}, nil
}

func (a *AdminActor) EqualsRole(other vo2.Role) bool {
	return other.Equals(vo2.AdminRole)
}

func (a *AdminActor) ToPrimitives() *vo2.ActorPrimitives {
	return &vo2.ActorPrimitives{
		Id:   a.Id.String(),
		Role: vo2.AdminRole.String(),
	}
}
