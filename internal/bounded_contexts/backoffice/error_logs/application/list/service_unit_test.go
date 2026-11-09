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

	t.Run("returns the latest errors capped at the limit", func(t *testing.T) {
		t.Parallel()
		repo := &stubRepository{result: []error_logs.ErrorLog{{ErrorId: "e1", Code: "X"}}}

		result, err := list.NewService(repo).List(context.Background(), list.NewQuery())
		require.NoError(t, err)
		assert.Equal(t, repo.result, result)
		assert.Equal(t, list.Limit, repo.gotLimit)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := &stubRepository{err: fmt.Errorf("db error")}

		_, err := list.NewService(repo).List(context.Background(), list.NewQuery())
		assert.Error(t, err)
	})
}
