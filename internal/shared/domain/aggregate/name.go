package aggregate

import (
	sharedtypes "github.com/isaac-olle/gas-backoffice/internal/shared/types"
)

type Name = sharedtypes.AggregateName

const (
	Admin             = sharedtypes.AggregateAdmin
	User              = sharedtypes.AggregateUser
	ProtectedResource = sharedtypes.AggregateProtectedResource
	Scope             = sharedtypes.AggregateScope
	Client            = sharedtypes.AggregateClient

	ClientAssertionKey = sharedtypes.AggregateClientAssertionKey
	AdminAssertionKey  = sharedtypes.AggregateAdminAssertionKey
	KeyPair            = sharedtypes.AggregateKeyPair
	ClientSecret       = sharedtypes.AggregateClientSecret

	AuthCodeSession          = sharedtypes.AggregateAuthCodeSession
	ClientCredentialsSession = sharedtypes.AggregateClientCredentialsSession
)
