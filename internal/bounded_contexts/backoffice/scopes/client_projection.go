package scopes

import "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

type ClientProjection struct {
	ClientId       ids.ClientId
	SubscriptionId ids.SubscriptionId
}
