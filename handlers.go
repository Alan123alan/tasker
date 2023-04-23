package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"
)

func HandleGet(DB *sql.DB, getCmd *flag.FlagSet, all *bool, id *string) {
	getCmd.Parse(os.Args[2:])
	if *all == false && *id == "" {
		fmt.Println("id is required or specify --all for all to do")
		getCmd.PrintDefaults()
		os.Exit(1)
	}

	// if *all {
	// 	tasks := getTasks(DB)
	// 	printToDoTable(tasks)
	// 	return
	// }

	if *id != "" {
		tasks := getTask(DB, *id)
		printToDoTable(tasks)
	}
}

func HandleAdd(DB *sql.DB, addCmd *flag.FlagSet, id *string, description *string) {
	addCmd.Parse(os.Args[2:])
	validateNewTask(addCmd, description)
	task := Task{*description, 0, time.Now(), time.Now()}
	saveTask(DB, task)
}

// func HandleUpdate(DB *sql.DB, updateCmd *flag.FlagSet, id *string, description *string, status *Status) {
// 	updateCmd.Parse(os.Args[2:])
// 	validateTaskUpdate(updateCmd, id)
// 	updateTask(DB, *id, *description, *status)
// }

func HandleCreate(DB *sql.DB, createCmd *flag.FlagSet) {
	createTaskTable(DB)
}

func HandleHelp() {

}
