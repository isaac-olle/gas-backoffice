package protected_resources

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

type Ownership struct {
	TenantId ids2.TenantId
	AdminId  ids2.AdminId
}

func NewOwnership(tenantId ids2.TenantId, adminId ids2.AdminId) Ownership {
	return Ownership{TenantId: tenantId, AdminId: adminId}
}

func NewOwnershipFromStrings(tenantId string, adminId string) (Ownership, error) {
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
	return NewOwnership(tId, aId), nil
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}
