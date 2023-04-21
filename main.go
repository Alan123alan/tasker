package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	DB, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getAll := getCmd.Bool("all", false, "Get all tasks")
	getId := getCmd.String("id", "", "Get task by id")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDescription := addCmd.String("desc", "", "Task description")
	addId := addCmd.String("id", "", "Task Id")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateDescription := updateCmd.String("desc", "", "Task description")
	updateStatus := updateCmd.Int("status", 0, "Task status")
	updateId := updateCmd.String("id", "", "Task id")

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'get', 'add', 'update' or 'create' commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		HandleGet(DB, getCmd, getAll, getId)
	case "add":
		HandleAdd(DB, addCmd, addId, addDescription)
	case "update":
		HandleUpdate(DB, updateCmd, updateId, updateDescription, (*Status)(updateStatus))
	case "create":
		HandleCreate(DB, createCmd)
	case "help":
	default:
		fmt.Printf("'%v' is not a valid command. See './tasker -help'.", os.Args[1])

	}
}
