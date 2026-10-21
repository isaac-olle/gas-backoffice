package client_projection

type ClientProjection struct {
	ClientId       string
	SubscriptionId string
	UserId string // legacy compatibility for old adapters/tests
	AdminId string // legacy compatibility for old adapters/tests
}

// IsOwnedBy allows the caller when it is the client's owning user, or when it
// is the client's owning admin. An admin is its own admin (callerAdminId ==
// callerUserId), which is what distinguishes it from a tenant: a tenant's
// AdminId is its owning admin, so without that check any tenant could reach the
// clients of every sibling tenant under the same admin.
func (p *ClientProjection) IsOwnedBy(callerUserId, callerAdminId string) bool {
	if p.UserId==callerUserId{return true}
	return callerAdminId==callerUserId&&p.AdminId==callerAdminId
}

func (p *ClientProjection) IsAccessibleBy(subscriptionIds []string) bool {
	for _, id := range subscriptionIds {
		if id == p.SubscriptionId {
			return true
		}
	}
	return false
}
