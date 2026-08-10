package handlerstest

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

const (
	AdminID   = "b1b2c3d4-e5f6-7890-abcd-ef1234567890"
	UserId    = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	ClientID  = "c1b2c3d4-e5f6-7890-abcd-ef1234567890"
	PrID      = "d1b2c3d4-e5f6-7890-abcd-ef1234567890"
	ScopeID   = "e1b2c3d4-e5f6-7890-abcd-ef1234567890"
	SessionID = "f1b2c3d4-e5f6-7890-abcd-ef1234567890"
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

func ValidUserOwnershipExtractor() func(context.Context) (users.Ownership, bool) {
	return func(_ context.Context) (users.Ownership, bool) {
		userId, _ := ids.NewUserIdFromString(UserId)
		return users.NewOwnership(userId), true
	}
}

func MissingUserOwnershipExtractor() func(context.Context) (users.Ownership, bool) {
	return func(_ context.Context) (users.Ownership, bool) {
		return users.Ownership{}, false
	}
}

func ValidPrOwnershipExtractor() func(context.Context) (protected_resources.Ownership, bool) {
	return func(_ context.Context) (protected_resources.Ownership, bool) {
		userId, _ := ids.NewUserIdFromString(UserId)
		return protected_resources.NewOwnership(userId), true
	}
}

func MissingPrOwnershipExtractor() func(context.Context) (protected_resources.Ownership, bool) {
	return func(_ context.Context) (protected_resources.Ownership, bool) {
		return protected_resources.Ownership{}, false
	}
}
