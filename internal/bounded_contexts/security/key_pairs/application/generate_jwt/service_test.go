package generate_jwt_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/domain_services"
	kpmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestGenerateJWT_HappyPath(t *testing.T) {
	t.Parallel()

	kp := newKeyPairMother(t)
	claims := map[string]any{"sub": "user1"}

	repo := kpmocks.NewMockRepository(t)
	repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(kp, nil)
	gen := &mockTokenGenerator{token: "signed.jwt.token"}
	finder := domain_services.NewKeyPairFinder(repo)
	svc := generate_jwt.NewService(finder, gen)

	qry := generate_jwt.NewQuery(kp.ToPrimitives().Id, claims, time.Hour)
	result, err := svc.GenerateJWT(context.Background(), qry)

	require.NoError(t, err)
	assert.Equal(t, "signed.jwt.token", result)
	assert.Equal(t, kp, gen.gotKeyPair, "generator must receive the fetched KeyPair")
	assert.Equal(t, claims, gen.gotClaims, "generator must receive the original claims")
}

func TestGenerateJWT_ErrorCases(t *testing.T) {
	t.Parallel()

	kp := newKeyPairMother(t)

	tests := []struct {
		name         string
		keyPairId    string
		repoResult   *key_pairs.KeyPair
		generatorErr error
		wantErrCode  app_errors.ErrorCode
	}{
		{
			name:        "returns not found error for invalid key pair id",
			keyPairId:   "not-a-uuid",
			repoResult:  nil,
			wantErrCode: key_pairs.ErrKeyPairNotFound,
		},
		{
			name:        "returns not found error when key pair does not exist",
			keyPairId:   uuid.New().String(),
			repoResult:  nil,
			wantErrCode: key_pairs.ErrKeyPairNotFound,
		},
		{
			name:         "returns not found error when generator fails",
			keyPairId:    kp.ToPrimitives().Id,
			repoResult:   kp,
			generatorErr: errors.New("sign failed"),
			wantErrCode:  key_pairs.ErrKeyPairNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := kpmocks.NewMockRepository(t)
			if tc.keyPairId != "not-a-uuid" {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(tc.repoResult, nil).Maybe()
			}
			gen := &mockTokenGenerator{err: tc.generatorErr}
			finder := domain_services.NewKeyPairFinder(repo)
			svc := generate_jwt.NewService(finder, gen)

			qry := generate_jwt.NewQuery(tc.keyPairId, map[string]any{"sub": "user1"}, time.Hour)
			_, err := svc.GenerateJWT(context.Background(), qry)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
