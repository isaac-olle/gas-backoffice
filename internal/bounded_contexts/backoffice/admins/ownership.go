package admins

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	AdminId ids.AdminId
}

func newOwnership(adminId ids.AdminId) Ownership {
	return Ownership{AdminId: adminId}
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}

func NewOwnership(adminId string) (Ownership, error) {
	id, err := ids.NewAdminIdFromString(adminId)
	if err != nil {
		return Ownership{}, app_errors.AddTrace(err)
	}
	return newOwnership(id), nil
}
