package error_managers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	errorLogCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

// PostgresErrorManager persists every managed AppError as a row in
// public.errors, so failures can be inspected after the fact without
// depending on log retention. It dispatches the error_logs create command
// rather than writing to Postgres itself.
type PostgresErrorManager struct {
	commandExecutor command.CommandExecutor
}

func NewPostgresErrorManager(commandExecutor command.CommandExecutor) *PostgresErrorManager {
	return &PostgresErrorManager{commandExecutor: commandExecutor}
}

func (m *PostgresErrorManager) Manage(ctx context.Context, err *app_errors.AppError) {
	var stackTrace any
	if err.StackTrace != nil {
		stackTrace = err.StackTrace.ToJSON()
	}

	causes := make([]string, 0, len(err.Cause))
	for _, c := range err.Cause {
		causes = append(causes, c.Error())
	}

	occurredAt := err.OccurredAt
	var authorId, authorRole string
	if author, ok := vo.AuthorFromContext(ctx); ok && author != nil {
		authorId, authorRole = author.Id.String(), author.Role.String()
	}
	cmd := errorLogCreate.NewCommand(
		err.ErrorId, err.Type.String(), string(err.Code), err.Description,
		err.Details, causes, stackTrace, err.RequestId, &occurredAt,
	).WithAuthor(authorId, authorRole)

	if execErr := m.commandExecutor.SingleExecution(context.Background(), cmd); execErr != nil {
		logPersistenceFailure(os.Stderr, err.ErrorId, execErr)
	}
}

func logPersistenceFailure(writer io.Writer, errorId string, err error) {
	var details any = err
	var appError *app_errors.AppError
	if errors.As(err, &appError) {
		details = appError.ToPrimitives()
	}
	fmt.Fprintf(writer, "PostgresErrorManager: failed to persist error %s: %+v\n", errorId, details)
}
