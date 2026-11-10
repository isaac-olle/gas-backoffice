package list_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs/application/list"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubRepository struct {
	result   []error_logs.ErrorLog
	err      error
	gotLimit int
}

func (s *stubRepository) Save(_ context.Context, _ error_logs.ErrorLog) error {
	return nil
}

func (s *stubRepository) Latest(_ context.Context, limit int) ([]error_logs.ErrorLog, error) {
	s.gotLimit = limit
	return s.result, s.err
}

func TestList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		repoResult []error_logs.ErrorLog
		repoErr    error
		wantErr    bool
	}{
		{
			name:       "returns the latest errors capped at the limit",
			repoResult: []error_logs.ErrorLog{{ErrorId: "e1", Code: "X"}},
		},
		{
			name:    "returns error when repo errors",
			repoErr: fmt.Errorf("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &stubRepository{result: tt.repoResult, err: tt.repoErr}

			result, err := list.NewService(repo).List(context.Background(), list.NewQuery())

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, repo.result, result)
			assert.Equal(t, list.Limit, repo.gotLimit)
		})
	}
}
