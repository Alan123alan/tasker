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
		fmt.Println("Id\t\t\tDescription\t\t\tStatus\t\t\tStarted\t\t\tFinished")
		for _, todo := range todos {
			fmt.Printf("%v\t\t\t%v\t\t\t%v\t\t\t%v\t\t\t%v\n", todo.Id, todo.Description, todo.IsComplete, todo.Started, todo.Finished)
		}

		return
	}

	if *id != "" {
		todos := getTodos(db)
		id := *id
		for _, todo := range todos {
			if todo.Id == id {
				fmt.Println("Id\t\t\tDescription\t\t\tStatus\t\t\tStarted\t\t\tFinished")
				fmt.Printf("%v\t\t\t%v\t\t\t%v\t\t\t%v\t\t\t%v\n", todo.Id, todo.Description, todo.IsComplete, todo.Started, todo.Finished)
			}
		}
	}
}

func HandleAdd(db *sql.DB, addCmd *flag.FlagSet, id *string, description *string) {
	addCmd.Parse(os.Args[2:])
	ValidateToDo(addCmd, description)
	todo := Todo{*id, *description, false, time.Now(), time.Now()}
	fmt.Println(todo)
	// videos := getVideos()
	// videos = append(videos, video)
	saveToDo(db, todo)
}

func ValidateToDo(addCmd *flag.FlagSet, desc *string) {
	if *desc == "" {
		fmt.Println("Description is required to add a To do")
		addCmd.PrintDefaults()
		os.Exit(1)
	}
}

func HandleHelp() {

}

func main() {
	db, err := sql.Open("sqlite", "./myDb.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// db.Exec("CREATE TABLE todos( id int NOT NULL, description string, iscomplete bit, started datetime, finished datetime, PRIMARY KEY (id));")
	// db.Exec("INSERT INTO videos(id, title, description, image_url, url) values (?,?,?,?,?)")
	// stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	// var version string
	// err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	// var video video
	// err = db.QueryRow("SELECT * FROM videos").Scan(&video)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(version)
	// fmt.Println(video)
	// fmt.Println(stmt)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getAll := getCmd.Bool("all", false, "Get all to do")
	getId := getCmd.String("id", "", "Get to do by id")

	// addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	// addId := addCmd.String("id", "", "todo Id")
	// addTitle := addCmd.String("title", "", "To do title")
	// addDescription := addCmd.String("desc", "", "To do description")
	// addUrl := addCmd.String("url", "", "To do URL")
	// addImageUrl := addCmd.String("imageurl", "", "Video image URL")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDescription := addCmd.String("desc", "", "To do description")
	addId := addCmd.String("id", "", "todo Id")

	if len(os.Args) < 2 {
		fmt.Println("expected 'get' or 'add' commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		HandleGet(db, getCmd, getAll, getId)
	case "add":
		HandleAdd(db, addCmd, addId, addDescription)
	case "help":
		HandleHelp()
	default:
		fmt.Printf("'%v' is not a valid command. See './go-basic-cli -help'.", os.Args[1])

	}
}
