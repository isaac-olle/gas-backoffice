package workers

import (
	"context"
	"sync"
)

type PoolStats struct {
	PendingTasks     int
	ActiveWorkers    int
	AvailableWorkers int
}

type Pool struct {
	availableWorkers *Workers
	activeWorkers    *Workers
	mu               sync.Mutex
	pendingTasks     PendingTasks
	totalWorkers     int
}

// NewPool inicialitza N workers
func NewPool(numWorkers int) *Pool {
	p := &Pool{
		availableWorkers: NewWorkers(numWorkers),
		activeWorkers:    NewWorkers(0),
		pendingTasks:     NewPendingTasks(),
		totalWorkers:     numWorkers,
	}
	return p
}

// Process envia una feina a un worker disponible o l'enfila si no n'hi ha
func (p *Pool) Process(ctx context.Context, task func(ctx context.Context)) {
	p.mu.Lock()
	worker := p.availableWorkers.PullFirst()
	if worker == nil {
		p.pendingTasks.Push(task)
		p.mu.Unlock()
		return
	}

	worker.SetActive()
	p.activeWorkers.Push(worker)
	p.mu.Unlock()

	go func(w *Worker) {
		defer func() {
			w.SetUnActive()

			p.mu.Lock()
			p.activeWorkers.PullByID(w.ID)

			nextTask := p.pendingTasks.PullFirst()
			if nextTask != nil {
				p.mu.Unlock()
				p.Process(ctx, nextTask)
				return
			}

			p.availableWorkers.Push(w)
			p.mu.Unlock()
		}()

		w.Work(ctx, task)
	}(worker)
}

func (p *Pool) Stats() PoolStats {
	p.mu.Lock()
	defer p.mu.Unlock()
	return PoolStats{
		PendingTasks:     p.pendingTasks.Length(),
		ActiveWorkers:    p.activeWorkers.Length(),
		AvailableWorkers: p.availableWorkers.Length(),
	}
}

func (p *Pool) Scale(delta int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if delta == 0 {
		return
	}

	if delta > 0 {
		for range delta {
			p.addWorker()
		}
		return
	}

	if delta < 0 {
		for range -delta {
			p.removeWorker()
		}
		return
	}

}

func (p *Pool) addWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	w := NewWorker()
	p.availableWorkers.Push(w)
	p.totalWorkers++
}

func (p *Pool) removeWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	w := p.availableWorkers.PullFirst()
	if w != nil {
		w.Stop()
		p.totalWorkers--
	}
}
