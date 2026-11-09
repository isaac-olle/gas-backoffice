package domain_services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/domain_services"
	keyPairInfra "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/infrastructure"
	kpmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/mocks"
	app_errors2 "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func newKeyPairMother(t *testing.T) *key_pairs.KeyPair {
	t.Helper()
	privateKeyPEM, publicKeyPEM, err := keyPairInfra.NewPEMKeyGenerator().Generate("RS256")
	require.NoError(t, err)
	kp, err := key_pairs.NewKeyPair(
		uuid.New().String(), uuid.New().String(), "RS256", privateKeyPEM, publicKeyPEM,
	)
	require.NoError(t, err)
	return kp
}

func TestKeyPairFinder_Find(t *testing.T) {
	t.Parallel()

	kp := newKeyPairMother(t)

	tests := []struct {
		name        string
		repoResult  *key_pairs.KeyPair
		repoErr     error
		wantErr     bool
		wantErrType app_errors2.ErrorType
	}{
		{
			name:       "finds existing key pair",
			repoResult: kp,
		},
		{
			name:        "returns not found when repo returns nil",
			repoResult:  nil,
			wantErr:     true,
			wantErrType: app_errors2.NotFoundErrorType,
		},
		{
			name:        "wraps repository error as internal error",
			repoErr:     errors.New("db error"),
			wantErr:     true,
			wantErrType: app_errors2.InternalErrorType,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := kpmocks.NewMockRepository(t)
			repo.EXPECT().SearchById(mock.Anything, mock.AnythingOfType("ids.KeyPairId")).Return(tc.repoResult, tc.repoErr)

			finder := domain_services.NewKeyPairFinder(repo)

			found, err := finder.Find(context.Background(), ids.NewKeyPairId())

			if tc.wantErr {
				require.Error(t, err)
				var appErr *app_errors2.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tc.wantErrType, appErr.Type)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, kp, found)
		})
	}
}
