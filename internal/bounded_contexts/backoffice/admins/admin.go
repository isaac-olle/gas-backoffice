package admins

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

var adminAggregateName = aggregate2.Admin

type Admin struct {
	*aggregate2.Root
	id        ids2.AdminId
	name      vos.AdminName
	ownership Ownership
	userIds   *vo.SetVo[ids2.UserId]
	isSystem  bool
}

func FromPrimitives(p AdminPrimitives) (*Admin, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewAdminIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name := vos.AdminNameFromPrimitives(p.Name)

	userIds := vo.NewEmptySetVo[ids2.UserId]()
	for _, rawId := range p.UserIds {
		userId, tErr := ids2.NewUserIdFromString(rawId)
		if tErr != nil {
			errArr.AddError(tErr)
		} else {
			userIds.Add(userId)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	a := &Admin{
		Root:      aggregate2.NewRoot(adminAggregateName),
		id:        id,
		name:      name,
		ownership: newOwnership(id),
		userIds:   userIds,
		isSystem:  p.IsSystem,
	}
	return a, nil
}

func (a *Admin) IsSystemResource() bool { return a.isSystem }

func (a *Admin) AddUserId(userId ids2.UserId) {
	a.userIds.Add(userId)
}

func (a *Admin) RemoveUserId(userId ids2.UserId) {
	a.userIds.Remove(userId)
}
