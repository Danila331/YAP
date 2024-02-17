package models

import (
	"fmt"
	"github/Danila331/YAP/internal/pkg"
	"github/Danila331/YAP/internal/store"
)

// Структура арифметического выражения
type Task struct {
	Id         int    `db:"id" json:"id"`
	IdServer   int    `db:"idservera" json:"idservera"` // Add field to database
	Expression string `db:"expression" json:"expression"`
	Result     string `db:"result" json:"result"`
	Status     string `db:"status" json:"status"`
	StartDate  string `db:"startdate" json:"startdate"`
	Processed  bool   `db:"processed" json:"processed"` // Вообще хз ля чего я это добавил
}

// Интерфейс для реализации работы с базой данных
type TaskInterface interface {
	Create() error                       // Функция для занесения новой операции в BD
	ReadAll() ([]Task, error)            // Функция для получения операций списком
	ReadById(taskId int) (Task, error)   // Функция для получения операции по ID
	Update(taskId, taskResult int) error // Функция для обновления статусы и результата операции
}

func (t *Task) Create() error {

	conn, err := store.ConnectToDatabase()
	defer conn.Close()

	if err != nil {
		return err
	}

	queryRow := fmt.Sprintf("INSERT INTO tasks(expression, result, status, startdate,idservera) VALUES (?, ?, ?, ?,?)")

	_, err = conn.Exec(queryRow, t.Expression, "", "pending", t.StartDate, pkg.Random(1, 5))

	if err != nil {
		return err
	}

	return nil
}

func (t *Task) ReadAll() ([]Task, error) {
	var tasks []Task

	conn, err := store.ConnectToDatabase()
	defer conn.Close()

	if err != nil {
		return []Task{}, err
	}

	queryRow := fmt.Sprintf("SELECT * FROM tasks")

	rows, err := conn.Query(queryRow)
	for rows.Next() {
		var task Task

		err = rows.Scan(
			&task.Id,
			&task.Expression,
			&task.Result,
			&task.Status,
			&task.StartDate,
			&task.Processed,
			&task.IdServer,
		)

		if err != nil {
			return []Task{}, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *Task) ReadById() (Task, error) {
	var task Task

	conn, err := store.ConnectToDatabase()
	defer conn.Close()

	if err != nil {
		return Task{}, err
	}

	queryRow := fmt.Sprintf("SELECT id, expression, result, status, startdate FROM tasks WHERE id = ?")
	row := conn.QueryRow(queryRow, t.Id)

	err = row.Scan(&task.Id, &task.Expression, &task.Result, &task.Status, &task.StartDate)

	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (t *Task) Update() error {
	conn, err := store.ConnectToDatabase()
	defer conn.Close()

	if err != nil {
		return err
	}

	queryRow := fmt.Sprintf("UPDATE tasks SET idservera = ?, status = ?, result = ?, processed = 1 WHERE id = ?")
	_, err = conn.Exec(queryRow, t.IdServer, t.Status, t.Result, t.Id)

	if err != nil {
		return err
	}

	return nil
}
