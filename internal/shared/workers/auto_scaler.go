package workers

import "time"

type AutoScaler struct {
	pool                *Pool
	minWorkers          int
	maxWorkers          int
	maxPendingTasks     int
	maxAvailableWorkers int
	tickInterval        time.Duration
}

func NewAutoScaler(pool *Pool, minWorkers, maxWorkers, maxPendingTasks, maxAvailableWorkers int, tickInterval time.Duration) *AutoScaler {
	return &AutoScaler{
		pool:                pool,
		minWorkers:          minWorkers,
		maxWorkers:          maxWorkers,
		maxPendingTasks:     maxPendingTasks,
		maxAvailableWorkers: maxAvailableWorkers,
		tickInterval:        tickInterval,
	}
}

func (a *AutoScaler) Start() {
	ticker := time.NewTicker(a.tickInterval)
	defer ticker.Stop()

	for range ticker.C {
		delta := a.decide(a.pool.Stats())
		if delta != 0 {
			a.pool.Scale(delta)
		}
	}
}

func (a *AutoScaler) decide(stats PoolStats) int {
	total := stats.ActiveWorkers + stats.AvailableWorkers

	if stats.PendingTasks > a.maxPendingTasks && total < a.maxWorkers {
		return 1
	}
	if stats.AvailableWorkers > a.maxAvailableWorkers && total > a.minWorkers {
		return -1
	}
	return 0
}
