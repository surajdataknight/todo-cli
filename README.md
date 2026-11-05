# todo-cli

A beginner-friendly command-line To-Do List written in Go.

It supports multiple named lists (e.g. `Work`, `Personal`), task priorities, and file persistence — all driven by `os.Args` commands with `switch` dispatch.

## Features

- **One command** — install once, then run `todo` from any folder
- **Multiple named lists** — tasks are grouped under a list name you choose
- **View all lists** — `lists` shows every list name and task count
- **Priorities** — `low`, `medium`, `high` (via `iota` constants), sorted high→low when listed
- **File persistence** — tasks are saved to `todos.txt` and reloaded on each run (`todos.txt` is gitignored as user data)
- **Commands** — `add`, `list`, `lists`, `done`, `delete`, `edit`; every command takes a list name

## Prerequisites

Go 1.22 or later.

## Install

```powershell
go install github.com/surajdataknight/todo-cli/cmd/todo@latest
```

This downloads, builds, and copies the `todo` executable into your Go bin directory — `go env GOBIN` if it is set, otherwise `$(go env GOPATH)\bin`. Make sure that directory is on your `PATH`, then confirm the install:

```powershell
todo version
```

To build from a local clone instead:

```powershell
git clone https://github.com/surajdataknight/todo-cli.git
cd todo-cli
go install ./cmd/todo
```

## Usage

```text
todo add <list> "Task text" [low|medium|high]   Add a new task to a list
todo list <list>                                 List all tasks in a list
todo lists                                       List all lists with task counts
todo done <list> <number>                        Mark a task as done
todo delete <list> <number>                      Remove a task
todo edit <list> <number> "new text"             Change a task's text
todo version                                     Show the version
todo help                                        Show this help
```

Priority defaults to `medium` when not given.

### Examples

```powershell
# Add tasks to different lists
todo add Work "Buy milk" high
todo add Work "Email boss"
todo add Personal "Gym"

# List a list (sorted by priority)
todo list Work

# Mark task 1 done, delete task 2, edit task 3's text
todo done Work 1
todo delete Work 2
todo edit Personal 1 "Gym workout"

# See every list and how many tasks are in each
todo lists

# Check which version is installed
todo version
```

Run `todo help` (or `todo` with no arguments) to see the usage message.

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
