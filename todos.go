package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type Todo struct {
	Id          string
	Description string
	IsComplete  bool
	Started     time.Time
	Finished    time.Time
}

func printToDoTable(todos []Todo) {
	headerFormat := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFormat := color.New(color.FgYellow).SprintfFunc()

	table := table.New("ID", "Description", "Status", "Started", "Finished")
	table.WithHeaderFormatter(headerFormat).WithFirstColumnFormatter(columnFormat)

	for _, todo := range todos {
		table.AddRow(todo.Id, todo.Description, todo.IsComplete, todo.Started, todo.Finished)
	}

	table.Print()
}

func saveToDo(db *sql.DB, todo Todo) {
	stmnt, err := db.Prepare("INSERT INTO todos(id, description, iscomplete, started, finished) values (?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	res, err := stmnt.Exec(todo.Id, todo.Description, todo.IsComplete, todo.Started, todo.Finished)
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
	rows, err := db.Query("SELECT id, description, iscomplete, started, finished FROM todos")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var todo Todo
	for rows.Next() {
		err := rows.Scan(&todo.Id, &todo.Description, &todo.IsComplete, &todo.Started, &todo.Finished)
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
