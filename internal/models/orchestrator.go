package models

import (
	"database/sql"
	"fmt"
	"sync"
)

// Orchestrator представляет оркестратора задач
type Orchestrator struct {
	db            *sql.DB
	workerChange  chan *Worker
	workers       []*Worker
	workerMutex   sync.Mutex
	workerCounter int
}

// NewOrchestrator создает нового оркестратора
func NewOrchestrator(db *sql.DB) *Orchestrator {
	return &Orchestrator{
		db:           db,
		workerChange: make(chan *Worker),
	}
}

// Start запускает оркестратора
func (o *Orchestrator) Start() {
	for i := 0; i < 1; i++ {
		o.workerCounter++
		worker := NewWorker(o.workerCounter, o.db, o.workerChange)
		o.workers = append(o.workers, worker)
		worker.Start()
		fmt.Println("Worker start work")
	}

	go o.monitorWorkers()
}

// monitorWorkers мониторит состояние воркеров и обрабатывает события изменения состояния
func (o *Orchestrator) monitorWorkers() {
	for {
		select {
		case worker := <-o.workerChange:
			o.workerMutex.Lock()
			defer o.workerMutex.Unlock()

			for i, w := range o.workers {
				if w == worker {
					o.workers[i] = NewWorker(o.workerCounter, o.db, o.workerChange)
					o.workers[i].Start()
					o.workerCounter++
					break
				}
			}
		}
	}
}
