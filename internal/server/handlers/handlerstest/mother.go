package handlerstest

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const (
	AdminID        = "b1b2c3d4-e5f6-7890-abcd-ef1234567890"
	UserId         = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	ClientID       = "c1b2c3d4-e5f6-7890-abcd-ef1234567890"
	PrID           = "d1b2c3d4-e5f6-7890-abcd-ef1234567890"
	ScopeID        = "e1b2c3d4-e5f6-7890-abcd-ef1234567890"
	SessionID      = "f1b2c3d4-e5f6-7890-abcd-ef1234567890"
	SubscriptionID = "01b2c3d4-e5f6-7890-abcd-ef1234567890"
)

func ValidSessionIdExtractor() func(context.Context) (ids.SessionId, bool) {
	return func(_ context.Context) (ids.SessionId, bool) {
		id, _ := ids.NewSessionIdFromString(SessionID)
		return id, true
	}
}

func MissingSessionIdExtractor() func(context.Context) (ids.SessionId, bool) {
	return func(_ context.Context) (ids.SessionId, bool) {
		return ids.SessionId{}, false
	}
}

func ValidUserIdExtractor() func(context.Context) (ids.UserId, bool) {
	return func(_ context.Context) (ids.UserId, bool) {
		id, _ := ids.NewUserIdFromString(UserId)
		return id, true
	}
}

func MissingUserIdExtractor() func(context.Context) (ids.UserId, bool) {
	return func(_ context.Context) (ids.UserId, bool) {
		return ids.UserId{}, false
	}
}

func ValidAuthorExtractor() func(context.Context) (*authors.UserAuthor, bool) {
	return func(_ context.Context) (*authors.UserAuthor, bool) {
		id, _ := ids.NewUserIdFromString(UserId)
		adminId, _ := ids.NewUserIdFromString(AdminID)
		return &authors.UserAuthor{Id: id, AdminId: adminId, Role: vo.TenantRole}, true
	}
}

func MissingAuthorExtractor() func(context.Context) (*authors.UserAuthor, bool) {
	return func(_ context.Context) (*authors.UserAuthor, bool) {
		return nil, false
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

func ValidSubscriptionIdExtractor() func(context.Context) (ids.SubscriptionId, bool) {
	return func(_ context.Context) (ids.SubscriptionId, bool) {
		id, _ := ids.NewSubscriptionIdFromString(SubscriptionID)
		return id, true
	}
}

func MissingSubscriptionIdExtractor() func(context.Context) (ids.SubscriptionId, bool) {
	return func(_ context.Context) (ids.SubscriptionId, bool) {
		return ids.SubscriptionId{}, false
	}
}

func NewScopeIdExtractor() func(context.Context) (ids.ScopeId, bool) {
	return func(_ context.Context) (ids.ScopeId, bool) {
		return ids.NewScopeId(), true
	}
}

func ValidUserOwnershipExtractor() func(context.Context) (tenants.Ownership, bool) {
	return func(_ context.Context) (tenants.Ownership, bool) {
		userId, _ := ids.NewUserIdFromString(UserId)
		return tenants.NewOwnership(userId), true
	}
}

func MissingUserOwnershipExtractor() func(context.Context) (tenants.Ownership, bool) {
	return func(_ context.Context) (tenants.Ownership, bool) {
		return tenants.Ownership{}, false
	}
}

func ValidPrOwnershipExtractor() func(context.Context) (protected_resources.Ownership, bool) {
	return func(_ context.Context) (protected_resources.Ownership, bool) {
		subscriptionId, _ := ids.NewSubscriptionIdFromString(SubscriptionID)
		return protected_resources.NewOwnership(subscriptionId), true
	}
}

func MissingPrOwnershipExtractor() func(context.Context) (protected_resources.Ownership, bool) {
	return func(_ context.Context) (protected_resources.Ownership, bool) {
		return protected_resources.Ownership{}, false
	}
}
