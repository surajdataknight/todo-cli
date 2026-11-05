# todo-cli

A beginner-friendly command-line To-Do List written in Go.

It supports multiple named lists (e.g. `Work`, `Personal`), task priorities, and file persistence. Every command works both as a one-shot `todo ...` call and at an interactive `todo>` prompt.

## Features

- **One command** — install once, then run `todo` from any folder
- **Interactive mode** — run `todo` on its own for a `todo>` prompt that accepts the same commands
- **One data file** — tasks live in your user config folder, independent of the current directory
- **Multiple named lists** — tasks are grouped under a list name you choose
- **View all lists** — `lists` shows every list name and task count
- **Priorities** — `low`, `medium`, `high` (via `iota` constants), sorted high→low when listed
- **File persistence** — tasks are saved to a plain text file and reloaded on each run
- **Commands** — `add`, `list`, `lists`, `done`, `delete`, `edit`, plus `where`, `version`, and `help`; every task command takes a list name

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
todo                                             Start the interactive prompt
todo add <list> "Task text" [low|medium|high]   Add a new task to a list
todo list <list>                                 List all tasks in a list
todo lists                                       List all lists with task counts
todo all                                         List all tasks in every list
todo where                                       Show the data file and task count
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

# See every list and every task in each list
todo all

# Show the data file and how many tasks are stored
todo where

# Check which version is installed
todo version
```

Run `todo help` to see the usage message.

Task numbers are the ones printed by `todo list`. They stay the same between
runs: tasks are shown highest priority first, equal priorities keep the order
you added them in, and the data file itself stays in insertion order.

### Interactive mode

Run `todo` with no arguments in a terminal to get a reusable prompt:

```text
$ todo
Interactive mode. Type "help" for commands, "exit" to quit.
todo> add Work "Buy milk" high
Added to "Work": Buy milk (priority: high)
todo> list Work

Tasks in "Work":
1. [ ] Buy milk (priority: high)
Total: 1 task(s)
todo> exit
```

The prompt accepts exactly the same commands as the one-shot form, without the
leading `todo`. Double or single quotes keep text with spaces together, so
`add Work "buy milk (2%)" high` stores the full text. A mistyped command prints
an error and returns to the prompt; `exit`, `quit`, and end-of-input (Ctrl+D, or
Ctrl+Z then Enter on Windows) leave cleanly.

When input comes from a pipe or a redirected file (`todo < commands.txt`),
`todo` prints the help instead of opening the prompt, so it never waits for
input that will not arrive.

## Data location

Tasks are stored in a single file, so the folder you run `todo` from does not matter. The default is `todo/todos.txt` inside your user config directory:

| Platform | Default path |
|---|---|
| Windows | `%AppData%\todo\todos.txt` |
| macOS | `~/Library/Application Support/todo/todos.txt` |
| Linux | `~/.config/todo/todos.txt` |

Run `todo where` to print the exact path in use and the number of stored tasks.

To keep tasks somewhere else, set the `TODO_FILE` environment variable before running a command:

```powershell
$env:TODO_FILE = "D:\tmp\todos.txt"
todo where
```

On the first run with the new layout, if the default file does not exist yet, `todo` imports an existing `todos.txt` from the current folder or from the folder containing the executable and prints `Imported N task(s) from <path>`. An existing data file is never overwritten.

## Data format

Each line in the data file is `list|done|priority|text`, for example:

```text
Work|false|high|Buy milk
Personal|true|medium|Gym
```

The text may itself contain `|`; the first three separators are the ones that matter. The file is created automatically on the first `add`/`done`/`delete`/`edit`.

## Concepts demonstrated

- `struct`, `slice`, `map[string][]Task`
- `iota` constants and a `String()` method
- `sort.SliceStable` with a custom "less" function
- `os.Args` with `switch` command dispatch
- `os.Stdin.Stat` and `os.ModeCharDevice` to detect an interactive terminal
- a small quote-aware tokenizer built on `strings.Builder`
- `os.UserConfigDir`, `filepath.Join`, and `os.MkdirAll` for portable paths
- `strings.TrimSpace` input validation
- `bufio.Scanner`/`Writer`, `os.Create`/`Open`, `defer file.Close()`
- `strconv.Atoi`, `strconv.FormatBool`, `strconv.ParseBool`
