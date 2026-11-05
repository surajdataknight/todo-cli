package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// errCommandFailed means a command already printed a helpful message.
// main uses it only to choose the process exit code.
var errCommandFailed = errors.New("command failed")

func addTask(list, text string, priority Priority) error {
	// lists[list] returns a nil slice when the list does not exist yet,
	// and append works on a nil slice, so this creates the list automatically.
	lists[list] = append(lists[list], Task{Text: text, Done: false, Priority: priority})
	fmt.Printf("Added to %q: %s (priority: %s)\n", list, text, priority)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
		return err
	}
	return nil
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

func doneTask(list string, num int) error {
	tasks := lists[list]
	if num < 1 || num > len(tasks) {
		fmt.Println("Invalid task number.")
		return errCommandFailed
	}
	tasks[num-1].Done = true
	lists[list] = tasks // write back so the map holds the updated slice
	fmt.Printf("Marked done in %q: %s\n", list, tasks[num-1].Text)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
		return err
	}
	return nil
}

// deleteTask removes task number `num` from `list` and saves to disk.
// tasks[:num-1] is everything before the removed task,
// tasks[num:] is everything after it, and append joins them back together.
func deleteTask(list string, num int) error {
	tasks := lists[list]
	if num < 1 || num > len(tasks) {
		fmt.Println("Invalid task number.")
		return errCommandFailed
	}
	deleted := tasks[num-1].Text
	tasks = append(tasks[:num-1], tasks[num:]...)
	lists[list] = tasks // write back because append may change the slice's length
	fmt.Printf("Deleted from %q: %s\n", list, deleted)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
		return err
	}
	return nil
}

// editTask changes the text of task number `num` in `list` to `newText`.
// The slice element is modified in place.
func editTask(list string, num int, newText string) error {
	tasks := lists[list]
	if num < 1 || num > len(tasks) {
		fmt.Println("Invalid task number.")
		return errCommandFailed
	}
	oldText := tasks[num-1].Text
	tasks[num-1].Text = newText
	lists[list] = tasks // write back so the map holds the updated slice
	fmt.Printf("Edited task %d in %q:\n  before: %s\n  after:  %s\n", num, list, oldText, newText)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
		return err
	}
	return nil
}

// usage prints how to use the program.
func usage() {
	fmt.Println("To-Do List CLI")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  todo add <list> \"Task text\" [low|medium|high]   Add a new task to a list")
	fmt.Println("  todo list <list>                                 List all tasks in a list")
	fmt.Println("  todo done <list> <number>                        Mark a task as done")
	fmt.Println("  todo delete <list> <number>                      Remove a task")
	fmt.Println("  todo edit <list> <number> \"new text\"             Change a task's text")
	fmt.Println("  todo lists                                       List all lists")
	fmt.Println("  todo version                                     Show the version")
	fmt.Println("  todo help                                        Show this help")
	fmt.Println("")
	fmt.Println("Priority defaults to medium when not given.")
	fmt.Println("Example: todo add Work \"Buy milk\" high")
}

// run executes one command and returns a non-nil error when it fails.
// main and the interactive prompt both call this same function.
func run(args []string) error {
	if len(args) == 0 {
		usage()
		return nil
	}

	// switch picks which command to run based on the first argument.
	switch command := args[0]; command {
	case "add":
		if len(args) < 3 {
			fmt.Println("Error: add needs a list name and task text.")
			fmt.Println(`Example: todo add Work "Buy milk"`)
			return errCommandFailed
		}
		list := args[1]
		text := strings.TrimSpace(args[2])
		if text == "" {
			fmt.Println("Error: task text cannot be empty.")
			return errCommandFailed
		}
		// Optional priority; defaults to Medium.
		priority := Medium
		if len(args) >= 4 {
			p, ok := parsePriority(args[3])
			if !ok {
				fmt.Println("Error: unknown priority:", args[3])
				fmt.Println(`Example: todo add Work "Buy milk" high`)
				return errCommandFailed
			}
			priority = p
		}
		return addTask(list, text, priority)

	case "list":
		if len(args) < 2 {
			fmt.Println("Error: list needs a list name.")
			fmt.Println("Example: todo list Work")
			return errCommandFailed
		}
		listTasks(args[1])
		return nil

	case "done":
		if len(args) < 3 {
			fmt.Println("Error: done needs a list name and task number.")
			fmt.Println("Example: todo done Work 1")
			return errCommandFailed
		}
		list := args[1]
		num, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Error: task number must be a number.")
			return errCommandFailed
		}
		return doneTask(list, num)

	case "delete":
		if len(args) < 3 {
			fmt.Println("Error: delete needs a list name and task number.")
			fmt.Println("Example: todo delete Work 1")
			return errCommandFailed
		}
		list := args[1]
		num, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Error: task number must be a number.")
			return errCommandFailed
		}
		return deleteTask(list, num)

	case "edit":
		if len(args) < 4 {
			fmt.Println("Error: edit needs a list name, task number, and new text.")
			fmt.Println(`Example: todo edit Work 1 "Buy milk (2%)"`)
			return errCommandFailed
		}
		list := args[1]
		num, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Error: task number must be a number.")
			return errCommandFailed
		}
		newText := strings.TrimSpace(args[3])
		if newText == "" {
			fmt.Println("Error: new text cannot be empty.")
			return errCommandFailed
		}
		return editTask(list, num, newText)

	case "lists":
		listLists()
		return nil

	case "version":
		fmt.Println("todo", version)
		return nil

	case "help":
		usage()
		return nil

	default:
		fmt.Println("Unknown command:", command)
		usage()
		return errCommandFailed
	}
}
