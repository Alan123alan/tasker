package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type Status int

const (
	Started Status = iota
	Blocked
	Completed
)

func (enum Status) String() string {
	switch enum {
	case Started:
		return "Started"
	case Blocked:
		return "Blocked"
	case Completed:
		return "Completed"
	default:
		return fmt.Sprintf("%d", int(enum))
	}
}

type Task struct {
	Description string
	Status      Status
	StartedAt   time.Time
	CompletedAt time.Time
}

type TaskModel struct {
	Id          string
	Description string
	Status      Status
	StartedAt   time.Time
	CompletedAt time.Time
}

func createTaskTable(db *sql.DB) {
	stmnt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, description TEXT, status INTEGER, started_at DATETIME, completed_at DATETIME);`)
	if err != nil {
		log.Fatal(err)
	}
	stmnt.Exec()
}

func printToDoTable(tasks []TaskModel) {
	headerFormat := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFormat := color.New(color.FgYellow).SprintfFunc()

	table := table.New("ID", "Description", "Status", "Started", "Finished")
	table.WithHeaderFormatter(headerFormat).WithFirstColumnFormatter(columnFormat)

	for _, task := range tasks {
		table.AddRow(task.Id, task.Description, Status(task.Status), task.StartedAt, task.CompletedAt)
	}

	table.Print()
}

func saveTask(db *sql.DB, task Task) {
	stmnt, err := db.Prepare("INSERT INTO tasks(description, status, started_at, completed_at) values (?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmnt.Exec(task.Description, task.Status, task.StartedAt, task.CompletedAt)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Rows afected: %v", rowsAffected)
}

// func getTasks(db *sql.DB) (tasks []TaskModel) {
// 	rows, err := db.Query("SELECT id, description, status, started_at, completed_at FROM tasks")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	var task TaskModel
// 	for rows.Next() {
// 		err := rows.Scan(&task.Id, &task.Description, &task.Status, &task.StartedAt, &task.CompletedAt)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		tasks = append(tasks, task)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return tasks
// }

func getTask(db *sql.DB, id string) (tasks []TaskModel) {
	rows, err := db.Query("SELECT id, description, status, started_at, completed_at FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var task TaskModel
	for rows.Next() {
		err := rows.Scan(&task.Id, &task.Description, &task.Status, &task.StartedAt, &task.CompletedAt)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return tasks
}

func getTasks(DB *sql.DB) tea.Cmd {
	return func() tea.Msg {
		rows, err := DB.Query("SELECT id, description, status, started_at, completed_at FROM tasks")
		if err != nil {
			return errMsg{err}
		}
		defer rows.Close()
		var task TaskModel
		tasks := []TaskModel{}
		for rows.Next() {
			err := rows.Scan(&task.Id, &task.Description, &task.Status, &task.StartedAt, &task.CompletedAt)
			if err != nil {
				return errMsg{err}
			}
			tasks = append(tasks, task)
		}
		err = rows.Err()
		if err != nil {
			return errMsg{err}
		}
		return taskMsg(tasks)

	}
}

type taskMsg []TaskModel

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

// func updateTask(db *sql.DB, id string, description string, status Status) {
// 	task := getTask(db, id)[0]
// 	if description != "" {
// 		task.Description = description
// 	}
// 	if status != Started {
// 		task.Status = status
// 		if status == Completed {
// 			task.CompletedAt = time.Now()
// 		}
// 	}
// 	stmnt, err := db.Prepare(`UPDATE tasks SET description=?, status=?, started_at=?, completed_at=? WHERE id = ?`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	res, err := stmnt.Exec(task.Description, task.Status, task.StartedAt, task.CompletedAt, task.Id)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Rows afected: %v", rowsAffected)
// }
