package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Priority is how important a task is.
// iota gives each constant a unique number, starting at 0:
//
//	Low  = 0
//	Medium = 1
//	High = 2
type Priority int

const (
	Low Priority = iota
	Medium
	High
)

// String turns a Priority into a readable word (e.g. High -> "high").
func (p Priority) String() string {
	switch p {
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	default:
		return "unknown"
	}
}

// parsePriority turns a word like "high" into a Priority.
// The second return value tells whether the word was recognized.
func parsePriority(s string) (Priority, bool) {
	switch strings.ToLower(s) {
	case "low":
		return Low, true
	case "medium":
		return Medium, true
	case "high":
		return High, true
	default:
		return Medium, false
	}
}

// Task is a single to-do item.
type Task struct {
	Text     string    // what the task says
	Done     bool      // true if finished, false if not
	Priority Priority  // how important the task is (low/medium/high)
}

// tasks holds every to-do in memory while the program runs.
var tasks []Task

// filename is where we persist tasks so they survive closing the program.
const filename = "todos.txt"

func addTask(text string, priority Priority) {
	tasks = append(tasks, Task{Text: text, Done: false, Priority: priority})
	fmt.Printf("Added: %s (priority: %s)\n", text, priority)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
	}
}

func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks yet.")
		return
	}

	// sort.Slice orders the slice using a custom "less" function.
	// Here we put higher priority first (High before Medium before Low).
	// This sorts the slice in place, so the printed order matches the
	// order used by done/delete.
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})

	fmt.Println("\nYour tasks:")
	for i, t := range tasks {
		marker := " "
		if t.Done {
			marker = "x"
		}
		fmt.Printf("%d. [%s] %s (priority: %s)\n", i+1, marker, t.Text, t.Priority)
	}
	fmt.Printf("Total: %d task(s)\n", len(tasks))
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

// deleteTask removes task number `num` from the list and saves to disk.
// tasks[:num-1] is everything before the removed task,
// tasks[num:] is everything after it, and append joins them back together.
func deleteTask(num int) {
	if num < 1 || num > len(tasks) {
		fmt.Println("Invalid task number.")
		return
	}
	deleted := tasks[num-1].Text
	tasks = append(tasks[:num-1], tasks[num:]...)
	fmt.Println("Deleted:", deleted)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
	}
}

// loadTasks reads todos.txt and rebuilds the in-memory slice.
// If the file doesn't exist yet, we simply start empty.
// Each line now looks like: true|high|Buy milk   (done|priority|text)
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
		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			continue // skip malformed lines
		}
		done, _ := strconv.ParseBool(parts[0])
		priority, _ := parsePriority(parts[1])
		text := parts[2]
		tasks = append(tasks, Task{Text: text, Done: done, Priority: priority})
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
		line := strconv.FormatBool(t.Done) + "|" + t.Priority.String() + "|" + t.Text + "\n"
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
	fmt.Println("  go run main.go add \"Task text\" [low|medium|high]   Add a new task")
	fmt.Println("  go run main.go list                                 List all tasks")
	fmt.Println("  go run main.go done <number>                        Mark a task as done")
	fmt.Println("  go run main.go delete <number>                      Remove a task")
	fmt.Println("")
	fmt.Println("Priority defaults to medium when not given.")
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
		text := strings.TrimSpace(args[1])
		if text == "" {
			fmt.Println("Error: task text cannot be empty.")
			return
		}
		// Optional priority; defaults to Medium.
		priority := Medium
		if len(args) >= 3 {
			p, ok := parsePriority(args[2])
			if !ok {
				fmt.Println("Error: unknown priority:", args[2])
				fmt.Println(`Example: go run main.go add "Buy milk" high`)
				return
			}
			priority = p
		}
		addTask(text, priority)

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

	case "delete":
		if len(args) < 2 {
			fmt.Println("Error: delete needs a task number.")
			fmt.Println("Example: go run main.go delete 1")
			return
		}
		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: task number must be a number.")
			return
		}
		deleteTask(num)

	default:
		fmt.Println("Unknown command:", command)
		usage()
	}
}
