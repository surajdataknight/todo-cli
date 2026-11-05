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

// addTask appends a new task and saves to disk immediately.
func addTask(text string) {
	tasks = append(tasks, Task{Text: text, Done: false})
	fmt.Println("Added:", text)
	if err := saveTasks(); err != nil {
		fmt.Println("Warning: could not save tasks:", err)
	}
}

// listTasks prints every task with a number and a status marker.
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

// doneTask marks task number `num` as finished and saves to disk.
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

// printMenu shows the options and prompts for a choice.
func printMenu() {
	fmt.Println("\n--- To-Do List ---")
	fmt.Println("1. Add task")
	fmt.Println("2. List tasks")
	fmt.Println("3. Mark task done")
	fmt.Println("4. Quit")
	fmt.Print("Choose: ")
}

// readLine reads one full line of input (supports spaces).
func readLine() string {
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(line)
}

func main() {
	// Load any saved tasks before showing the menu.
	loadTasks()

	for {
		printMenu()

		choiceStr := readLine()
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Please enter a number.")
			continue
		}

		if choice == 4 {
			fmt.Println("Goodbye!")
			break
		}

		if choice == 1 {
			fmt.Print("Task text: ")
			text := readLine()
			if text == "" {
				fmt.Println("Task text cannot be empty.")
				continue
			}
			addTask(text)
		} else if choice == 2 {
			listTasks()
		} else if choice == 3 {
			fmt.Print("Task number: ")
			numStr := readLine()
			num, err := strconv.Atoi(numStr)
			if err != nil {
				fmt.Println("Please enter a valid number.")
				continue
			}
			doneTask(num)
		} else {
			fmt.Println("Invalid choice.")
		}
	}
}
