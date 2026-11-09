package hard_delete_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// secretFixture holds the ids of a secret owned by a tenant (userId) under an admin.
type secretFixture struct {
	id       ids.ClientSecretId
	clientId ids.ClientId
	adminId  ids.UserId
	userId   ids.UserId
}

func newSecretFixture(t *testing.T) secretFixture {
	t.Helper()
	id, err := ids.NewClientSecretIdFromString(uuid.New().String())
	require.NoError(t, err)
	clientId, err := ids.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	return secretFixture{id: id, clientId: clientId, adminId: ids.NewUserId(), userId: ids.NewUserId()}
}

func aSecret(t *testing.T, f secretFixture) *client_secrets.ClientSecret {
	t.Helper()
	cs, err := client_secrets.FromPrimitives(client_secrets.ClientSecretPrimitives{
		Id:         f.id.String(),
		ClientId:   f.clientId.String(),
		UserId:     f.userId.String(),
		AdminId:    f.adminId.String(),
		SecretHash: "hash",
	})
	require.NoError(t, err)
	return cs
}
