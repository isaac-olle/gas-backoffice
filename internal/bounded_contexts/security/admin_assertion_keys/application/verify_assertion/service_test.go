package verify_assertion_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys/application/verify_assertion"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys/domain_services"
	aakmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

const jwtBearerGrant = "urn:ietf:params:oauth:grant-type:jwt-bearer"

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestVerifyAdminAssertionService_HappyPath(t *testing.T) {
	t.Parallel()

	key := newAdminAssertionKeyMother(t)
	p := key.ToPrimitives()
	const assertion = "valid.assertion.token"
	validClaims := map[string]any{
		"sub": uuid.New().String(),
		"iss": p.AdminId,
		"aud": "gas",
	}

	repo := aakmocks.NewMockRepository(t)
	repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(key, nil)
	finder := domain_services.NewAdminAssertionKeyFinder(repo)
	verifier := &mockAssertionVerifier{claims: validClaims}
	decoder := &mockJwtDecoder{kid: p.Id}
	svc := verify_assertion.NewService(finder, verifier, decoder, "gas")

	_, err := svc.VerifyAssertion(context.Background(), verify_assertion.NewQuery(assertion, jwtBearerGrant))

	require.NoError(t, err)
	assert.Equal(t, assertion, decoder.gotAssertion, "decoder must receive the original assertion")
	assert.Equal(t, assertion, verifier.gotAssertion, "verifier must receive the original assertion")
	assert.Equal(t, key, verifier.gotKey, "verifier must receive the fetched key")
}

func TestVerifyAdminAssertionService_ErrorCases(t *testing.T) {
	t.Parallel()

	key := newAdminAssertionKeyMother(t)

	tests := []struct {
		name           string
		assertion      string
		grantType      string
		decoderKid     string
		decoderErr     error
		repoResult     *admin_assertion_keys.AdminAssertionKey
		verifierClaims map[string]any
		verifierErr    error
		wantErrCode    app_errors.ErrorCode
	}{
		{
			name:        "rejects unsupported grant type",
			assertion:   "tok",
			grantType:   "client_credentials",
			wantErrCode: "UNSUPPORTED_GRANT_TYPE",
		},
		{
			name:        "returns invalid assertion error when decoder fails",
			assertion:   "tok",
			grantType:   jwtBearerGrant,
			decoderErr:  errors.New("bad token"),
			wantErrCode: admin_assertion_keys.ErrInvalidAdminAssertion,
		},
		{
			name:        "returns not found error when key does not exist",
			assertion:   "tok",
			grantType:   jwtBearerGrant,
			decoderKid:  uuid.New().String(),
			repoResult:  nil,
			wantErrCode: admin_assertion_keys.ErrAdminAssertionKeyNotFound,
		},
		{
			name:        "returns invalid assertion error when verifier rejects signature",
			assertion:   "tok",
			grantType:   jwtBearerGrant,
			decoderKid:  key.ToPrimitives().Id,
			repoResult:  key,
			verifierErr: errors.New("invalid sig"),
			wantErrCode: admin_assertion_keys.ErrInvalidAdminAssertion,
		},
		{
			name:        "returns invalid assertion error when kid is not a valid UUID",
			assertion:   "tok",
			grantType:   jwtBearerGrant,
			decoderKid:  "not-a-uuid",
			wantErrCode: admin_assertion_keys.ErrInvalidAdminAssertion,
		},
		{
			name:       "returns bad request when claims are missing required fields",
			assertion:  "tok",
			grantType:  jwtBearerGrant,
			decoderKid: key.ToPrimitives().Id,
			repoResult: key,
			verifierClaims: map[string]any{
				"sub": 12345,
				"iss": key.ToPrimitives().AdminId,
				"aud": "gas",
			},
			wantErrCode: admin_assertion_keys.ErrInvalidAdminAssertion,
		},
		{
			name:       "returns invalid assertion when issuer does not match admin id",
			assertion:  "tok",
			grantType:  jwtBearerGrant,
			decoderKid: key.ToPrimitives().Id,
			repoResult: key,
			verifierClaims: map[string]any{
				"sub": uuid.New().String(),
				"iss": uuid.New().String(),
				"aud": "gas",
			},
			wantErrCode: admin_assertion_keys.ErrInvalidAdminAssertion,
		},
		{
			name:       "returns invalid assertion when audience does not match app name",
			assertion:  "tok",
			grantType:  jwtBearerGrant,
			decoderKid: key.ToPrimitives().Id,
			repoResult: key,
			verifierClaims: map[string]any{
				"sub": uuid.New().String(),
				"iss": key.ToPrimitives().AdminId,
				"aud": "wrong-app",
			},
			wantErrCode: admin_assertion_keys.ErrInvalidAdminAssertion,
		},
		{
			name:       "returns invalid assertion when sub is not a valid actor UUID",
			assertion:  "tok",
			grantType:  jwtBearerGrant,
			decoderKid: key.ToPrimitives().Id,
			repoResult: key,
			verifierClaims: map[string]any{
				"sub": "not-a-uuid",
				"iss": key.ToPrimitives().AdminId,
				"aud": "gas",
			},
			wantErrCode: admin_assertion_keys.ErrInvalidAdminAssertion,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := aakmocks.NewMockRepository(t)
			if tc.decoderErr == nil && tc.grantType == jwtBearerGrant && tc.decoderKid != "not-a-uuid" {
				repo.EXPECT().SearchById(mock.Anything, mock.Anything).Return(tc.repoResult, nil).Maybe()
			}
			finder := domain_services.NewAdminAssertionKeyFinder(repo)
			verifier := &mockAssertionVerifier{claims: tc.verifierClaims, err: tc.verifierErr}
			decoder := &mockJwtDecoder{kid: tc.decoderKid, err: tc.decoderErr}
			svc := verify_assertion.NewService(finder, verifier, decoder, "gas")

			_, err := svc.VerifyAssertion(context.Background(), verify_assertion.NewQuery(tc.assertion, tc.grantType))

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
