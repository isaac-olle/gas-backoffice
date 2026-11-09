package create_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs/application/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubRepository struct {
	err    error
	gotLog error_logs.ErrorLog
}

func (s *stubRepository) Save(_ context.Context, errorLog error_logs.ErrorLog) error {
	s.gotLog = errorLog
	return s.err
}

func (s *stubRepository) Latest(_ context.Context, _ int) ([]error_logs.ErrorLog, error) {
	return nil, nil
}

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("saves the error log built from the command", func(t *testing.T) {
		t.Parallel()
		repo := &stubRepository{}
		occurredAt := time.Now()
		cmd := create.NewCommand(
			"e1", "INTERNAL", "X", "boom",
			map[string]any{"k": "v"}, []string{"cause1"}, map[string]any{"trace": true},
			"req-1", &occurredAt,
		)

		err := create.NewService(repo).Create(context.Background(), cmd)
		require.NoError(t, err)
		assert.Equal(t, error_logs.ErrorLog{
			ErrorId:     "e1",
			Type:        "INTERNAL",
			Code:        "X",
			Description: "boom",
			Details:     map[string]any{"k": "v"},
			Cause:       []string{"cause1"},
			StackTrace:  map[string]any{"trace": true},
			RequestId:   "req-1",
			OccurredAt:  &occurredAt,
		}, repo.gotLog)
	})

	t.Run("returns error when repo errors", func(t *testing.T) {
		t.Parallel()
		repo := &stubRepository{err: fmt.Errorf("db error")}
		cmd := create.NewCommand("e1", "INTERNAL", "X", "boom", nil, nil, nil, "", nil)

		err := create.NewService(repo).Create(context.Background(), cmd)
		assert.Error(t, err)
	})
}
