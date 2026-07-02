package admin_assertion_keys

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

type Ownership struct {
	AdminId ids.AdminId
}

func newOwnership(adminId ids.AdminId) Ownership {
	return Ownership{AdminId: adminId}
}

func NewOwnershipFromVos(adminId ids.AdminId) Ownership {
	return newOwnership(adminId)
}

func NewOwnership(adminId string) (Ownership, error) {
	aId, err := ids.NewAdminIdFromString(adminId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid admin id").
			WithCause(err)
	}
	return newOwnership(aId), nil
}
