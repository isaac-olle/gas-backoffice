package actors

import (
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type ClientActor struct {
	Id ids.ClientId
}

func NewClientActorFromString(id string) (*ClientActor, error) {
	clientId, err := ids.NewClientIdFromString(id)
	if err != nil {
		return nil, err
	}
	return &ClientActor{Id: clientId}, nil
}

func (a *ClientActor) EqualsRole(other vo2.Role) bool {
	return other.Equals(vo2.ClientRole)
}

func (a *ClientActor) ToPrimitives() *vo2.ActorPrimitives {
	return &vo2.ActorPrimitives{
		Id:   a.Id.String(),
		Role: vo2.ClientRole.String(),
	}
}
