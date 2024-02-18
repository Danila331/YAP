package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
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
	dbmutex      sync.Mutex
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
	ticker := time.NewTicker(time.Duration(w.id*1) * time.Minute)
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

	rows, err := w.db.Query("SELECT id, expression FROM tasks WHERE processed = 0 AND status = ? LIMIT 5", "pending")
	if err != nil {
		log.Printf("Worker %d: Failed to query tasks: %v", w.id, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		fmt.Println(rows)
		if err := rows.Scan(&task.Id, &task.Expression); err != nil {
			log.Printf("Worker %d: Failed to scan task: %v", w.id, err)
			continue
		}
		w.taskQueue <- &task
	}
}

// calculateExpression вычисляет значение выражения
func (w *Worker) calculateExpression(expression string) string {
	var timer TimeOperations
	timer, err := timer.Read()
	if err != nil {
		return ""
	}
	countPulse := strings.Count(expression, "+") * timer.TimePulse
	countMinus := strings.Count(expression, "-") * timer.TimeMinus
	countProz := strings.Count(expression, "*") * timer.TimeProz
	countDel := strings.Count(expression, "/") * timer.TimeDel
	resultTime := time.Duration(countPulse + countMinus + countDel + countProz)
	evaluator, _ := govaluate.NewEvaluableExpression(expression)

	result, _ := evaluator.Evaluate(nil)
	resultFloat, _ := result.(float64)
	// Имитация долгих вычислений
	resultStr := strconv.FormatFloat(resultFloat, 'f', -1, 64)
	time.Sleep(resultTime * time.Second)
	return resultStr
}

// saveResult сохраняет результат в базе данных
func (w *Worker) saveResult(taskID int, result string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	_, err := w.db.Exec("UPDATE tasks SET result = ?, processed = 1, status = 'ok' WHERE id = ?", result, taskID)
	if err != nil {
		log.Printf("Worker %d: Failed to save result: %v", w.id, err)
	}
}
