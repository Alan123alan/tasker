package main

import (
	"database/sql"
	"sort"

	// "flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "modernc.org/sqlite"
)

const (
	Get    string = "get"
	Add    string = "add"
	Create string = "create"
	Update string = "update"
)

// type Commands map[Command]Subcommands
// type Subcommands []Subcommand
// type Command string
// type Subcommand string

type Model struct {
	commands map[string][]string
	cursor   int
	selected map[int]struct{}
	DB       *sql.DB
	tasks    []TaskModel
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
	commands := make(map[string][]string)
	commands[Get] = []string{"All", "ID"}
	commands[Add] = []string{}
	commands[Update] = []string{}
	commands[Create] = []string{"table", "task"}
	return Model{
		commands: commands,
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
			if m.cursor < len(m.commands)-1 {
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

	return mainView(m)

}

func mainView(m Model) string {
	//Ask user wich command does it want to execute
	s := "Which command do you want to execute?\n\n"

	//Get list of available commands and sort them alphabetically
	commands := make([]string, 0, len(m.commands))
	for command := range m.commands {
		commands = append(commands, command)
	}
	sort.Strings(commands)

	// Iterate over our commands
	for index, command := range commands {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == index {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := "+" // not selected
		if _, ok := m.selected[index]; ok {
			checked = "-" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %v\n", cursor, checked, command)
		//if the command is selected and there is vailable subcommands, display them
		if checked == "-" {
			// subcommands := m.commands[command]
			// sort.Strings(subcommands)
			// for _, subcommand := range subcommands {
			// 	s += fmt.Sprintf("    - %v\n", subcommand)
			// }
			switch command {
			case Create:
				return createCommandView()
			case Get:
				return getCommandView(m)
			}
		}
	}
	// The footer
	s += "\nPress q to quit.\n"
	return s
}

func createCommandView() string {
	s := "Which table do you want to create in database?\n\n"
	s += "\nPress q to quit.\n"
	return s
}

func addCommandView() string {
	s := "Enter new task:\n\n"
	s += "\nPress q to quit.\n"
	return s
}

func getCommandView(m Model) string {
	tasks := getTasks(m.DB)
	s := "Get list of tasks:\n\n"
	for _, task := range tasks {
		s += fmt.Sprintf("%s [%s] %v\n", task.Id, task.Description, task.Status)
	}
	s += "\nPress q to quit.\n"
	return s
}

func updateCommandView() string {
	s := "Update existing task:\n\n"
	s += "\nPress q to quit.\n"
	return s
}
