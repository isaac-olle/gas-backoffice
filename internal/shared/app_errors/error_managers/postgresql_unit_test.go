package error_managers_test

import (
	"context"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors/error_managers"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type errorLogRepository struct{ saved error_logs.ErrorLog }

func (r *errorLogRepository) Save(_ context.Context, record error_logs.ErrorLog) error {
	r.saved = record
	return nil
}

func (*errorLogRepository) Latest(context.Context, int) ([]error_logs.ErrorLog, error) {
	return nil, nil
}

type errorLogExecutor struct{ handler *create.CommandHandler }

func (e errorLogExecutor) SingleExecution(ctx context.Context, cmd command.Command) error {
	return e.handler.Handle(ctx, cmd)
}

func (e errorLogExecutor) MultipleExecution(ctx context.Context, commands ...command.Command) error {
	for _, cmd := range commands {
		if err := e.SingleExecution(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func TestPostgresErrorManagerMetadata(t *testing.T) {
	for _, authenticated := range []bool{false, true} {
		name := "anonymous"
		if authenticated {
			name = "authenticated"
		}
		t.Run(name, func(t *testing.T) {
			repository := &errorLogRepository{}
			executor := errorLogExecutor{handler: create.NewCommandHandler(create.NewService(repository))}
			ctx := context.Background()
			authorID := ids.NewUserId().String()
			if authenticated {
				author, err := vo.NewUserAuthorFromString(authorID, authorID, vo.SubscriptionUserRole)
				require.NoError(t, err)
				ctx = vo.WithAuthor(ctx, author)
			}
			err := app_errors.BadRequestError("TEST")
			before := err.Error()
			error_managers.NewPostgresErrorManager(executor).Manage(ctx, err)
			require.Equal(t, err.ErrorId, repository.saved.ErrorId)
			require.Equal(t, before, err.Error(), "logging must not change AppError")
			if authenticated {
				require.Equal(t, authorID, repository.saved.AuthorId)
				require.Equal(t, vo.SubscriptionUserRole.String(), repository.saved.AuthorRole)
			} else {
				require.Empty(t, repository.saved.AuthorId)
				require.Empty(t, repository.saved.AuthorRole)
			}
		})
	}
}
