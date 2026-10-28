package users

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

var userAggregateName aggregate.Name = aggregate.User

type User struct {
	*aggregate.Root
	id              ids.UserId
	name            vos.UserName
	email           vos.Email
	systemRole      vo.Role
	isSystem        bool
	version         int
	authorId        ids.UserId
	authorRole      vo.Role
	subscriptionIds *vo.SetVo[ids.SubscriptionId]
	memberships     []vo.SubscriptionMembership
}

func NewUser(id ids.UserId, name, email string, role vo.Role, isSystem bool) (*User, error) {
	es := app_errors.NewEmptyErrorArray()
	n, e := vos.NewUserName(name)
	if e != nil {
		es.AddError(e)
	}
	em, e := vos.NewEmail(email)
	if e != nil {
		es.AddError(app_errors.BadRequestError(ErrInvalidEmail).WithDescription("invalid email").WithCause(e, false))
	}
	if !role.IsSystemAdministrator() && !role.IsSubscriptionUser() {
		es.AddError(app_errors.BadRequestError(ErrInvalidEmail).WithDescription("invalid system role"))
	}
	if es.HasErrors() {
		return nil, es
	}
	u := &User{Root: aggregate.NewRoot(userAggregateName), id: id, name: n, email: em, systemRole: role, isSystem: isSystem, version: 1, subscriptionIds: vo.NewEmptySetVo[ids.SubscriptionId](), memberships: []vo.SubscriptionMembership{}}
	u.Record(backoffice.NewUserCreated(id.String(), n.String()))
	return u, nil
}
func NewSystemAdministrator(id ids.UserId, name, email string) (*User, error) {
	return NewUser(id, name, email, vo.SystemAdministratorRole, true)
}
func NewSubscriptionUser(id ids.UserId, name, email string) (*User, error) {
	return NewUser(id, name, email, vo.SubscriptionUserRole, false)
}
func FromPrimitives(p UserPrimitives) (*User, error) {
	es := app_errors.NewEmptyErrorArray()
	id, e := ids.NewUserIdFromString(p.Id)
	if e != nil {
		es.AddError(e)
	}
	role, e := vo.NewRole(p.SystemRole)
	if e != nil {
		es.AddError(e)
	}
	var aid ids.UserId
	if p.AuthorId != "" {
		aid, e = ids.NewUserIdFromString(p.AuthorId)
		if e != nil {
			es.AddError(e)
		}
	}
	ar := vo.Role("")
	if p.AuthorRole != "" {
		ar, e = vo.NewRole(p.AuthorRole)
		if e != nil {
			es.AddError(e)
		}
	}
	if es.HasErrors() {
		return nil, es
	}
	set := vo.NewEmptySetVo[ids.SubscriptionId]()
	ms := make([]vo.SubscriptionMembership, 0, len(p.Memberships))
	for _, m := range p.Memberships {
		set.Add(m.SubscriptionId)
		ms = append(ms, m)
	}
	for _, raw := range p.SubscriptionIds {
		sid, x := ids.NewSubscriptionIdFromString(raw)
		if x != nil {
			es.AddError(x)
		} else {
			set.Add(sid)
		}
	}
	if es.HasErrors() {
		return nil, es
	}
	v := p.Version
	if v == 0 {
		v = 1
	}
	return &User{Root: aggregate.NewRoot(userAggregateName), id: id, name: vos.UserNameFromPrimitives(p.Name), email: vos.EmailFromPrimitives(p.Email), systemRole: role, isSystem: p.IsSystem, version: v, authorId: aid, authorRole: ar, subscriptionIds: set, memberships: ms}, nil
}
func (u *User) IsSystemResource() bool                     { return u.isSystem }
func (u *User) IsSystemAdministrator() bool                { return u.systemRole.IsSystemAdministrator() }
func (u *User) SystemRole() vo.Role                        { return u.systemRole }
func (u *User) AddSubscriptionId(id ids.SubscriptionId)    { u.subscriptionIds.Add(id) }
func (u *User) RemoveSubscriptionId(id ids.SubscriptionId) { u.subscriptionIds.Remove(id) }
func (u *User) AddMembership(m vo.SubscriptionMembership) {
	for i, x := range u.memberships {
		if x.SubscriptionId == m.SubscriptionId {
			u.memberships[i] = m
			return
		}
	}
	u.memberships = append(u.memberships, m)
	u.subscriptionIds.Add(m.SubscriptionId)
}
func (u *User) RemoveMembership(id ids.SubscriptionId) {
	for i, m := range u.memberships {
		if m.SubscriptionId == id {
			u.memberships = append(u.memberships[:i], u.memberships[i+1:]...)
			break
		}
	}
	u.subscriptionIds.Remove(id)
}
func (u *User) IsSelf(a *authors.UserAuthor) bool { return a != nil && u.id == a.Id }
func (u *User) EnsureSelf(id ids.UserId) error {
	if u.id != id {
		return app_errors.ForbiddenError(ErrInvalidOwnership).WithDescription("caller does not have access to this user")
	}
	return nil
}
func (u *User) EnsureUserAccess(a *authors.UserAuthor) error {
	if a == nil || (!a.Role.IsSystemAdministrator() && !u.IsSelf(a)) {
		return app_errors.ForbiddenError(ErrInvalidOwnership).WithDescription("caller does not have access to this user")
	}
	return nil
}
func (u *User) ChangeName(n vos.UserName) error {
	u.name = n
	u.Record(backoffice.NewUserNameChanged(u.id.String(), n.String()))
	return nil
}
func (u *User) Delete() { u.Record(backoffice.NewUserDeleted(u.id.String())) }
