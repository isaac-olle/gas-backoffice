package domain_services_test

import (
	"context"
	"errors"
	app_errors2 "github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
)

// --- stubs ---

type mockSessionRepo struct {
	byId   *session.AuthorizationCodeSession
	byCode *session.AuthorizationCodeSession
	err    error
}

func (m *mockSessionRepo) Save(_ context.Context, _ *session.AuthorizationCodeSession) error {
	return nil
}
func (m *mockSessionRepo) FindById(_ context.Context, _ ids.SessionId) (*session.AuthorizationCodeSession, error) {
	return m.byId, m.err
}
func (m *mockSessionRepo) FindByCode(_ context.Context, _ string) (*session.AuthorizationCodeSession, error) {
	return m.byCode, m.err
}
func (m *mockSessionRepo) SoftDelete(_ context.Context, _ ids.SessionId) error { return nil }
func (m *mockSessionRepo) HardDelete(_ context.Context, _ ids.SessionId) error { return nil }

type mockClientRepo struct {
	result *session.Client
	err    error
}

func (m *mockClientRepo) FindById(_ context.Context, _ string) (*session.Client, error) {
	return m.result, m.err
}

type codeGen struct{}

func (c *codeGen) Generate() (string, error) { return "gen-code", nil }

func makeSession(t *testing.T) *session.AuthorizationCodeSession {
	t.Helper()
	s, err := session.NewAuthorizationCodeSession(
		ids.NewSessionId(), uuid.New().String(),
		"https://example.com/cb", "openid", "state",
		"code", &codeGen{}, "", "",
	)
	require.NoError(t, err)
	return s
}

// --- SessionFinder tests ---

func TestSessionFinder_FindById(t *testing.T) {
	t.Parallel()

	sess := makeSession(t)

	tests := []struct {
		name        string
		repoResult  *session.AuthorizationCodeSession
		repoErr     error
		wantErr     bool
		wantErrType app_errors2.ErrorType
	}{
		{
			name:       "returns session when found",
			repoResult: sess,
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
			repo := &mockSessionRepo{byId: tc.repoResult, err: tc.repoErr}
			finder := domain_services.NewSessionFinder(repo)

			found, err := finder.FindById(context.Background(), ids.NewSessionId())
			if tc.wantErr {
				require.Error(t, err)
				var appErr *app_errors2.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tc.wantErrType, appErr.Type)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, sess, found)
		})
	}
}

func TestSessionFinder_FindByCode(t *testing.T) {
	t.Parallel()

	sess := makeSession(t)

	tests := []struct {
		name        string
		repoResult  *session.AuthorizationCodeSession
		repoErr     error
		wantErr     bool
		wantErrType app_errors2.ErrorType
	}{
		{
			name:       "returns session when code matches",
			repoResult: sess,
		},
		{
			name:        "returns bad request when code is not found",
			repoResult:  nil,
			wantErr:     true,
			wantErrType: app_errors2.BadRequestErrorType,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockSessionRepo{byCode: tc.repoResult, err: tc.repoErr}
			finder := domain_services.NewSessionFinder(repo)

			found, err := finder.FindByCode(context.Background(), "some-code")
			if tc.wantErr {
				require.Error(t, err)
				var appErr *app_errors2.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tc.wantErrType, appErr.Type)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, sess, found)
		})
	}
}

// --- ClientFinder tests ---

func TestClientFinder_Find(t *testing.T) {
	t.Parallel()

	client := session.NewClient(
		uuid.New().String(), "TestClient",
		uuid.New().String(),
		[]string{}, []string{}, []*session.Scope{},
	)

	tests := []struct {
		name        string
		repoResult  *session.Client
		repoErr     error
		wantErr     bool
		wantErrType app_errors2.ErrorType
	}{
		{
			name:       "returns client when found",
			repoResult: client,
		},
		{
			name:        "returns not found when repo returns nil",
			repoResult:  nil,
			wantErr:     true,
			wantErrType: app_errors2.NotFoundErrorType,
		},
		{
			name:        "wraps repository error as internal error",
			repoErr:     errors.New("timeout"),
			wantErr:     true,
			wantErrType: app_errors2.InternalErrorType,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockClientRepo{result: tc.repoResult, err: tc.repoErr}
			finder := domain_services.NewClientFinder(repo)

			found, err := finder.Find(context.Background(), uuid.New().String())
			if tc.wantErr {
				require.Error(t, err)
				var appErr *app_errors2.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tc.wantErrType, appErr.Type)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, client, found)
		})
	}
}
