package hard_delete_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys/application/hard_delete"
	cakmocks "github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys/domain_services/mocks"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestHardDeleteClientAssertionKeyService_HappyPath(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()

	var deletedId ids.ClientAssertionKeyId
	repo := cakmocks.NewMockClientAssertionKeyRepository(t)
	repo.EXPECT().HardDelete(mock.Anything, mock.Anything).Run(func(_ context.Context, id ids.ClientAssertionKeyId) {
		deletedId = id
	}).Return(nil)

	svc := hard_delete.NewService(repo)

	require.NoError(t, svc.HardDelete(context.Background(), hard_delete.NewCommand(id)))

	assert.Equal(t, id, deletedId.String(), "repo must receive the correct ID for deletion")
}

func TestHardDeleteClientAssertionKeyService_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		id            string
		hardDeleteErr error
		wantErrCode   app_errors.ErrorCode
	}{
		{
			name:        "rejects invalid id",
			id:          "not-a-valid-uuid",
			wantErrCode: client_assertion_keys.ErrClientAssertionKeyNotFound,
		},
		{
			name:          "wraps repository hard delete error",
			id:            uuid.New().String(),
			hardDeleteErr: errors.New("db error"),
			wantErrCode:   client_assertion_keys.ErrClientAssertionKeyNotSaved,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := cakmocks.NewMockClientAssertionKeyRepository(t)
			if tc.hardDeleteErr != nil {
				repo.EXPECT().HardDelete(mock.Anything, mock.Anything).Return(tc.hardDeleteErr)
			}

			svc := hard_delete.NewService(repo)

			cmd := hard_delete.NewCommand(tc.id)
			err := svc.HardDelete(context.Background(), cmd)

			assertErrCode(t, err, tc.wantErrCode)
		})
	}
}
