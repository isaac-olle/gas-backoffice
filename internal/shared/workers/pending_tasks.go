package workers

import (
	"context"
)

// PendingTasks és un slice de funcions
type PendingTasks []func(ctx context.Context)

// NewPendingTasks crea una nova cua de pending tasks buida
func NewPendingTasks() PendingTasks {
	return PendingTasks{}
}

// PullFirst treu la primera feina o retorna nil si està buida
func (pt *PendingTasks) PullFirst() func(ctx context.Context) {
	if len(*pt) == 0 {
		return nil
	}

	task := (*pt)[0]
	*pt = (*pt)[1:]
	return task
}

// Push afegeix una feina al final de la cua
func (pt *PendingTasks) Push(task func(ctx context.Context)) {
	*pt = append(*pt, task)
}

func (pt *PendingTasks) Length() int {
	return len(*pt)
}
