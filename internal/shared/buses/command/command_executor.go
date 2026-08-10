package command

import "context"

// CommandExecutor is the handler-facing façade for executing commands. Handlers
// depend on this, never on the raw Bus: SingleExecution runs one command,
// MultipleExecution runs several as one atomic unit. The transaction is hidden
// from the caller.
type CommandExecutor interface {
	SingleExecution(ctx context.Context, cmd Command) error
	MultipleExecution(ctx context.Context, commands ...Command) error
}
