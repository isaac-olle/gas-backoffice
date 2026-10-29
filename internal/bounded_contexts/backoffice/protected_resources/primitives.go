package protected_resources

import "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

type ProtectedResourcePrimitives struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	SubscriptionId string   `json:"subscription_id"`
	IsSystem       bool     `json:"is_system"`
	Enabled        bool     `json:"enabled"`
	ScopeIds       []string `json:"scope_ids"`
	Version        int      `json:"version"`
	AuthorId       string   `json:"author_id"`
	AuthorRole     string   `json:"author_role"`
}

type OwnershipPrimitives struct {
	SubscriptionId string `json:"subscription_id"`
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{SubscriptionId: o.SubscriptionId.String()}
}

func (pr *ProtectedResource) ToPrimitives() ProtectedResourcePrimitives {
	scopeIds := make([]string, 0, pr.scopeIds.Size())
	for _, sc := range pr.scopeIds.Values() {
		scopeIds = append(scopeIds, sc.String())
	}

	return ProtectedResourcePrimitives{
		Id:             pr.id.String(),
		Name:           pr.name.String(),
		SubscriptionId: pr.ownership.SubscriptionId.String(),
		IsSystem:       pr.isSystem,
		Enabled:        pr.enabled,
		ScopeIds:       scopeIds,
		Version:        pr.version,
		AuthorId:       authorIDString(pr.authorId),
		AuthorRole:     pr.authorRole.String(),
	}
}

func authorIDString(id ids.UserId) string {
	if id == (ids.UserId{}) {
		return ""
	}
	return id.String()
}
