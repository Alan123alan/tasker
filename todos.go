package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type video struct {
	Id          string
	Title       string
	Description string
	Imageurl    string
	Url         string
}

type Todo struct {
	Id          string
	Description string
	IsComplete  bool
	Started     time.Time
	Finished    time.Time
}

func getVideos() (videos []video) {
	videoBytes, err := ioutil.ReadFile("./videos.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(videoBytes, &videos)
	if err != nil {
		panic(err)
	}
	return videos
}

func saveVideos(videos []video) {
	videoBytes, err := json.Marshal(videos)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("./videos.json", videoBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func saveVideoToDB(db *sql.DB, video video) {
	db.Exec("INSERT INTO videos(id, title, description, image_url, url) values (?,?,?,?,?)", video.Id, video.Title, video.Description, video.Imageurl, video.Url)
}

func getVideosFromDB(db *sql.DB) (videos []video) {
	rows, err := db.Query("SELECT id, title, description, image_url, url FROM videos")
	if err != nil {
		panic(err)
	}
	var video video
	for rows.Next() {
		rows.Scan(&video.Id, &video.Title, &video.Description, &video.Imageurl, &video.Url)
		videos = append(videos, video)
	}
	fmt.Println(videos)
	return videos
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
	// db.Exec("INSERT INTO todos(description, ispending, iscompleted, started, finished) values (?,?,?,?,?)", todo.Description, todo.IsComplete, todo.Started, todo.Finished)
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
