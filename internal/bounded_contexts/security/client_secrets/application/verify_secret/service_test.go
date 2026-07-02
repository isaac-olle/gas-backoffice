package verify_secret_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets/application/verify_secret"
)

func TestVerifySecretService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		clientId    string
		plainSecret string
		verifierErr error
		wantErr     bool
	}{
		{
			name:        "returns nil when secret is valid",
			clientId:    "client-id",
			plainSecret: "plain-secret",
		},
		{
			name:        "returns error when secret is invalid",
			clientId:    "client-id",
			plainSecret: "wrong-secret",
			verifierErr: errors.New("invalid secret"),
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			verifier := &mockSecretVerifier{err: tc.verifierErr}
			svc := verify_secret.NewService(verifier)

			result, err := svc.VerifySecret(context.Background(), verify_secret.NewQuery(tc.clientId, tc.plainSecret))

			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Nil(t, result)
		})
	}
}
