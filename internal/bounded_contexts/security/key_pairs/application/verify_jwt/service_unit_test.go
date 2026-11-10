package verify_jwt_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/domain_services"
	kpmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

const testIssuer = "gas"

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestVerifyJWT_HappyPath(t *testing.T) {
	t.Parallel()

	kp := newKeyPairMother(t)
	validClaims := newValidClaimsMother(testIssuer)
	const token = "some.jwt.token"

	repo := kpmocks.NewMockRepository(t)
	repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(kp, nil)
	decoder := &mockJwtDecoder{kid: kp.ToPrimitives().Id}
	verifier := &mockTokenVerifier{claims: validClaims}
	finder := domain_services.NewKeyPairFinder(repo)
	svc := verify_jwt.NewService(finder, verifier, decoder)

	result, err := svc.VerifyJWT(context.Background(), verify_jwt.NewQuery(token).WithExpectedClaim("iss", testIssuer).WithVerifyExpiration())

	require.NoError(t, err)
	assert.Equal(t, validClaims, result)
	assert.Equal(t, token, decoder.gotToken, "decoder must receive the original token")
	assert.Equal(t, token, verifier.gotToken, "verifier must receive the original token")
	assert.Equal(t, kp, verifier.gotKeyPair, "verifier must receive the fetched KeyPair")
	assert.Equal(t, testIssuer, verifier.gotExpectedClaims["iss"], "verifier must receive the expected claims")
	assert.True(t, verifier.gotVerifyExpiration, "verifier must receive the verify-expiration flag")
}

func TestVerifyJWT_ClaimsToReturn(t *testing.T) {
	t.Parallel()

	kp := newKeyPairMother(t)
	claims := map[string]any{
		"sub": uuid.New().String(),
		"address": map[string]any{
			"city": "Barcelona",
		},
	}

	tests := []struct {
		name           string
		claimsTemplate map[string]any
		wantResult     map[string]any
		wantErrCode    app_errors.ErrorCode
	}{
		{
			name:           "returns only the requested top-level claim",
			claimsTemplate: map[string]any{"sub": nil},
			wantResult:     map[string]any{"sub": claims["sub"]},
		},
		{
			name:           "descends into a nested claim",
			claimsTemplate: map[string]any{"address": map[string]any{"city": nil}},
			wantResult:     map[string]any{"address": map[string]any{"city": "Barcelona"}},
		},
		{
			name:           "returns invalid token error when a requested claim is missing",
			claimsTemplate: map[string]any{"missing": nil},
			wantErrCode:    key_pairs.ErrInvalidToken,
		},
		{
			name:           "returns invalid token error when descending into a non-object claim",
			claimsTemplate: map[string]any{"sub": map[string]any{"nested": nil}},
			wantErrCode:    key_pairs.ErrInvalidToken,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := kpmocks.NewMockRepository(t)
			repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(kp, nil)
			decoder := &mockJwtDecoder{kid: kp.ToPrimitives().Id}
			verifier := &mockTokenVerifier{claims: claims}
			finder := domain_services.NewKeyPairFinder(repo)
			svc := verify_jwt.NewService(finder, verifier, decoder)

			result, err := svc.VerifyJWT(context.Background(), verify_jwt.NewQuery("tok").WithClaimsToReturn(tc.claimsTemplate))

			if tc.wantErrCode != "" {
				assertErrCode(t, err, tc.wantErrCode)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestVerifyJWT_ErrorCases(t *testing.T) {
	t.Parallel()

	kp := newKeyPairMother(t)
	validClaims := newValidClaimsMother(testIssuer)

	tests := []struct {
		name           string
		token          string
		decoderKid     string
		decoderErr     error
		repoResult     *key_pairs.KeyPair
		verifierClaims map[string]any
		verifierErr    error
		wantErrCode    app_errors.ErrorCode
	}{
		{
			name:        "returns invalid token error when decoder fails",
			token:       "bad.tok",
			decoderErr:  errors.New("malformed token"),
			wantErrCode: key_pairs.ErrInvalidToken,
		},
		{
			name:        "returns invalid token error when kid is not a valid UUID",
			token:       "tok",
			decoderKid:  "not-a-uuid",
			wantErrCode: key_pairs.ErrInvalidToken,
		},
		{
			name:        "returns not found error when key pair does not exist",
			token:       "tok",
			decoderKid:  uuid.New().String(),
			repoResult:  nil,
			wantErrCode: key_pairs.ErrKeyPairNotFound,
		},
		{
			name:           "returns invalid token error when verifier rejects the token",
			token:          "bad.sig",
			decoderKid:     kp.ToPrimitives().Id,
			repoResult:     kp,
			verifierClaims: validClaims,
			verifierErr:    errors.New("invalid signature"),
			wantErrCode:    key_pairs.ErrInvalidToken,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := kpmocks.NewMockRepository(t)
			if tc.decoderErr == nil && tc.decoderKid != "not-a-uuid" {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(tc.repoResult, nil).Maybe()
			}
			decoder := &mockJwtDecoder{kid: tc.decoderKid, err: tc.decoderErr}
			verifier := &mockTokenVerifier{claims: tc.verifierClaims, err: tc.verifierErr}
			finder := domain_services.NewKeyPairFinder(repo)
			svc := verify_jwt.NewService(finder, verifier, decoder)

			_, err := svc.VerifyJWT(context.Background(), verify_jwt.NewQuery(tc.token))

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
