package verify_assertion_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/verify_assertion"
)

func TestVerifyClientAssertionService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		assertion      string
		verifierClaims map[string]any
		verifierErr    error
		wantErr        bool
		wantClientID   string
	}{
		{
			name:           "extracts client id from valid assertion claims",
			assertion:      "some.jwt",
			verifierClaims: map[string]any{"sub": "client-abc"},
			wantClientID:   "client-abc",
		},
		{
			name:           "returns empty client id when sub claim is missing",
			assertion:      "tok",
			verifierClaims: map[string]any{},
			wantClientID:   "",
		},
		{
			name:        "returns error when verifier rejects assertion",
			assertion:   "bad.jwt",
			verifierErr: errors.New("bad signature"),
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			verifier := &mockAssertionVerifier{claims: tc.verifierClaims, err: tc.verifierErr}
			svc := verify_assertion.NewService(verifier)

			result, err := svc.VerifyAssertion(context.Background(), verify_assertion.NewQuery(tc.assertion))

			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			claims, ok := result.(*verify_assertion.ClientClaims)
			require.True(t, ok)
			assert.Equal(t, tc.wantClientID, claims.ClientID)
		})
	}
}
