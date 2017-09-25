package main

import (
	"database/sql"
	"time"
)

type Task struct {
	Id         string `json:"id"`
	Rule       string `json:"rule"`
	TimeZone   string `json:"timeZone"`
	Epsilon    int    `json:"epsilon"`
	MaxRetries int    `json:"maxRetries"`
}

type DbTask struct {
	*Task
	Completed   bool       `json:"completed"`
	NextFireAt  *time.Time `json:"nextFireAt,omitempty"`
	NextRetryAt *time.Time `json:"nextFireAt,omitempty"`
}

func getRunnableTasks(db sql.DB, firingAtTo time.Time, limit int) ([]DbTask, error) {
	stmt, err := db.Prepare(`
		select id, rule, timeZone, epsilon, maxRetries, completed from tasks
		where completed == 0 and (nextFireAt <= ? or nextFireAt is null)
		order by nextFireAt asc
		limit ?
	`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(firingAtTo.UTC().Unix(), limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return initTasksFrom(rows)
}

func initTasksFrom(rows *sql.Rows) ([]DbTask, error) {
	var tasks = []DbTask{}

	for rows.Next() {
		var id string
		var rule string
		var timeZone string
		var epsilon int
		var maxRetries int
		var completed bool

		err := rows.Scan(&id, &rule, &timeZone, &epsilon, &maxRetries, &completed)

		if err != nil {
			return nil, err
		}

		var task = &DbTask{
			Task:        &Task{Id: id, Rule: rule, TimeZone: timeZone, Epsilon: epsilon},
			Completed:   completed,
			NextFireAt:  nil,
			NextRetryAt: nil}

		tasks = append(tasks, *task)
	}

	err := rows.Err()

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func getTasks(db sql.DB) ([]DbTask, error) {
	rows, err := db.Query("select id, rule, timeZone, epsilon, maxRetries, completed from tasks")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return initTasksFrom(rows)
}

func getTask(db sql.DB, taskId string) (*DbTask, error) {
	stmt, err := db.Prepare(`
		select id, rule, timeZone, epsilon, maxRetries, completed from tasks
	    where id = ? limit 1
 	`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var id string
	var rule string
	var timeZone string
	var epsilon int
	var maxRetries int
	var completed bool

	row := stmt.QueryRow(taskId)

	err = row.Scan(&id, &rule, &timeZone, &epsilon, &maxRetries, &completed)

	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var task = &DbTask{Task: &Task{Id: id, Rule: rule, TimeZone: timeZone, Epsilon: epsilon}, Completed: completed, NextFireAt: nil, NextRetryAt: nil}

	return task, nil
}

func updateTaskNextFireAt(db sql.DB, taskId string, nextFireAt time.Time, completed bool) error {
	stmt, err := db.Prepare(`
		update tasks
		set nextFireAt = ?, completed = ?
		where id = ?
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(nextFireAt.UTC().Unix(), completed, taskId)

	if err != nil {
		return err
	}

	return nil
}

func addTask(db sql.DB, task Task) (*DbTask, error) {
	dbTask := &DbTask{Task: &task, Completed: false, NextFireAt: nil, NextRetryAt: nil}

	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(`insert into tasks(id, rule, timeZone, epsilon, maxRetries, completed) values(?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(&dbTask.Id, &dbTask.Rule, &dbTask.TimeZone, &dbTask.Epsilon, &dbTask.MaxRetries, &dbTask.Completed)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return dbTask, nil
}
