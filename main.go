package main

import (
	"database/sql"
	// "flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "modernc.org/sqlite"
)

type Model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	DB       *sql.DB
	tasks    taskMsg
	err      error
}

func main() {
	DB, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	p := tea.NewProgram(initialModel(DB))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	// getAll := getCmd.Bool("all", false, "Get all tasks")
	// getId := getCmd.String("id", "", "Get task by id")

	// addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	// addDescription := addCmd.String("desc", "", "Task description")
	// addId := addCmd.String("id", "", "Task Id")

	// updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	// updateDescription := updateCmd.String("desc", "", "Task description")
	// updateStatus := updateCmd.Int("status", 0, "Task status")
	// updateId := updateCmd.String("id", "", "Task id")

	// createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	// if len(os.Args) < 2 {
	// 	fmt.Println("expected 'get', 'add', 'update' or 'create' commands")
	// 	os.Exit(1)
	// }

	// switch os.Args[1] {
	// case "get":
	// 	HandleGet(DB, getCmd, getAll, getId)
	// case "add":
	// 	HandleAdd(DB, addCmd, addId, addDescription)
	// case "update":
	// 	HandleUpdate(DB, updateCmd, updateId, updateDescription, (*Status)(updateStatus))
	// case "create":
	// 	HandleCreate(DB, createCmd)
	// case "help":
	// default:
	// 	fmt.Printf("'%v' is not a valid command. See './tasker -help'.", os.Args[1])

	// }
}

func initialModel(DB *sql.DB) Model {
	return Model{
		choices:  []string{"get", "add", "update", "create", "help"},
		selected: make(map[int]struct{}),
		DB:       DB,
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
	// return getTasks(m.DB)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case taskMsg:
		// The server returned a status message. Save it to our model. Also
		// tell the Bubble Tea runtime we want to exit because we have nothing
		// else to do. We'll still be able to render a final view with our
		// status message.
		m.tasks = taskMsg(msg)
		return m, tea.Quit

	case errMsg:
		// There was an error. Note it in the model. And tell the runtime
		// we're done and want to quit.
		m.err = msg
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
				m.tasks = getTasks(m.DB)
				return m, nil
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	// If there's an error, print it out and don't do anything else.
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	// Tell the user we're doing something.
	// s := "Checking tasks database table ... "

	// Send off whatever we came up with above for rendering.
	// return "\n" + s + "\n\n"
	// The header
	s := "Which operation do you want to perform?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// When the database responds with the tasks, add it to the current line.
	if len(m.tasks) > 0 {
		for _, task := range m.tasks {
			s += fmt.Sprintf("\n%v | %v | %v\n", task.Id, task.Description, task.Status)
		}
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
