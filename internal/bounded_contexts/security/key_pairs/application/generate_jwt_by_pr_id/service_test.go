package generate_jwt_by_pr_id_test

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
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt_by_pr_id"
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

func TestGenerateJWTByPrId_HappyPath(t *testing.T) {
	t.Parallel()

	kp := newKeyPairMother(t)
	prId := kp.ToPrimitives().ProtectedResourceId
	claims := map[string]any{"sub": "user1"}

	repo := kpmocks.NewMockRepository(t)
	repo.EXPECT().SearchByProtectedResourceId(mock.Anything, mock.Anything).Return(kp, nil)
	gen := &mockTokenGenerator{token: "signed.jwt.token"}
	finder := domain_services.NewKeyPairFinder(repo)
	svc := generate_jwt_by_pr_id.NewService(finder, gen)

	qry := generate_jwt_by_pr_id.NewQuery(prId, claims, time.Hour)
	result, err := svc.GenerateJWT(context.Background(), qry)

	require.NoError(t, err)
	assert.Equal(t, "signed.jwt.token", result)
	assert.Equal(t, kp, gen.gotKeyPair)
	assert.Equal(t, claims, gen.gotClaims)
}

func TestGenerateJWTByPrId_ErrorCases(t *testing.T) {
	t.Parallel()

	kp := newKeyPairMother(t)
	prId := kp.ToPrimitives().ProtectedResourceId

	tests := []struct {
		name               string
		protectedResoureId string
		repoResult         *key_pairs.KeyPair
		generatorErr       error
		wantErrCode        app_errors.ErrorCode
	}{
		{
			name:               "returns error for invalid protected resource id",
			protectedResoureId: "not-a-uuid",
			wantErrCode:        key_pairs.ErrKeyPairNotFound,
		},
		{
			name:               "returns not found when no key pair exists for the PR",
			protectedResoureId: uuid.New().String(),
			repoResult:         nil,
			wantErrCode:        key_pairs.ErrKeyPairNotFound,
		},
		{
			name:               "returns error when token generator fails",
			protectedResoureId: prId,
			repoResult:         kp,
			generatorErr:       errors.New("sign failed"),
			wantErrCode:        key_pairs.ErrKeyPairNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := kpmocks.NewMockRepository(t)
			if tc.protectedResoureId != "not-a-uuid" {
				repo.EXPECT().SearchByProtectedResourceId(mock.Anything, mock.Anything).Return(tc.repoResult, nil).Maybe()
			}
			gen := &mockTokenGenerator{err: tc.generatorErr}
			finder := domain_services.NewKeyPairFinder(repo)
			svc := generate_jwt_by_pr_id.NewService(finder, gen)

			qry := generate_jwt_by_pr_id.NewQuery(tc.protectedResoureId, map[string]any{"sub": "user1"}, time.Hour)
			_, err := svc.GenerateJWT(context.Background(), qry)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
