package types

type AggregateName string

const (
	AggregateAdmin             AggregateName = "admin"
	AggregateUser              AggregateName = "user"
	AggregateProtectedResource AggregateName = "protected-resource"
	AggregateScope             AggregateName = "scope"
	AggregateClient            AggregateName = "client"

	AggregateClientAssertionKey AggregateName = "client-assertion-key"
	AggregateAdminAssertionKey  AggregateName = "admin-assertion-key"
	AggregateKeyPair            AggregateName = "key-pair"
	AggregateClientSecret       AggregateName = "client-secret"

	AggregateAuthCodeSession          AggregateName = "auth-code-session"
	AggregateClientCredentialsSession AggregateName = "client-credentials-session"
)

func (a AggregateName) String() string {
	return string(a)
}
