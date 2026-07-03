package verify_assertion_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// stubPubKeyHex is a valid hex-encoded RSA 2048-bit PKIX public key for tests.
const stubPubKeyHex = "30820122300d06092a864886f70d01010105000382010f003082010a0282010100ad18fdd1191ba86551729d71622eb72fe64a6aa322dc8545866e9a75233e48d26c76212470c8eb5eef222a5c4dd7d91f24537f7c394ac19b8381d28bab6fd742a2a940722d17a92035bbb37e43566cb93379b5a7e6ca240bd06a812b47e089e86f93e583c4f8bf63e8ec9227f34f7c6c6b7134f2ce6c6ffe85e5fb2503590f2a1c5a94160f81b806bbe08a67ca187e1acef0662c747875ac82ae28bb191d885b5918e0c9a15fce1334b0aa167cd975dd06b30c86247567caba54bc42456879abfbdec8f2f0c1bc53c9db0d9d117375f810b97491ea98b8ebeaa57ee83aaae7a1a7f0371eb2d471b624b6d27a7867670c48282ab60a75975cc74a330efeb197590203010001"

// mockAssertionVerifier stubs admin_assertion_keys.AssertionVerifier.
type mockAssertionVerifier struct {
	claims       map[string]any
	err          error
	gotAssertion string
	gotKey       *admin_assertion_keys.AdminAssertionKey
}

func (m *mockAssertionVerifier) Verify(_ context.Context, assertion string, key *admin_assertion_keys.AdminAssertionKey) (map[string]any, error) {
	m.gotAssertion = assertion
	m.gotKey = key
	return m.claims, m.err
}

// mockJwtDecoder stubs admin_assertion_keys.JwtDecoder.
type mockJwtDecoder struct {
	kid          string
	err          error
	gotAssertion string
}

func (m *mockJwtDecoder) ExtractKeyId(_ context.Context, assertion string) (string, error) {
	m.gotAssertion = assertion
	return m.kid, m.err
}

// newAdminAssertionKeyMother builds a valid AdminAssertionKey for tests.
func newAdminAssertionKeyMother(t *testing.T) *admin_assertion_keys.AdminAssertionKey {
	t.Helper()
	adminId, err := ids2.NewAdminIdFromString(uuid.New().String())
	require.NoError(t, err)
	key, err := admin_assertion_keys.NewAdminAssertionKey(
		uuid.New().String(), adminId, "RS256", stubPubKeyHex,
	)
	require.NoError(t, err)
	return key
}
