package protected_resources

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

var protectedResourceAggregateName aggregate2.Name = aggregate2.ProtectedResource

type ProtectedResource struct {
	*aggregate2.Root
	id        ids2.ProtectedResourceId
	name      vos.ProtectedResourceName
	ownership Ownership
	isSystem  bool
	enabled   bool
	scopeIds  *vo.SetVo[ids2.ScopeId]
	version   int
}

func NewProtectedResource(prId ids2.ProtectedResourceId, name string, ownership Ownership, isSystem bool) (*ProtectedResource, error) {
	errArr := app_errors.NewEmptyErrorArray()

	prName, err := vos.NewProtectedResourceName(name)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	pr := &ProtectedResource{
		Root:      aggregate2.NewRoot(protectedResourceAggregateName),
		id:        prId,
		name:      prName,
		ownership: ownership,
		isSystem:  isSystem,
		enabled:   true,
		scopeIds:  vo.NewEmptySetVo[ids2.ScopeId](),
		version:   1,
	}
	pr.Record(backoffice.NewProtectedResourceCreated(pr.id.String(), pr.name.String(), pr.ownership.SubscriptionId.String()))
	return pr, nil
}

func (pr *ProtectedResource) UpdateName(newName vos.ProtectedResourceName) {
	pr.name = newName
	pr.Record(backoffice.NewProtectedResourceNameUpdated(pr.id.String(), pr.name.String()))
}

func FromPrimitives(p ProtectedResourcePrimitives) (*ProtectedResource, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewProtectedResourceIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name := vos.ProtectedResourceNameFromPrimitives(p.Name)

	ownership, err := NewOwnershipFromStrings(p.SubscriptionId)
	if err != nil {
		errArr.AddError(err)
	}

	scopeIds := vo.NewEmptySetVo[ids2.ScopeId]()
	for _, raw := range p.ScopeIds {
		scopeId, scErr := ids2.NewScopeIdFromString(raw)
		if scErr != nil {
			errArr.AddError(scErr)
		} else {
			scopeIds.Add(scopeId)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	// A persisted protected resource always has version >= 1 (the table
	// defaults new rows to 1, and Save always writes a positive version). A
	// zero only ever comes from a hand-built ProtectedResourcePrimitives{} in
	// a test that didn't set it.
	version := p.Version
	if version == 0 {
		version = 1
	}

	pr := &ProtectedResource{
		Root:      aggregate2.NewRoot(protectedResourceAggregateName),
		id:        id,
		name:      name,
		ownership: ownership,
		isSystem:  p.IsSystem,
		enabled:   p.Enabled,
		scopeIds:  scopeIds,
		version:   version,
	}
	return pr, nil
}

func (pr *ProtectedResource) IsEnabled() bool { return pr.enabled }

func (pr *ProtectedResource) Enable() {
	pr.enabled = true
	pr.Record(backoffice.NewProtectedResourceEnabled(pr.id.String()))
}

func (pr *ProtectedResource) Disable() {
	pr.enabled = false
	pr.Record(backoffice.NewProtectedResourceDisabled(pr.id.String()))
}

func (pr *ProtectedResource) IsSystemResource() bool { return pr.isSystem }

// AddScopeId and RemoveScopeId keep the denormalized, event-driven read cache
// of this protected resource's scopes in sync. They record no event: the
// primary domain action is the Scope's own creation/deletion.
func (pr *ProtectedResource) AddScopeId(scopeId ids2.ScopeId) {
	pr.scopeIds.Add(scopeId)
}

func (pr *ProtectedResource) RemoveScopeId(scopeId ids2.ScopeId) {
	pr.scopeIds.Remove(scopeId)
}

func (pr *ProtectedResource) Delete() {
	pr.Record(backoffice.NewProtectedResourceDeleted(pr.id.String(), pr.ownership.SubscriptionId.String()))
}

// EnsureOwnership allows the caller when the protected resource belongs to
// the given subscription. Membership of the caller in that subscription is
// validated earlier, at the subscription path param middleware.
func (pr *ProtectedResource) EnsureOwnership(subscriptionIds []ids2.SubscriptionId) error {
	for _, id := range subscriptionIds {
		if id == pr.ownership.SubscriptionId {
			return nil
		}
	}
	return app_errors.ForbiddenError(ErrInvalidOwnership).
		WithDescription("protected resource does not belong to the provided subscription")
}
