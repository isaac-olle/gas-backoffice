package domain_services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks/domain_services"
	"github.com/stretchr/testify/require"
)

type signingKeyRepositoryStub struct {
	keys []*jwks.SigningKeyProjection
	err  error
}

func (r *signingKeyRepositoryStub) SearchActive(context.Context) ([]*jwks.SigningKeyProjection, error) {
	return r.keys, r.err
}

func TestSigningKeyFinderFindActive(t *testing.T) {
	keys := []*jwks.SigningKeyProjection{{}}
	tests := []struct {
		name    string
		keys    []*jwks.SigningKeyProjection
		repoErr error
		wantErr bool
	}{
		{name: "returns active signing keys", keys: keys},
		{name: "returns empty result when no keys are active", keys: []*jwks.SigningKeyProjection{}},
		{name: "wraps repository errors", repoErr: errors.New("query failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain_services.NewSigningKeyFinder(&signingKeyRepositoryStub{keys: tt.keys, err: tt.repoErr}).FindActive(context.Background())
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				require.Equal(t, tt.keys, got)
			}
		})
	}
}
