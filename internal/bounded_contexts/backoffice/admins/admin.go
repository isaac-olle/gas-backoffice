package admins

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

var adminAggregateName aggregate.Name = aggregate.Admin

// Admin is a manager: it owns tenants and controls their lifecycle. It has no
// role field and no adminId field (unlike tenants.Tenant) because there is no
// ownership hierarchy above an Admin — see EnsureSelf.
type Admin struct {
	*aggregate.Root
	id              ids.UserId
	name            vos.UserName
	email           vos.Email
	isSystem        bool
	version         int
	subscriptionIds *vo.SetVo[ids.SubscriptionId]
}

func NewAdmin(id ids.UserId, name string, email string, isSystem bool) (*Admin, error) {
	errArr := app_errors.NewEmptyErrorArray()

	adminName, err := vos.NewUserName(name)
	if err != nil {
		errArr.AddError(err)
	}

	adminEmail, err := vos.NewEmail(email)
	if err != nil {
		errArr.AddError(app_errors.BadRequestError(ErrInvalidEmail).WithDescription("invalid email").WithCause(err, false))
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	a := &Admin{
		Root:            aggregate.NewRoot(adminAggregateName),
		id:              id,
		name:            adminName,
		email:           adminEmail,
		isSystem:        isSystem,
		version:         1,
		subscriptionIds: vo.NewEmptySetVo[ids.SubscriptionId](),
	}
	a.Record(backoffice.NewAdminCreated(a.id.String(), a.name.String()))
	return a, nil
}

func FromPrimitives(p AdminPrimitives) (*Admin, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids.NewUserIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name := vos.UserNameFromPrimitives(p.Name)
	email := vos.EmailFromPrimitives(p.Email)

	if errArr.HasErrors() {
		return nil, errArr
	}

	// A persisted admin always has version >= 1 (the table defaults new rows
	// to 1, and Save always writes a positive version). A zero only ever comes
	// from a hand-built AdminPrimitives{} in a test that didn't set it.
	version := p.Version
	if version == 0 {
		version = 1
	}
	subscriptionIds := vo.NewEmptySetVo[ids.SubscriptionId]()
	for _, raw := range p.SubscriptionIds {
		sid, parseErr := ids.NewSubscriptionIdFromString(raw)
		if parseErr != nil {
			errArr.AddError(parseErr)
		} else {
			subscriptionIds.Add(sid)
		}
	}
	if errArr.HasErrors() {
		return nil, errArr
	}

	a := &Admin{
		Root:            aggregate.NewRoot(adminAggregateName),
		id:              id,
		name:            name,
		email:           email,
		isSystem:        p.IsSystem,
		version:         version,
		subscriptionIds: subscriptionIds,
	}
	return a, nil
}

func (a *Admin) IsSystemResource() bool { return a.isSystem }

func (a *Admin) AddSubscriptionId(id ids.SubscriptionId)    { a.subscriptionIds.Add(id) }
func (a *Admin) RemoveSubscriptionId(id ids.SubscriptionId) { a.subscriptionIds.Remove(id) }

// EnsureSelf is the only access check an Admin needs: there is no hierarchy
// above an Admin, so a caller may only act on their own admin record.
func (a *Admin) EnsureSelf(callerId ids.UserId) error {
	if a.id != callerId {
		return app_errors.ForbiddenError(ErrInvalidOwnership).
			WithDescription("caller does not have access to this admin")
	}
	return nil
}

func (a *Admin) ChangeName(newName vos.UserName) error {
	a.name = newName
	a.Record(backoffice.NewAdminNameChanged(a.id.String(), a.name.String()))
	return nil
}
