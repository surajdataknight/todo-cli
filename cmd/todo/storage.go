package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// lists maps a list name (e.g. "Work") to that list's tasks.
// Replacing the single tasks slice with a map lets us keep many named lists.
var lists map[string][]Task

// filename is where we persist tasks so they survive closing the program.
const filename = "todos.txt"

// loadTasks reads todos.txt and rebuilds the in-memory map.
// Each line looks like: Work|true|high|Buy milk   (list|done|priority|text)
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
