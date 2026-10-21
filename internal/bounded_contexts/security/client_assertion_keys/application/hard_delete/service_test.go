package hard_delete_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/hard_delete"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/domain_services"
	cakmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/domain_services/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func newService(repo *cakmocks.MockClientAssertionKeyRepository) *hard_delete.Service {
	return hard_delete.NewService(domain_services.NewClientAssertionKeyFinder(repo), repo)
}

func TestHardDeleteClientAssertionKeyService_HappyPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		callerUserId  func(f keyFixture) string
		callerAdminId func(f keyFixture) string
	}{
		{
			name:          "owning tenant deletes its key",
			callerUserId:  func(f keyFixture) string { return f.userId.String() },
			callerAdminId: func(f keyFixture) string { return f.adminId.String() },
		},
		{
			name:          "owning admin deletes its tenant's key",
			callerUserId:  func(f keyFixture) string { return f.adminId.String() },
			callerAdminId: func(f keyFixture) string { return f.adminId.String() },
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := newKeyFixture(t)

			var deletedId ids.ClientAssertionKeyId
			repo := cakmocks.NewMockClientAssertionKeyRepository(t)
			repo.EXPECT().SearchById(mock.Anything, f.id).Return(aKey(t, f), nil)
			repo.EXPECT().HardDelete(mock.Anything, mock.Anything).Run(func(_ context.Context, id ids.ClientAssertionKeyId) {
				deletedId = id
			}).Return(nil)

			cmd := hard_delete.NewCommand(f.id.String(), f.clientId.String(), tc.callerAdminId(f), tc.callerUserId(f))
			require.NoError(t, newService(repo).HardDelete(context.Background(), cmd))

			assert.Equal(t, f.id, deletedId, "repo must receive the correct ID for deletion")
		})
	}
}

func TestHardDeleteClientAssertionKeyService_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		buildCmd    func(f keyFixture) *hard_delete.Command
		setupRepo   func(t *testing.T, repo *cakmocks.MockClientAssertionKeyRepository, f keyFixture)
		wantErrCode app_errors.ErrorCode
	}{
		{
			name: "rejects invalid id",
			buildCmd: func(f keyFixture) *hard_delete.Command {
				return hard_delete.NewCommand("not-a-valid-uuid", f.clientId.String(), f.adminId.String(), f.userId.String())
			},
			wantErrCode: client_assertion_keys.ErrClientAssertionKeyNotFound,
		},
		{
			name: "rejects invalid ownership ids",
			buildCmd: func(f keyFixture) *hard_delete.Command {
				return hard_delete.NewCommand(f.id.String(), f.clientId.String(), "not-a-uuid", f.userId.String())
			},
			wantErrCode: client_assertion_keys.ErrInvalidOwnership,
		},
		{
			name: "returns not found when key does not exist",
			buildCmd: func(f keyFixture) *hard_delete.Command {
				return hard_delete.NewCommand(f.id.String(), f.clientId.String(), f.adminId.String(), f.userId.String())
			},
			setupRepo: func(_ *testing.T, repo *cakmocks.MockClientAssertionKeyRepository, f keyFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(nil, nil)
			},
			wantErrCode: client_assertion_keys.ErrClientAssertionKeyNotFound,
		},
		{
			name: "wraps repository search error",
			buildCmd: func(f keyFixture) *hard_delete.Command {
				return hard_delete.NewCommand(f.id.String(), f.clientId.String(), f.adminId.String(), f.userId.String())
			},
			setupRepo: func(_ *testing.T, repo *cakmocks.MockClientAssertionKeyRepository, f keyFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(nil, errors.New("db error"))
			},
			wantErrCode: client_assertion_keys.ErrClientAssertionKeyNotFound,
		},
		{
			name: "returns forbidden for an unrelated user",
			buildCmd: func(f keyFixture) *hard_delete.Command {
				other := ids.NewUserId().String()
				return hard_delete.NewCommand(f.id.String(), f.clientId.String(), ids.NewUserId().String(), other)
			},
			setupRepo: func(t *testing.T, repo *cakmocks.MockClientAssertionKeyRepository, f keyFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(aKey(t, f), nil)
			},
			wantErrCode: client_assertion_keys.ErrInvalidOwnership,
		},
		{
			name: "returns forbidden for a sibling tenant under the same admin",
			buildCmd: func(f keyFixture) *hard_delete.Command {
				return hard_delete.NewCommand(f.id.String(), f.clientId.String(), f.adminId.String(), ids.NewUserId().String())
			},
			setupRepo: func(t *testing.T, repo *cakmocks.MockClientAssertionKeyRepository, f keyFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(aKey(t, f), nil)
			},
			wantErrCode: client_assertion_keys.ErrInvalidOwnership,
		},
		{
			name: "returns forbidden when the key belongs to a different client",
			buildCmd: func(f keyFixture) *hard_delete.Command {
				return hard_delete.NewCommand(f.id.String(), uuid.New().String(), f.adminId.String(), f.userId.String())
			},
			setupRepo: func(t *testing.T, repo *cakmocks.MockClientAssertionKeyRepository, f keyFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(aKey(t, f), nil)
			},
			wantErrCode: client_assertion_keys.ErrInvalidOwnership,
		},
		{
			name: "wraps repository hard delete error",
			buildCmd: func(f keyFixture) *hard_delete.Command {
				return hard_delete.NewCommand(f.id.String(), f.clientId.String(), f.adminId.String(), f.userId.String())
			},
			setupRepo: func(t *testing.T, repo *cakmocks.MockClientAssertionKeyRepository, f keyFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(aKey(t, f), nil)
				repo.EXPECT().HardDelete(mock.Anything, f.id).Return(errors.New("db error"))
			},
			wantErrCode: client_assertion_keys.ErrClientAssertionKeyNotSaved,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := newKeyFixture(t)
			repo := cakmocks.NewMockClientAssertionKeyRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(t, repo, f)
			}

			err := newService(repo).HardDelete(context.Background(), tc.buildCmd(f))

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
