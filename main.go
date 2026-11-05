package main

import "fmt"

// Task is a single to-do item.
// A struct lets us group related data together under one name.
type Task struct {
	Text string // what the task says
	Done bool   // true if finished, false if not
}

// tasks holds every to-do in memory while the program runs.
// It starts empty and is wiped when the program quits (Phase 1).
var tasks []Task

// addTask appends a new task to the slice.
func addTask(text string) {
	tasks = append(tasks, Task{Text: text, Done: false})
	fmt.Println("Added:", text)
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
		// i starts at 0, so we show i+1 for a human-friendly number
		fmt.Printf("%d. [%s] %s\n", i+1, marker, t.Text)
	}
}

// doneTask marks task number `num` as finished.
// The real slice index is num-1 because humans count from 1.
func doneTask(num int) {
	if num < 1 || num > len(tasks) {
		fmt.Println("Invalid task number.")
		return
	}
	tasks[num-1].Done = true
	fmt.Println("Marked done:", tasks[num-1].Text)
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

func main() {
	for {
		printMenu()

		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Please enter a number.")
			// throw away the bad input so the loop doesn't get stuck
			var dummy string
			fmt.Scan(&dummy)
			continue
		}

		if choice == 4 {
			fmt.Println("Goodbye!")
			break // exit the for loop, ending the program
		}

		if choice == 1 {
			fmt.Print("Task text: ")
			var text string
			fmt.Scan(&text)
			addTask(text)
		} else if choice == 2 {
			listTasks()
		} else if choice == 3 {
			fmt.Print("Task number: ")
			var num int
			fmt.Scan(&num)
			doneTask(num)
		} else {
			fmt.Println("Invalid choice.")
		}
	}
}
