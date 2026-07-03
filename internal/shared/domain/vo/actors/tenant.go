package actors

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type TenantActor struct {
	Id      ids2.TenantId
	AdminId ids2.AdminId
}

func NewTenantActorFromString(id string) (*TenantActor, error) {
	tenantId, err := ids2.NewTenantIdFromString(id)
	if err != nil {
		return nil, err
	}
	return &TenantActor{Id: tenantId}, nil
}

func (a *TenantActor) EqualsRole(other vo2.Role) bool {
	return other.Equals(vo2.TenantRole)
}

func (a *TenantActor) ToPrimitives() *vo2.ActorPrimitives {
	return &vo2.ActorPrimitives{
		Id:   a.Id.String(),
		Role: vo2.TenantRole.String(),
	}
}

func (a *TenantActor) Ownership() tenants.Ownership {
	return tenants.NewOwnership(a.AdminId)
}
