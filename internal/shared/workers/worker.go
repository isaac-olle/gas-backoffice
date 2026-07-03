package workers

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"sync"

	"github.com/google/uuid"
)

type Worker struct {
	ID       string
	ch       chan func(context.Context, event.Event)
	isActive bool
	mu       sync.Mutex
}

func NewWorker() *Worker {
	return &Worker{
		ID:       uuid.New().String(),
		ch:       make(chan func(context.Context, event.Event)),
		isActive: false,
	}
}

func (w *Worker) SetActive() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isActive = true
}

func (w *Worker) SetUnActive() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isActive = false
}

func (w *Worker) Active() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.isActive
}

func (w *Worker) Work(ctx context.Context, task func(context.Context)) {
	task(ctx)
}

func (w *Worker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.isActive = false

	select {
	case <-w.ch:
	default:
		close(w.ch)
	}
}
