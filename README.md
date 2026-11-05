# todo-cli

A beginner-friendly command-line To-Do List written in Go.

It supports multiple named lists (e.g. `Work`, `Personal`), task priorities, and file persistence — all driven by `os.Args` commands with `switch` dispatch.

## Features

- **Multiple named lists** — tasks are grouped under a list name you choose
- **View all lists** — `lists` shows every list name and task count
- **Priorities** — `low`, `medium`, `high` (via `iota` constants), sorted high→low when listed
- **File persistence** — tasks are saved to `todos.txt` and reloaded on each run (`todos.txt` is gitignored as user data)
- **Commands** — `add`, `list`, `lists`, `done`, `delete`, `edit`; every command takes a list name

## Prerequisites

Go 1.22 or later.

## Build

```powershell
cd todo-cli
go build -o todo.exe ./cmd/todo
```

## Usage

```text
go run ./cmd/todo add <list> "Task text" [low|medium|high]   Add a new task to a list
go run ./cmd/todo list <list>                                 List all tasks in a list
go run ./cmd/todo lists                                       List all lists with task counts
go run ./cmd/todo done <list> <number>                        Mark a task as done
go run ./cmd/todo delete <list> <number>                      Remove a task
go run ./cmd/todo edit <list> <number> "new text"             Change a task's text
```

Priority defaults to `medium` when not given.

### Examples

```powershell
# Add tasks to different lists
go run ./cmd/todo add Work "Buy milk" high
go run ./cmd/todo add Work "Email boss"
go run ./cmd/todo add Personal "Gym"

# List a list (sorted by priority)
go run ./cmd/todo list Work

# Mark task 1 done, delete task 2, edit task 3's text
go run ./cmd/todo done Work 1
go run ./cmd/todo delete Work 2
go run ./cmd/todo edit Personal 1 "Gym workout"

# See every list and how many tasks are in each
go run ./cmd/todo lists
```

Run with no arguments to see the usage message.

## Data format

Each line in `todos.txt` is `list|done|priority|text`, for example:

```text
Work|false|high|Buy milk
Personal|true|medium|Gym
```

The file is created automatically on the first `add`/`done`/`delete`/`edit` and is gitignored so your task data stays out of source control.

## Concepts demonstrated

- `struct`, `slice`, `map[string][]Task`
- `iota` constants and a `String()` method
- `sort.Slice` with a custom "less" function
- `os.Args` with `switch` command dispatch
- `strings.TrimSpace` input validation
- `bufio.Scanner`/`Writer`, `os.Create`/`Open`, `defer file.Close()`
- `strconv.Atoi`, `strconv.FormatBool`, `strconv.ParseBool`
