package hard_delete_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/application/hard_delete"
	csmocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/mocks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func ownerCmd(t *testing.T, f secretFixture) *hard_delete.Command {
	return hard_delete.NewCommand(f.id.String(), f.clientId.String(), authorFor(t, f.subscriptionId))
}

func newService(repo *csmocks.MockRepository) *hard_delete.Service {
	return hard_delete.NewService(repo, sharedDomainServices.NewSubscriptionAuthorizer())
}

func TestHardDeleteClientSecretService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		buildCmd    func(t *testing.T, f secretFixture) *hard_delete.Command
		setupRepo   func(t *testing.T, repo *csmocks.MockRepository, f secretFixture)
		wantErr     bool
		wantErrCode app_errors.ErrorCode
	}{
		{
			name:     "owning subscription deletes its secret",
			buildCmd: ownerCmd,
			setupRepo: func(t *testing.T, repo *csmocks.MockRepository, f secretFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(aSecret(t, f), nil)
				repo.EXPECT().HardDelete(mock.Anything, f.id).Return(nil)
			},
		},
		{
			name: "rejects invalid id",
			buildCmd: func(t *testing.T, f secretFixture) *hard_delete.Command {
				return hard_delete.NewCommand("not-a-valid-uuid", f.clientId.String(), authorFor(t, f.subscriptionId))
			},
			wantErr:     true,
			wantErrCode: client_secrets.ErrClientSecretNotFound,
		},
		{
			name: "rejects invalid client id",
			buildCmd: func(t *testing.T, f secretFixture) *hard_delete.Command {
				return hard_delete.NewCommand(f.id.String(), "not-a-uuid", authorFor(t, f.subscriptionId))
			},
			wantErr:     true,
			wantErrCode: client_secrets.ErrInvalidOwnership,
		},
		{
			name:     "propagates not found from repository",
			buildCmd: ownerCmd,
			setupRepo: func(_ *testing.T, repo *csmocks.MockRepository, f secretFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(nil,
					app_errors.NotFoundError(client_secrets.ErrClientSecretNotFound).WithDescription("client secret not found"))
			},
			wantErr:     true,
			wantErrCode: client_secrets.ErrClientSecretNotFound,
		},
		{
			name:     "returns not found when repository returns nil",
			buildCmd: ownerCmd,
			setupRepo: func(_ *testing.T, repo *csmocks.MockRepository, f secretFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(nil, nil)
			},
			wantErr:     true,
			wantErrCode: client_secrets.ErrClientSecretNotFound,
		},
		{
			name:     "wraps raw repository search error",
			buildCmd: ownerCmd,
			setupRepo: func(_ *testing.T, repo *csmocks.MockRepository, f secretFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(nil, errors.New("db error"))
			},
			wantErr:     true,
			wantErrCode: client_secrets.ErrClientSecretNotFound,
		},
		{
			name: "returns forbidden when author lacks subscription access",
			buildCmd: func(t *testing.T, f secretFixture) *hard_delete.Command {
				return hard_delete.NewCommand(f.id.String(), f.clientId.String(), authorFor(t, ids.NewSubscriptionId()))
			},
			setupRepo: func(t *testing.T, repo *csmocks.MockRepository, f secretFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(aSecret(t, f), nil)
			},
			wantErr:     true,
			wantErrCode: client_secrets.ErrInvalidOwnership,
		},
		{
			name: "returns forbidden when the secret belongs to a different client",
			buildCmd: func(t *testing.T, f secretFixture) *hard_delete.Command {
				return hard_delete.NewCommand(f.id.String(), uuid.New().String(), authorFor(t, f.subscriptionId))
			},
			setupRepo: func(t *testing.T, repo *csmocks.MockRepository, f secretFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(aSecret(t, f), nil)
			},
			wantErr:     true,
			wantErrCode: client_secrets.ErrInvalidOwnership,
		},
		{
			name:     "wraps repository hard delete error",
			buildCmd: ownerCmd,
			setupRepo: func(t *testing.T, repo *csmocks.MockRepository, f secretFixture) {
				repo.EXPECT().SearchById(mock.Anything, f.id).Return(aSecret(t, f), nil)
				repo.EXPECT().HardDelete(mock.Anything, f.id).Return(errors.New("db error"))
			},
			wantErr:     true,
			wantErrCode: client_secrets.ErrClientSecretNotSaved,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := newSecretFixture(t)
			repo := csmocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(t, repo, f)
			}

			err := newService(repo).HardDelete(context.Background(), tc.buildCmd(t, f))

			if tc.wantErr {
				assertErrCode(t, err, tc.wantErrCode)
				return
			}
			require.NoError(t, err)
		})
	}
}
