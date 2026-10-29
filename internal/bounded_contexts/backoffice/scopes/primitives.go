package scopes

import "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

type ScopePrimitives struct {
	Id                  string   `json:"id"`
	Name                string   `json:"name"`
	ProtectedResourceId string   `json:"protected_resource_id"`
	SubscriptionId      string   `json:"subscription_id"`
	AllowedClientIds    []string `json:"allowed_client_ids"`
	IsSystem            bool     `json:"is_system"`
	Enabled             bool     `json:"enabled"`
	Version             int      `json:"version"`
	AuthorId            string   `json:"author_id"`
	AuthorRole          string   `json:"author_role"`
}

type OwnershipPrimitives struct {
	ProtectedResourceId string `json:"protected_resource_id"`
	SubscriptionId      string `json:"subscription_id"`
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{
		ProtectedResourceId: o.ProtectedResourceId.String(),
		SubscriptionId:      o.SubscriptionId.String(),
	}
}

func (s *Scope) ToPrimitives() ScopePrimitives {
	allowedClientIds := make([]string, 0, s.allowedClientIds.Size())
	for _, id := range s.allowedClientIds.Values() {
		allowedClientIds = append(allowedClientIds, id.String())
	}
	return ScopePrimitives{
		Id:                  s.id.String(),
		Name:                s.name.String(),
		ProtectedResourceId: s.ownership.ProtectedResourceId.String(),
		SubscriptionId:      s.ownership.SubscriptionId.String(),
		AllowedClientIds:    allowedClientIds,
		IsSystem:            s.isSystem,
		Enabled:             s.enabled,
		Version:             s.version,
		AuthorId:            authorIDString(s.authorId),
		AuthorRole:          s.authorRole.String(),
	}
}

func authorIDString(id ids.UserId) string {
	if id == (ids.UserId{}) {
		return ""
	}
	return id.String()
}
