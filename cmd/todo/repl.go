package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// startREPL is the interactive counterpart to the one-shot CLI. It reads
// commands from the terminal until the user types exit/quit or closes the
// input, and calls the same run() dispatcher, so both modes understand
// exactly the same commands.
func startREPL() {
	fmt.Println(`Interactive mode. Type "help" for commands, "exit" to quit.`)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("todo> ")
		if !scanner.Scan() {
			// Ctrl+D / Ctrl+Z, or the end of piped input.
			fmt.Println()
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
			}
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // ignore blank lines, just show the prompt again
		}

		args, err := splitArgs(line)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		if len(args) == 0 {
			continue
		}
		if args[0] == "exit" || args[0] == "quit" {
			return
		}

		// A failed command prints its own message; the prompt stays alive.
		_ = run(args)
	}
}
