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

// lists maps a list name (e.g. "Work") to that list's tasks.
// Replacing the single tasks slice with a map lets us keep many named lists.
var lists map[string][]Task

// filename is where we persist tasks so they survive closing the program.
const filename = "todos.txt"

func addTask(list, text string, priority Priority) {
	// lists[list] returns a nil slice when the list does not exist yet,
	// and append works on a nil slice, so this creates the list automatically.
	lists[list] = append(lists[list], Task{Text: text, Done: false, Priority: priority})
	fmt.Printf("Added to %q: %s (priority: %s)\n", list, text, priority)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
	}
}

func listTasks(list string) {
	tasks := lists[list]
	if len(tasks) == 0 {
		fmt.Printf("No tasks in list %q yet.\n", list)
		return
	}

	// sort.Slice orders the slice using a custom "less" function.
	// Here we put higher priority first (High before Medium before Low).
	// This sorts the slice in place, so the printed order matches the
	// order used by done/delete/edit.
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})

	fmt.Printf("\nTasks in %q:\n", list)
	for i, t := range tasks {
		marker := " "
		if t.Done {
			marker = "x"
		}
		fmt.Printf("%d. [%s] %s (priority: %s)\n", i+1, marker, t.Text, t.Priority)
	}
	fmt.Printf("Total: %d task(s)\n", len(tasks))
}

// listLists prints every list name and how many tasks are in each.
func listLists() {
	if len(lists) == 0 {
		fmt.Println("No lists yet. Add a task to create one.")
		return
	}

	fmt.Println("All lists:")
	for listName, tasks := range lists {
		count := len(tasks)
		label := "tasks"
		if count == 1 {
			label = "task"
		}
		fmt.Printf("- %s (%d %s)\n", listName, count, label)
	}
}

func doneTask(list string, num int) {
	tasks := lists[list]
	if num < 1 || num > len(tasks) {
		fmt.Println("Invalid task number.")
		return
	}
	tasks[num-1].Done = true
	lists[list] = tasks // write back so the map holds the updated slice
	fmt.Printf("Marked done in %q: %s\n", list, tasks[num-1].Text)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
	}
}

// deleteTask removes task number `num` from `list` and saves to disk.
// tasks[:num-1] is everything before the removed task,
// tasks[num:] is everything after it, and append joins them back together.
func deleteTask(list string, num int) {
	tasks := lists[list]
	if num < 1 || num > len(tasks) {
		fmt.Println("Invalid task number.")
		return
	}
	deleted := tasks[num-1].Text
	tasks = append(tasks[:num-1], tasks[num:]...)
	lists[list] = tasks // write back because append may change the slice's length
	fmt.Printf("Deleted from %q: %s\n", list, deleted)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
	}
}

// editTask changes the text of task number `num` in `list` to `newText`.
// The slice element is modified in place.
func editTask(list string, num int, newText string) {
	tasks := lists[list]
	if num < 1 || num > len(tasks) {
		fmt.Println("Invalid task number.")
		return
	}
	oldText := tasks[num-1].Text
	tasks[num-1].Text = newText
	lists[list] = tasks // write back so the map holds the updated slice
	fmt.Printf("Edited task %d in %q:\n  before: %s\n  after:  %s\n", num, list, oldText, newText)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
	}
}

// loadTasks reads todos.txt and rebuilds the in-memory map.
// Each line now looks like: Work|true|high|Buy milk   (list|done|priority|text)
// SplitN with 4 keeps the task text intact even if it contains "|".
// If the file doesn't exist yet, we simply start empty.
func loadTasks() {
	lists = make(map[string][]Task)
	file, err := os.Open(filename)
	if err != nil {
		// No saved file yet — that's fine, first run.
		return
	}
	defer file.Close() // make sure the file closes when we leave this function

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue // skip malformed lines
		}
		listName := parts[0]
		done, _ := strconv.ParseBool(parts[1])
		priority, _ := parsePriority(parts[2])
		text := parts[3]
		lists[listName] = append(lists[listName], Task{Text: text, Done: done, Priority: priority})
	}
}

// saveTasks writes the current in-memory lists to todos.txt.
func saveTasks() error {
	file, err := os.Create(filename) // creates or overwrites
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for listName, tasks := range lists {
		for _, t := range tasks {
			line := listName + "|" + strconv.FormatBool(t.Done) + "|" + t.Priority.String() + "|" + t.Text + "\n"
			if _, err := writer.WriteString(line); err != nil {
				return err
			}
		}
	}
	return writer.Flush() // push buffered data to disk
}

// usage prints how to run the program.
func usage() {
	fmt.Println("To-Do List CLI")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run main.go add <list> \"Task text\" [low|medium|high]   Add a new task to a list")
	fmt.Println("  go run main.go list <list>                                 List all tasks in a list")
	fmt.Println("  go run main.go done <list> <number>                        Mark a task as done")
	fmt.Println("  go run main.go delete <list> <number>                      Remove a task")
	fmt.Println("  go run main.go edit <list> <number> \"new text\"             Change a task's text")
	fmt.Println("  go run main.go lists                                       List all lists")
	fmt.Println("")
	fmt.Println("Priority defaults to medium when not given.")
	fmt.Println("Example: go run main.go add Work \"Buy milk\" high")
}

func main() {
	// Load any saved lists before doing anything.
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
		if len(args) < 3 {
			fmt.Println("Error: add needs a list name and task text.")
			fmt.Println(`Example: go run main.go add Work "Buy milk"`)
			return
		}
		list := args[1]
		text := strings.TrimSpace(args[2])
		if text == "" {
			fmt.Println("Error: task text cannot be empty.")
			return
		}
		// Optional priority; defaults to Medium.
		priority := Medium
		if len(args) >= 4 {
			p, ok := parsePriority(args[3])
			if !ok {
				fmt.Println("Error: unknown priority:", args[3])
				fmt.Println(`Example: go run main.go add Work "Buy milk" high`)
				return
			}
			priority = p
		}
		addTask(list, text, priority)

	case "list":
		if len(args) < 2 {
			fmt.Println("Error: list needs a list name.")
			fmt.Println("Example: go run main.go list Work")
			return
		}
		listTasks(args[1])

	case "done":
		if len(args) < 3 {
			fmt.Println("Error: done needs a list name and task number.")
			fmt.Println("Example: go run main.go done Work 1")
			return
		}
		list := args[1]
		num, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Error: task number must be a number.")
			return
		}
		doneTask(list, num)

	case "delete":
		if len(args) < 3 {
			fmt.Println("Error: delete needs a list name and task number.")
			fmt.Println("Example: go run main.go delete Work 1")
			return
		}
		list := args[1]
		num, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Error: task number must be a number.")
			return
		}
		deleteTask(list, num)

	case "edit":
		if len(args) < 4 {
			fmt.Println("Error: edit needs a list name, task number, and new text.")
			fmt.Println(`Example: go run main.go edit Work 1 "Buy milk (2%)"`)
			return
		}
		list := args[1]
		num, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Error: task number must be a number.")
			return
		}
		newText := strings.TrimSpace(args[3])
		if newText == "" {
			fmt.Println("Error: new text cannot be empty.")
			return
		}
		editTask(list, num, newText)

	case "lists":
		listLists()

	default:
		fmt.Println("Unknown command:", command)
		usage()
	}
}
