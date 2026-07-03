package scopes

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	ProtectedResourceId ids2.ProtectedResourceId
	TenantId            ids2.TenantId
	AdminId             ids2.AdminId
}

func newOwnership(protectedResourceId ids2.ProtectedResourceId, tenantId ids2.TenantId, adminId ids2.AdminId) Ownership {
	return Ownership{ProtectedResourceId: protectedResourceId, TenantId: tenantId, AdminId: adminId}
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}

func NewOwnershipFromVos(protectedResourceId ids2.ProtectedResourceId, tenantId ids2.TenantId, adminId ids2.AdminId) Ownership {
	return newOwnership(protectedResourceId, tenantId, adminId)
}

func NewOwnership(protectedResourceId, tenantId, adminId string) (Ownership, error) {
	prId, err := ids2.NewProtectedResourceIdFromString(protectedResourceId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid protected resource id").
			WithCause(err)
	}
	tId, err := ids2.NewTenantIdFromString(tenantId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid tenant id").
			WithCause(err)
	}
	aId, err := ids2.NewAdminIdFromString(adminId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid admin id").
			WithCause(err)
	}
	return newOwnership(prId, tId, aId), nil
}
