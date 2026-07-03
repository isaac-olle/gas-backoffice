package soft_delete_test

import (
	"context"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/soft_delete"
	scopemocks "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/mocks"
)

func TestSoftDelete(t *testing.T) {
	t.Parallel()

	t.Run("calls SoftDelete with the correct scope id", func(t *testing.T) {
		t.Parallel()
		repo := scopemocks.NewMockRepository(t)
		scopeId := ids.NewScopeId()

		repo.EXPECT().SoftDelete(mock.Anything, scopeId).Return(nil)

		svc := soft_delete.NewService(repo)
		require.NoError(t, svc.SoftDelete(context.Background(), aCommand(t, withScopeId(scopeId.String()))))
	})

	errCases := []struct {
		desc      string
		cmd       func(t *testing.T) *soft_delete.Command
		setupRepo func(repo *scopemocks.MockRepository)
		wantCode  app_errors.ErrorCode
	}{
		{
			desc:     "returns ErrScopeNotFound for invalid scope id format",
			cmd:      func(t *testing.T) *soft_delete.Command { return aCommand(t, withScopeId("not-a-uuid")) },
			wantCode: scopes.ErrScopeNotFound,
		},
		{
			desc: "returns ErrScopeNotSaved when repository fails",
			cmd:  func(t *testing.T) *soft_delete.Command { return aCommand(t) },
			setupRepo: func(repo *scopemocks.MockRepository) {
				repo.EXPECT().
					SoftDelete(mock.Anything, mock.AnythingOfType("ids.ScopeId")).
					Return(fmt.Errorf("db error"))
			},
			wantCode: scopes.ErrScopeNotSaved,
		},
	}

	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			repo := scopemocks.NewMockRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			svc := soft_delete.NewService(repo)
			err := svc.SoftDelete(context.Background(), tc.cmd(t))

			assertErrCode(t, err, tc.wantCode)
		})
	}
}

// --- helpers ---

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}
