package types

type AggregateName string

const (
	AggregateAdmin             AggregateName = "admin"
	AggregateUser              AggregateName = "user"
	AggregateProtectedResource AggregateName = "protected-resource"
	AggregateScope             AggregateName = "scope"
	AggregateClient            AggregateName = "client"
	AggregateSession           AggregateName = "session"

	AggregateClientAssertionKey AggregateName = "client-assertion-key"
	AggregateKeyPair            AggregateName = "key-pair"
	AggregateClientSecret       AggregateName = "client-secret"

	AggregateIdentityProvider AggregateName = "identity-provider"

	AggregateAuthCodeSession          AggregateName = "auth-code-session"
	AggregateClientCredentialsSession AggregateName = "client-credentials-session"

	AggregateUserIdentity AggregateName = "user-identity"
)

func (a AggregateName) String() string {
	return string(a)
}
