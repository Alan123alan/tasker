package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

func HandleGet(db *sql.DB, getCmd *flag.FlagSet, all *bool, id *string) {
	getCmd.Parse(os.Args[2:])
	if *all == false && *id == "" {
		fmt.Println("id is required or specify --all for all to do")
		getCmd.PrintDefaults()
		os.Exit(1)
	}

	if *all {
		todos := getTodos(db)
		printToDoTable(todos)
		return
	}

	if *id != "" {
		todos := getTodo(db, *id)
		printToDoTable(todos)
	}
}

func HandleAdd(db *sql.DB, addCmd *flag.FlagSet, id *string, description *string) {
	addCmd.Parse(os.Args[2:])
	ValidateToDo(addCmd, description)
	todo := Todo{*id, *description, 0, time.Now(), time.Now()}
	saveToDo(db, todo)
}

func HandleUpdate(db *sql.DB, updateCmd *flag.FlagSet, id *string, description *string, status *Status) {
	updateCmd.Parse(os.Args[2:])
	validateUpdate(updateCmd, id)
	updateTodo(db, *id, *description, *status)
}

func ValidateToDo(addCmd *flag.FlagSet, desc *string) {
	if *desc == "" {
		fmt.Println("Description is required to add a To do")
		addCmd.PrintDefaults()
		os.Exit(1)
	}
}

func validateUpdate(updateCmd *flag.FlagSet, id *string) {
	if *id == "" {
		fmt.Println("id is required to update a To do")
		updateCmd.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	db, err := sql.Open("sqlite", "./myDb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// stmnt, err := db.Prepare(`CREATE TABLE todos (id INT PRIMARY KEY, description STRING, status INT, started_at DATETIME, completed_at DATETIME);`)
	// stmnt.Exec()

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getAll := getCmd.Bool("all", false, "Get all to do")
	getId := getCmd.String("id", "", "Get to do by id")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDescription := addCmd.String("desc", "", "To do description")
	addId := addCmd.String("id", "", "To do Id")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateDescription := updateCmd.String("desc", "", "To do description")
	updateStatus := updateCmd.Int("status", 0, "To do status")
	updateId := updateCmd.String("id", "", "To do id")

	if len(os.Args) < 2 {
		fmt.Println("expected 'get', 'add' or 'updated' commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		HandleGet(db, getCmd, getAll, getId)
	case "add":
		HandleAdd(db, addCmd, addId, addDescription)
	case "update":
		HandleUpdate(db, updateCmd, updateId, updateDescription, (*Status)(updateStatus))
	case "help":
	default:
		fmt.Printf("'%v' is not a valid command. See './go-basic-cli -help'.", os.Args[1])

	}
}
