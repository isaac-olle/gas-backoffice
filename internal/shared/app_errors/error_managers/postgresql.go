package error_managers

import (
	"context"
	"fmt"
	"os"

	errorLogCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
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

func (m *PostgresErrorManager) Manage(err *app_errors.AppError) {
	var stackTrace any
	if err.StackTrace != nil {
		stackTrace = err.StackTrace.ToJSON()
	}

	causes := make([]string, 0, len(err.Cause))
	for _, c := range err.Cause {
		causes = append(causes, c.Error())
	}

	occurredAt := err.OccurredAt
	cmd := errorLogCreate.NewCommand(
		err.ErrorId, err.Type.String(), string(err.Code), err.Description,
		err.Details, causes, stackTrace, err.RequestId, &occurredAt,
	)

	if execErr := m.commandExecutor.SingleExecution(context.Background(), cmd); execErr != nil {
		fmt.Fprintf(os.Stderr, "PostgresErrorManager: failed to persist error %s: %v\n", err.ErrorId, execErr)
	}
}
