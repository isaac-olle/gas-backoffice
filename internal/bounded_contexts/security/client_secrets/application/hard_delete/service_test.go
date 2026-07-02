package hard_delete_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets"
	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets/application/hard_delete"
	csmocks "github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets/mocks"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
)

func assertErrCode(t *testing.T, err error, code app_errors.ErrorCode) {
	t.Helper()
	var appErr *app_errors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code, appErr.Code)
}

func TestHardDeleteClientSecretService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		id            string
		hardDeleteErr error
		wantErr       bool
		wantErrCode   app_errors.ErrorCode
	}{
		{
			name: "deletes client secret by id",
			id:   uuid.New().String(),
		},
		{
			name:        "rejects invalid id",
			id:          "not-a-valid-uuid",
			wantErr:     true,
			wantErrCode: client_secrets.ErrClientSecretNotFound,
		},
		{
			name:          "wraps repository hard delete error",
			id:            uuid.New().String(),
			hardDeleteErr: errors.New("db error"),
			wantErr:       true,
			wantErrCode:   client_secrets.ErrClientSecretNotSaved,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := csmocks.NewMockRepository(t)
			if tc.id != "not-a-valid-uuid" {
				repo.EXPECT().HardDelete(mock.Anything, mock.Anything).Return(tc.hardDeleteErr)
			}

			svc := hard_delete.NewService(repo)

			cmd := hard_delete.NewCommand(tc.id)
			err := svc.HardDelete(context.Background(), cmd)

			if tc.wantErr {
				assertErrCode(t, err, tc.wantErrCode)
				return
			}
			require.NoError(t, err)
		})
	}
}
