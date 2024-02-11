package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
)

// Worker представляет воркера
type Worker struct {
	id           int
	db           *sql.DB
	taskQueue    chan *Task
	workerChange chan *Worker
	stopCh       chan struct{}
	mutex        sync.Mutex
	dbMutex      sync.Mutex
}

// NewWorker создает нового воркера
func NewWorker(id int, db *sql.DB, workerChange chan *Worker) *Worker {
	return &Worker{
		id:           id,
		db:           db,
		taskQueue:    make(chan *Task),
		workerChange: workerChange,
		stopCh:       make(chan struct{}),
	}
}

// Start запускает воркера
func (w *Worker) Start() {
	for i := 0; i < 5; i++ {
		go w.processTask()
	}

	go w.monitorTasks()
}

// Stop останавливает воркера
func (w *Worker) Stop() {
	close(w.stopCh)
}

// processTask обрабатывает задачи вычисления выражений
func (w *Worker) processTask() {
	for {
		select {
		case <-w.stopCh:
			return
		case task := <-w.taskQueue:
			if !task.Processed {
				result := w.calculateExpression(task.Expression)
				w.saveResult(task.Id, result)
			}
		}
	}
}

// monitorTasks периодически проверяет базу данных на наличие новых задач
func (w *Worker) monitorTasks() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.fetchTasks()
		}
	}
}

// fetchTasks извлекает задачи из базы данных и отправляет их в канал задач воркера
func (w *Worker) fetchTasks() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.dbMutex.Lock()
	defer w.dbMutex.Unlock()
	rows, err := w.db.Query("SELECT id, expression FROM tasks WHERE processed = 0 LIMIT 5")
	if err != nil {
		log.Printf("Worker %d: Failed to query tasks: %v", w.id, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Expression); err != nil {
			log.Printf("Worker %d: Failed to scan task: %v", w.id, err)
			continue
		}
		task.IdServer = w.id
		task.Status = "processed"
		time.Sleep(60 * time.Second)
		if err := task.Update(); err != nil {
			log.Printf("Worker %d: Failed to Update: %v", w.id, err)
			continue
		}
		w.taskQueue <- &task
	}
}

// calculateExpression вычисляет значение выражения
func (w *Worker) calculateExpression(expression string) string {
	resultChan := make(chan string)

	go func() {
		evaluator, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			resultChan <- fmt.Sprintf("Error: %v", err)
			return
		}

		result, err := evaluator.Evaluate(nil)
		if err != nil {
			resultChan <- fmt.Sprintf("Error: %v", err)
			return
		}

		resultFloat, ok := result.(float64)
		if !ok {
			resultChan <- "Error: result is not a float64"
			return
		}

		time.Sleep(5 * time.Second) // Имитация долгих вычислений

		resultStr := strconv.FormatFloat(resultFloat, 'f', -1, 64)
		resultChan <- resultStr
	}()

	select {
	case result := <-resultChan:
		return result
	case <-time.After(10 * time.Second): // Ограничение времени выполнения
		return "Error: calculation timed out"
	}
}

// saveResult сохраняет результат в базе данных
func (w *Worker) saveResult(taskID int, result string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.dbMutex.Lock()
	defer w.dbMutex.Unlock()
	_, err := w.db.Exec("UPDATE tasks SET status = 'completed' result = ?, processed = 1 WHERE id = ?", result, taskID)
	if err != nil {
		log.Printf("Worker %d: Failed to save result: %v", w.id, err)
	}
}
