package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// lists maps a list name (e.g. "Work") to that list's tasks.
// Replacing the single tasks slice with a map lets us keep many named lists.
var lists map[string][]Task

// dataFile is the resolved path of the file that stores the tasks.
// main sets it once at startup (through initStorage) so every read and write
// in this process uses the same file, no matter which folder the command runs
// in.
var dataFile string

// dataPath returns the one file the CLI reads and writes.
//
// Before this existed, the program opened "todos.txt" relative to the current
// working directory. A globally installed todo therefore created a new, empty
// file in every folder it ran from. An absolute path fixes that:
//
//  1. TODO_FILE, when set, wins — handy for tests or for keeping data
//     somewhere else.
//  2. Otherwise <user config dir>/todo/todos.txt, for example
//     C:\Users\<you>\AppData\Roaming\todo\todos.txt on Windows.
func dataPath() (string, error) {
	if override := strings.TrimSpace(os.Getenv("TODO_FILE")); override != "" {
		// Abs keeps `todo where` showing a stable, full path even when the
		// user supplies a relative TODO_FILE.
		return filepath.Abs(override)
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "todo", "todos.txt"), nil
}

// initStorage resolves the data path and, on the first run with the new
// layout, imports tasks left behind by an older version.
func initStorage() error {
	path, err := dataPath()
	if err != nil {
		return err
	}
	dataFile = path

	// Importing old tasks is a convenience. If it fails, keep going so the
	// user can still work with the new file instead of being locked out.
	if err := migrateLegacyTasks(); err != nil {
		fmt.Println("Warning: could not import old tasks:", err)
	}
	return nil
}

// migrateLegacyTasks copies an old todos.txt into the new data file once.
// It only runs when the new file does not exist yet, so it can never
// overwrite tasks the user already has.
func migrateLegacyTasks() error {
	if _, err := os.Stat(dataFile); err == nil {
		return nil // data file already exists — nothing to import
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// Older versions wrote todos.txt next to wherever they were run from.
	// Check the current folder first, then the folder holding the executable
	// (the Phase 1 binary may have been run from there).
	candidates := []string{"todos.txt"}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "todos.txt"))
	}

	for _, candidate := range candidates {
		data, err := os.ReadFile(candidate)
		if errors.Is(err, os.ErrNotExist) {
			continue // try the next location
		}
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(dataFile), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dataFile, data, 0o644); err != nil {
			return err
		}
		fmt.Printf("Imported %d task(s) from %s\n", countTaskLines(data), candidate)
		return nil
	}
	return nil // no old file found — normal first run
}

// countTaskLines counts the lines loadTasks would accept.
func countTaskLines(data []byte) int {
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		if _, _, ok := parseTaskLine(strings.TrimSuffix(line, "\r")); ok {
			count++
		}
	}
	return count
}

// parseTaskLine splits one stored line: Work|true|high|Buy milk.
// SplitN with 4 keeps the task text intact even if it contains "|".
func parseTaskLine(line string) (string, Task, bool) {
	parts := strings.SplitN(line, "|", 4)
	if len(parts) != 4 {
		return "", Task{}, false
	}
	done, _ := strconv.ParseBool(parts[1])
	priority, _ := parsePriority(parts[2])
	return parts[0], Task{Text: parts[3], Done: done, Priority: priority}, true
}

// loadTasks reads the data file and rebuilds the in-memory map.
// If the file doesn't exist yet, we simply start empty.
func loadTasks() {
	lists = make(map[string][]Task)
	file, err := os.Open(dataFile)
	if err != nil {
		// No saved file yet — that's fine, first run.
		return
	}
	defer file.Close() // make sure the file closes when we leave this function

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		listName, task, ok := parseTaskLine(scanner.Text())
		if !ok {
			continue // skip malformed lines
		}
		lists[listName] = append(lists[listName], task)
	}
}

// saveTasks writes the current in-memory lists to the data file.
func saveTasks() error {
	// The todo folder may not exist on the very first save.
	if err := os.MkdirAll(filepath.Dir(dataFile), 0o755); err != nil {
		return err
	}

	file, err := os.Create(dataFile) // creates or overwrites
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
