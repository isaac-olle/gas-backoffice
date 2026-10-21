package tenants

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

var tenantAggregateName aggregate.Name = aggregate.User

type Tenant struct {
	*aggregate.Root
	id              ids.UserId
	adminId         ids.UserId
	name            vos.UserName
	email           vos.Email
	ownership       Ownership
	isSystem        bool
	version         int
	subscriptionIds *vo.SetVo[ids.SubscriptionId]
}

func NewTenant(userId ids.UserId, name string, adminId ids.UserId, email string, isSystem bool) (*Tenant, error) {
	errArr := app_errors.NewEmptyErrorArray()

	tenantName, err := vos.NewUserName(name)
	if err != nil {
		errArr.AddError(err)
	}

	tenantEmail, err := vos.NewEmail(email)
	if err != nil {
		errArr.AddError(app_errors.BadRequestError(ErrInvalidEmail).WithDescription("invalid email").WithCause(err, false))
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	t := &Tenant{
		Root:            aggregate.NewRoot(tenantAggregateName),
		id:              userId,
		adminId:         adminId,
		name:            tenantName,
		email:           tenantEmail,
		ownership:       NewOwnership(userId),
		isSystem:        isSystem,
		version:         1,
		subscriptionIds: vo.NewEmptySetVo[ids.SubscriptionId](),
	}
	t.Record(backoffice.NewUserCreated(t.id.String(), t.name.String()))
	return t, nil
}

func FromPrimitives(p TenantPrimitives) (*Tenant, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids.NewUserIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name := vos.UserNameFromPrimitives(p.Name)
	email := vos.EmailFromPrimitives(p.Email)

	adminId, err := ids.NewUserIdFromString(p.AdminId)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
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

	// A persisted tenant always has version >= 1 (the table defaults new rows
	// to 1, and Save always writes a positive version). A zero only ever comes
	// from a hand-built TenantPrimitives{} in a test that didn't set it.
	version := p.Version
	if version == 0 {
		version = 1
	}

	t := &Tenant{
		Root:            aggregate.NewRoot(tenantAggregateName),
		id:              id,
		adminId:         adminId,
		name:            name,
		email:           email,
		ownership:       NewOwnership(id),
		isSystem:        p.IsSystem,
		version:         version,
		subscriptionIds: subscriptionIds,
	}
	return t, nil
}

func (t *Tenant) IsSystemResource() bool { return t.isSystem }

func (t *Tenant) AddSubscriptionId(id ids.SubscriptionId)    { t.subscriptionIds.Add(id) }
func (t *Tenant) RemoveSubscriptionId(id ids.SubscriptionId) { t.subscriptionIds.Remove(id) }

func (t *Tenant) IsSelf(author *authors.UserAuthor) bool {
	return t.id == author.Id
}

func (t *Tenant) IsOwnedByAdmin(author *authors.UserAuthor) bool {
	return t.adminId == author.Id
}

func (t *Tenant) EnsureTenantAccess(author *authors.UserAuthor) error {
	var allowed bool
	switch {
	case author.Role.IsAdmin():
		allowed = t.IsOwnedByAdmin(author)
	case author.Role.IsTenant():
		allowed = t.IsSelf(author)
	}

	if !allowed {
		return app_errors.ForbiddenError(ErrInvalidOwnership).
			WithDescription("caller does not have access to this tenant")
	}
	return nil
}

func (t *Tenant) ChangeName(newName vos.UserName) error {
	t.name = newName
	t.Record(backoffice.NewUserNameChanged(t.id.String(), t.name.String()))
	return nil
}

func (t *Tenant) Delete() {
	t.Record(backoffice.NewUserDeleted(t.id.String()))
}
