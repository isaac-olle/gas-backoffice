package users

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

var userAggregateName aggregate.Name = aggregate.User

type User struct {
	*aggregate.Root
	id                   ids.UserId
	name                 vos.UserName
	role                 vo.Role
	ownership            Ownership
	clientIds            *vo.SetVo[ids.ClientId]
	protectedResourceIds *vo.SetVo[ids.ProtectedResourceId]
	isSystem             bool
}

func NewUser(userId ids.UserId, name string, ownership Ownership) (*User, error) {
	userName, err := vos.NewUserName(name)
	if err != nil {
		return nil, err
	}

	t := &User{
		Root:                 aggregate.NewRoot(userAggregateName),
		id:                   userId,
		name:                 userName,
		role:                 vo.UserRole,
		ownership:            ownership,
		clientIds:            vo.NewEmptySetVo[ids.ClientId](),
		protectedResourceIds: vo.NewEmptySetVo[ids.ProtectedResourceId](),
	}
	t.Record(backoffice.NewUserCreated(t.id.String(), t.name.String()))
	return t, nil
}

func FromPrimitives(p UserPrimitives) (*User, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids.NewUserIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name := vos.UserNameFromPrimitives(p.Name)

	role, err := vo.NewRole(p.Role)
	if err != nil {
		errArr.AddError(err)
	}

	clientIds := vo.NewEmptySetVo[ids.ClientId]()
	for _, rawId := range p.ClientIds {
		clientId, cErr := ids.NewClientIdFromString(rawId)
		if cErr != nil {
			errArr.AddError(cErr)
		} else {
			clientIds.Add(clientId)
		}
	}

	protectedResourceIds := vo.NewEmptySetVo[ids.ProtectedResourceId]()
	for _, rawId := range p.ProtectedResourceIds {
		prId, prErr := ids.NewProtectedResourceIdFromString(rawId)
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
		Root:                 aggregate.NewRoot(userAggregateName),
		id:                   id,
		name:                 name,
		role:                 role,
		ownership:            NewOwnership(id),
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

func (t *User) AddClientId(clientId ids.ClientId)    { t.clientIds.Add(clientId) }
func (t *User) RemoveClientId(clientId ids.ClientId) { t.clientIds.Remove(clientId) }
func (t *User) AddProtectedResourceId(prId ids.ProtectedResourceId) {
	t.protectedResourceIds.Add(prId)
}
func (t *User) RemoveProtectedResourceId(prId ids.ProtectedResourceId) {
	t.protectedResourceIds.Remove(prId)
}

func (t *User) Delete() {
	t.Record(backoffice.NewUserDeleted(t.id.String()))
}
