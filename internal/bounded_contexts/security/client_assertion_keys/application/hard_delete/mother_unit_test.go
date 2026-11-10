package hard_delete_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// validPubKeyHex is a valid PEM-encoded RSA 2048-bit PKIX public key for tests.
const validPubKeyHex = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArRj90RkbqGVRcp1xYi63
L+ZKaqMi3IVFhm6adSM+SNJsdiEkcMjrXu8iKlxN19kfJFN/fDlKwZuDgdKLq2/X
QqKpQHItF6kgNbuzfkNWbLkzebWn5sokC9BqgStH4Inob5Plg8T4v2Po7JIn8098
bGtxNPLObG/+heX7JQNZDyocWpQWD4G4BrvgimfKGH4azvBmLHR4dayCrii7GR2I
W1kY4MmhX84TNLCqFnzZdd0GswyGJHVnyrpUvEJFaHmr+97I8vDBvFPJ2w2dEXN1
+BC5dJHqmLjr6qV+6Dqq56Gn8DcestRxtiS20np4Z2cMSCgqtgp1l1zHSjMO/rGX
WQIDAQAB
-----END PUBLIC KEY-----
`

// keyFixture holds the ids of a key owned by a client under a subscription.
type keyFixture struct {
	id             ids.ClientAssertionKeyId
	clientId       ids.ClientId
	subscriptionId ids.SubscriptionId
}

func newKeyFixture(t *testing.T) keyFixture {
	t.Helper()
	id, err := ids.NewClientAssertionKeyIdFromString(uuid.New().String())
	require.NoError(t, err)
	clientId, err := ids.NewClientIdFromString(uuid.New().String())
	require.NoError(t, err)
	return keyFixture{id: id, clientId: clientId, subscriptionId: ids.NewSubscriptionId()}
}

func aKey(t *testing.T, f keyFixture) *client_assertion_keys.ClientAssertionKey {
	t.Helper()
	cak, err := client_assertion_keys.NewClientAssertionKeyForSubscription(f.id, f.clientId, f.subscriptionId, "RS256", validPubKeyHex)
	require.NoError(t, err)
	cak.PullDomainEvents()
	return cak
}

func authorFor(t *testing.T, subscriptionId ids.SubscriptionId) *vo.UserAuthor {
	t.Helper()
	author, err := vo.NewUserAuthorWithMemberships(ids.NewUserId().String(), ids.NewUserId().String(), vo.SubscriptionUserRole,
		[]vo.SubscriptionMembership{{SubscriptionId: subscriptionId, SubscriptionRole: vo.UserSubscriptionRole}})
	require.NoError(t, err)
	return author
}
