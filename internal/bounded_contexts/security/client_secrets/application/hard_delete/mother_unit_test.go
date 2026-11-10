package hard_delete_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// secretFixture holds the ids of a secret owned by a client under a subscription.
type secretFixture struct {
	id             ids.ClientSecretId
	clientId       ids.ClientId
	subscriptionId ids.SubscriptionId
}

func newSecretFixture(t *testing.T) secretFixture {
	t.Helper()
	id, err := ids.NewClientSecretIdFromString(uuid.New().String())
	require.NoError(t, err)
	clientId, err := ids.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	return secretFixture{id: id, clientId: clientId, subscriptionId: ids.NewSubscriptionId()}
}

func aSecret(t *testing.T, f secretFixture) *client_secrets.ClientSecret {
	t.Helper()
	cs, err := client_secrets.FromPrimitives(client_secrets.ClientSecretPrimitives{
		Id:             f.id.String(),
		ClientId:       f.clientId.String(),
		SubscriptionId: f.subscriptionId.String(),
		SecretHash:     "hash",
	})
	require.NoError(t, err)
	return cs
}

func authorFor(t *testing.T, subscriptionId ids.SubscriptionId) *vo.UserAuthor {
	t.Helper()
	author, err := vo.NewUserAuthorWithMemberships(ids.NewUserId().String(), ids.NewUserId().String(), vo.SubscriptionUserRole,
		[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
	require.NoError(t, err)
	return author
}
