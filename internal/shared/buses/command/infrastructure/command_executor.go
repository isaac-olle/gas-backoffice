package infrastructure

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

// CommandExecutor is the single command-execution façade for handlers: they
// always depend on it, so going from one command to several is a method change,
// not a re-wiring, and they never touch the raw bus.
//
//   - SingleExecution delegates straight to the command bus, whose transactional
//     decorator opens the transaction for that one command.
//   - MultipleExecution opens one unit of work around the whole batch (the bus's
//     decorator then reuses it), so all commands commit or none does.
type CommandExecutor struct {
	bus command.Bus
	uow persistance.UnitOfWork
}

var _ command.CommandExecutor = (*CommandExecutor)(nil)

func NewCommandExecutor(bus command.Bus, uow persistance.UnitOfWork) *CommandExecutor {
	return &CommandExecutor{bus: bus, uow: uow}
}

func (e *CommandExecutor) SingleExecution(ctx context.Context, cmd command.Command) error {
	return e.bus.Dispatch(ctx, cmd)
}

func (e *CommandExecutor) MultipleExecution(ctx context.Context, commands ...command.Command) error {
	return e.uow.Do(ctx, func(ctx context.Context) error {
		for _, cmd := range commands {
			if err := e.bus.Dispatch(ctx, cmd); err != nil {
				return err
			}
		}
		return nil
	})
}
