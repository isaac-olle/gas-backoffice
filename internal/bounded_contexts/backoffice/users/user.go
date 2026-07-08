package users

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

func NewUser(userId ids2.UserId, name string, ownership Ownership) (*User, error) {
	errArr := app_errors.NewEmptyErrorArray()

	userName, err := vos.NewUserName(name)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	t := &User{
		Root:                 aggregate2.NewRoot(userAggregateName),
		id:                   userId,
		name:                 userName,
		ownership:            ownership,
		clientIds:            vo.NewEmptySetVo[ids2.ClientId](),
		protectedResourceIds: vo.NewEmptySetVo[ids2.ProtectedResourceId](),
	}
	t.Record(backoffice.NewUserCreated(t.id.String(), t.name.String(), t.ownership.AdminId.String()))
	return t, nil
}

var userAggregateName aggregate2.Name = aggregate2.User

type User struct {
	*aggregate2.Root
	id                   ids2.UserId
	name                 vos.UserName
	ownership            Ownership
	clientIds            *vo.SetVo[ids2.ClientId]
	protectedResourceIds *vo.SetVo[ids2.ProtectedResourceId]
	isSystem             bool
}

func FromPrimitives(p UserPrimitives) (*User, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewUserIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name := vos.UserNameFromPrimitives(p.Name)

	adminId, err := ids2.NewAdminIdFromString(p.AdminId)
	if err != nil {
		errArr.AddError(err)
	}

	clientIds := vo.NewEmptySetVo[ids2.ClientId]()
	for _, rawId := range p.ClientIds {
		clientId, cErr := ids2.NewClientIdFromString(rawId)
		if cErr != nil {
			errArr.AddError(cErr)
		} else {
			clientIds.Add(clientId)
		}
	}

	protectedResourceIds := vo.NewEmptySetVo[ids2.ProtectedResourceId]()
	for _, rawId := range p.ProtectedResourceIds {
		prId, prErr := ids2.NewProtectedResourceIdFromString(rawId)
		if prErr != nil {
			errArr.AddError(prErr)
		} else {
			protectedResourceIds.Add(prId)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	t := &User{
		Root:                 aggregate2.NewRoot(userAggregateName),
		id:                   id,
		name:                 name,
		ownership:            NewOwnership(adminId),
		clientIds:            clientIds,
		protectedResourceIds: protectedResourceIds,
		isSystem:             p.IsSystem,
	}
	return t, nil
}

func (t *User) IsSystemResource() bool { return t.isSystem }

func (t *User) CheckOwnership(o Ownership) bool {
	return t.ownership.Equals(o)
}

func (t *User) ChangeName(newName vos.UserName) error {
	t.name = newName
	t.Record(backoffice.NewUserNameChanged(t.id.String(), t.name.String()))
	return nil
}

func (t *User) AddClientId(clientId ids2.ClientId)    { t.clientIds.Add(clientId) }
func (t *User) RemoveClientId(clientId ids2.ClientId) { t.clientIds.Remove(clientId) }
func (t *User) AddProtectedResourceId(prId ids2.ProtectedResourceId) {
	t.protectedResourceIds.Add(prId)
}
func (t *User) RemoveProtectedResourceId(prId ids2.ProtectedResourceId) {
	t.protectedResourceIds.Remove(prId)
}

func (t *User) Delete() {
	t.Record(
		backoffice.NewUserDeleted(t.id.String(), t.ownership.AdminId.String()),
	)
}
