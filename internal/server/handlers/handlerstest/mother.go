package handlerstest

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const (
	AdminID  = "b1b2c3d4-e5f6-7890-abcd-ef1234567890"
	TenantID = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	ClientID = "c1b2c3d4-e5f6-7890-abcd-ef1234567890"
	PrID     = "d1b2c3d4-e5f6-7890-abcd-ef1234567890"
	ScopeID  = "e1b2c3d4-e5f6-7890-abcd-ef1234567890"
)

func ValidAdminIdExtractor() func(context.Context) (ids.AdminId, bool) {
	return func(_ context.Context) (ids.AdminId, bool) {
		id, _ := ids.NewAdminIdFromString(AdminID)
		return id, true
	}
}

func MissingAdminIdExtractor() func(context.Context) (ids.AdminId, bool) {
	return func(_ context.Context) (ids.AdminId, bool) {
		return ids.AdminId{}, false
	}
}

func ValidTenantIdExtractor() func(context.Context) (ids.TenantId, bool) {
	return func(_ context.Context) (ids.TenantId, bool) {
		id, _ := ids.NewTenantIdFromString(TenantID)
		return id, true
	}
}

func MissingTenantIdExtractor() func(context.Context) (ids.TenantId, bool) {
	return func(_ context.Context) (ids.TenantId, bool) {
		return ids.TenantId{}, false
	}
}

func ValidClientIdExtractor() func(context.Context) (ids.ClientId, bool) {
	return func(_ context.Context) (ids.ClientId, bool) {
		id, _ := ids.NewClientIdFromString(ClientID)
		return id, true
	}
}

func MissingClientIdExtractor() func(context.Context) (ids.ClientId, bool) {
	return func(_ context.Context) (ids.ClientId, bool) {
		return ids.ClientId{}, false
	}
}

func ValidPrIdExtractor() func(context.Context) (ids.ProtectedResourceId, bool) {
	return func(_ context.Context) (ids.ProtectedResourceId, bool) {
		id, _ := ids.NewProtectedResourceIdFromString(PrID)
		return id, true
	}
}

func MissingPrIdExtractor() func(context.Context) (ids.ProtectedResourceId, bool) {
	return func(_ context.Context) (ids.ProtectedResourceId, bool) {
		return ids.ProtectedResourceId{}, false
	}
}

func ValidScopeIdExtractor() func(context.Context) (ids.ScopeId, bool) {
	return func(_ context.Context) (ids.ScopeId, bool) {
		id, _ := ids.NewScopeIdFromString(ScopeID)
		return id, true
	}
}

func MissingScopeIdExtractor() func(context.Context) (ids.ScopeId, bool) {
	return func(_ context.Context) (ids.ScopeId, bool) {
		return ids.ScopeId{}, false
	}
}

func NewScopeIdExtractor() func(context.Context) (ids.ScopeId, bool) {
	return func(_ context.Context) (ids.ScopeId, bool) {
		return ids.NewScopeId(), true
	}
}

func ValidTenantOwnershipExtractor() func(context.Context) (tenants.Ownership, bool) {
	return func(_ context.Context) (tenants.Ownership, bool) {
		adminId, _ := ids.NewAdminIdFromString(AdminID)
		return tenants.NewOwnership(adminId), true
	}
}

func MissingTenantOwnershipExtractor() func(context.Context) (tenants.Ownership, bool) {
	return func(_ context.Context) (tenants.Ownership, bool) {
		return tenants.Ownership{}, false
	}
}

func ValidPrOwnershipExtractor() func(context.Context) (protected_resources.Ownership, bool) {
	return func(_ context.Context) (protected_resources.Ownership, bool) {
		adminId, _ := ids.NewAdminIdFromString(AdminID)
		tenantId, _ := ids.NewTenantIdFromString(TenantID)
		return protected_resources.NewOwnership(tenantId, adminId), true
	}
}

func MissingPrOwnershipExtractor() func(context.Context) (protected_resources.Ownership, bool) {
	return func(_ context.Context) (protected_resources.Ownership, bool) {
		return protected_resources.Ownership{}, false
	}
}
