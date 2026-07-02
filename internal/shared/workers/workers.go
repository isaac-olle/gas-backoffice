package workers

import "sync"

type Workers struct {
	workers map[string]*Worker
	mu      sync.Mutex
}

// NewWorkers inicialitza N workers amb UUIDs
func NewWorkers(num int) Workers {
	w := Workers{
		workers: make(map[string]*Worker, num),
	}
	for i := 0; i < num; i++ {
		worker := NewWorker()
		w.workers[worker.ID] = worker
	}
	return w
}

// PullFirst retorna qualsevol worker i el treu del map
func (w *Workers) PullFirst() *Worker {
	w.mu.Lock()
	defer w.mu.Unlock()

	for id, worker := range w.workers {
		delete(w.workers, id)
		return worker
	}
	return nil
}

// Push afegeix o reemplaça un worker pel seu ID
func (w *Workers) Push(worker *Worker) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.workers[worker.ID] = worker
}

// PullByID treu un worker pel seu UUID
func (w *Workers) PullByID(ID string) *Worker {
	w.mu.Lock()
	defer w.mu.Unlock()

	worker, ok := w.workers[ID]
	if !ok {
		return nil
	}
	delete(w.workers, ID)
	return worker
}

func (w *Workers) Length() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.workers)
}
