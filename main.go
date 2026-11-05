package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Task is a single to-do item.
type Task struct {
	Text string // what the task says
	Done bool   // true if finished, false if not
}

// tasks holds every to-do in memory while the program runs.
var tasks []Task

// filename is where we persist tasks so they survive closing the program.
const filename = "todos.txt"

func addTask(text string) {
	tasks = append(tasks, Task{Text: text, Done: false})
	fmt.Println("Added:", text)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
	}
}

func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks yet.")
		return
	}

	fmt.Println("\nYour tasks:")
	for i, t := range tasks {
		marker := " "
		if t.Done {
			marker = "x"
		}
		fmt.Printf("%d. [%s] %s\n", i+1, marker, t.Text)
	}
}

func doneTask(num int) {
	if num < 1 || num > len(tasks) {
		fmt.Println("Invalid task number.")
		return
	}
	tasks[num-1].Done = true
	fmt.Println("Marked done:", tasks[num-1].Text)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
	}
}

// loadTasks reads todos.txt and rebuilds the in-memory slice.
// If the file doesn't exist yet, we simply start empty.
func loadTasks() {
	file, err := os.Open(filename)
	if err != nil {
		// No saved file yet — that's fine, first run.
		return
	}
	defer file.Close() // make sure the file closes when we leave this function

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Each line looks like: true|Buy milk
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue // skip malformed lines
		}
		done, _ := strconv.ParseBool(parts[0])
		text := parts[1]
		tasks = append(tasks, Task{Text: text, Done: done})
	}
}

// saveTasks writes the current in-memory tasks to todos.txt.
func saveTasks() error {
	file, err := os.Create(filename) // creates or overwrites
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, t := range tasks {
		line := strconv.FormatBool(t.Done) + "|" + t.Text + "\n"
		if _, err := writer.WriteString(line); err != nil {
			return err
		}
	}
	return writer.Flush() // push buffered data to disk
}

// usage prints how to run the program.
func usage() {
	fmt.Println("To-Do List CLI")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run main.go add \"Task text\"   Add a new task")
	fmt.Println("  go run main.go list              List all tasks")
	fmt.Println("  go run main.go done <number>     Mark a task as done")
}

func main() {
	// Load any saved tasks before doing anything.
	loadTasks()

	// os.Args holds the command-line arguments:
	//   os.Args[0] is the program name,
	//   os.Args[1] is the command (add/list/done),
	//   os.Args[2] and beyond are the command's arguments.
	args := os.Args[1:]

	if len(args) == 0 {
		usage()
		return
	}

	// switch picks which command to run based on the first argument.
	switch command := args[0]; command {
	case "add":
		if len(args) < 2 {
			fmt.Println("Error: add needs task text.")
			fmt.Println(`Example: go run main.go add "Buy milk"`)
			return
		}
		addTask(args[1])

	case "list":
		listTasks()

	case "done":
		if len(args) < 2 {
			fmt.Println("Error: done needs a task number.")
			fmt.Println("Example: go run main.go done 1")
			return
		}
		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: task number must be a number.")
			return
		}
		doneTask(num)

	default:
		fmt.Println("Unknown command:", command)
		usage()
	}
}
