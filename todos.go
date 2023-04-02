package main

import (
	"database/sql"
	"fmt"
	"time"

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

type Todo struct {
	Id          string
	Description string
	Status      Status
	StartedAt   time.Time
	CompletedAt time.Time
}

func printToDoTable(todos []Todo) {
	headerFormat := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFormat := color.New(color.FgYellow).SprintfFunc()

	table := table.New("ID", "Description", "Status", "Started", "Finished")
	table.WithHeaderFormatter(headerFormat).WithFirstColumnFormatter(columnFormat)

	for _, todo := range todos {
		table.AddRow(todo.Id, todo.Description, Status(todo.Status), todo.StartedAt, todo.CompletedAt)
	}

	table.Print()
}

func saveToDo(db *sql.DB, todo Todo) {
	stmnt, err := db.Prepare("INSERT INTO todos(id, description, status, started_at, completed_at) values (?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	res, err := stmnt.Exec(todo.Id, todo.Description, todo.Status, todo.StartedAt, todo.CompletedAt)
	if err != nil {
		panic(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Rows afected: %v", rowsAffected)
}

func getTodos(db *sql.DB) (todos []Todo) {
	rows, err := db.Query("SELECT id, description, status, started_at, completed_at FROM todos")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var todo Todo
	for rows.Next() {
		err := rows.Scan(&todo.Id, &todo.Description, &todo.Status, &todo.StartedAt, &todo.CompletedAt)
		if err != nil {
			panic(err)
		}
		todos = append(todos, todo)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return todos
}

func getTodo(db *sql.DB, id string) (todos []Todo) {
	rows, err := db.Query("SELECT id, description, status, started_at, completed_at FROM todos WHERE id = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var todo Todo
	for rows.Next() {
		err := rows.Scan(&todo.Id, &todo.Description, &todo.Status, &todo.StartedAt, &todo.CompletedAt)
		if err != nil {
			panic(err)
		}
		todos = append(todos, todo)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return todos
}

func updateTodo(db *sql.DB, id string, description string, status Status) {
	todo := getTodo(db, id)[0]
	if description != "" {
		todo.Description = description
	}
	if status != Started {
		todo.Status = status
		if status == Completed {
			todo.CompletedAt = time.Now()
		}
	}
	stmnt, err := db.Prepare(`UPDATE todos SET description=?, status=?, started_at=?, completed_at=? WHERE id = ?`)
	if err != nil {
		panic(err)
	}
	res, err := stmnt.Exec(todo.Description, todo.Status, todo.StartedAt, todo.CompletedAt, todo.Id)
	if err != nil {
		panic(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Rows afected: %v", rowsAffected)
}
