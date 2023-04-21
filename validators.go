package main

import (
	"flag"
	"fmt"
	"os"
)

func validateNewTask(addCmd *flag.FlagSet, desc *string) {
	if *desc == "" {
		fmt.Println("Description is required to add a To do")
		addCmd.PrintDefaults()
		os.Exit(1)
	}
}

func validateTaskUpdate(updateCmd *flag.FlagSet, id *string) {
	if *id == "" {
		fmt.Println("id is required to update a Task")
		updateCmd.PrintDefaults()
		os.Exit(1)
	}
}
