package main

import (
	"fmt"
	"os"
	"runtime/debug"
)

// version is the human-readable build version.
// Release builds can replace it with:
//
//	go build -ldflags "-X main.version=0.1.0" ./cmd/todo
var version = "dev"

// buildVersion returns the version to show to the user. An explicit build
// flag wins. Otherwise fall back to the module version that `go install`
// records in the binary, so an installed release does not report "dev".
func buildVersion() string {
	if version != "dev" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return version
	}

	// A build from a local clone records its VCS revision; keep calling that
	// "dev" so uncommitted work is not mistaken for a published release.
	// Downloads through `go install ...@version` have no VCS metadata and
	// carry the module's real version instead.
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return version
		}
	}
	return info.Main.Version
}

func main() {
	// Work out where tasks live before reading them. On the first run this
	// also imports a todos.txt left behind by an older version.
	if err := initStorage(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: could not find the todo data folder:", err)
		os.Exit(1)
	}

	// Load any saved lists before doing anything.
	loadTasks()

	// With no arguments, open the interactive prompt when the command was
	// typed at a terminal. If standard input is a pipe or a redirected file
	// (for example `todo < commands.txt`), print the usage instead so the
	// program never waits for input that will not arrive.
	if len(os.Args) == 1 {
		if stdinIsTerminal() {
			startREPL()
			return
		}
		_ = run(nil)
		return
	}

	// os.Args holds the command-line arguments:
	//   os.Args[0] is the program name,
	//   os.Args[1] is the command (add/list/done),
	//   os.Args[2] and beyond are the command's arguments.
	//
	// A non-nil error means the command already printed a helpful message;
	// exit with a non-zero status so scripts can detect the failure.
	if err := run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}

// stdinIsTerminal reports whether standard input is connected to a terminal.
// A pipe or a redirected file does not have the ModeCharDevice bit set, which
// is exactly when the interactive prompt must stay closed.
func stdinIsTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
