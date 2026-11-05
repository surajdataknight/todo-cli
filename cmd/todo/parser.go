package main

import (
	"errors"
	"strings"
	"unicode"
)

// splitArgs turns one line of interactive input into the same []string
// shape that os.Args provides. It is only used by the REPL; one-shot
// commands keep using the shell's own argument parsing.
//
// Rules:
//   - whitespace separates arguments
//   - "double quotes" and 'single quotes' group text containing spaces
//   - inside double quotes, a backslash escapes the next character
//   - an unmatched quote is an error
func splitArgs(line string) ([]string, error) {
	var args []string
	var current strings.Builder
	inSingle, inDouble := false, false
	escaped := false
	hasToken := false

	flush := func() {
		if hasToken {
			args = append(args, current.String())
			current.Reset()
			hasToken = false
		}
	}

	for _, r := range line {
		if inDouble {
			if escaped {
				current.WriteRune(r)
				escaped = false
				hasToken = true
				continue
			}
			switch r {
			case '\\':
				escaped = true
			case '"':
				inDouble = false
			default:
				current.WriteRune(r)
				hasToken = true
			}
			continue
		}

		if inSingle {
			if r == '\'' {
				inSingle = false
			} else {
				current.WriteRune(r)
				hasToken = true
			}
			continue
		}

		if unicode.IsSpace(r) {
			flush()
			continue
		}

		switch r {
		case '"':
			inDouble = true
			hasToken = true
		case '\'':
			inSingle = true
			hasToken = true
		default:
			current.WriteRune(r)
			hasToken = true
		}
	}

	// A backslash at the very end can only have come from inside double
	// quotes, so it also means the closing quote never arrived.
	if escaped || inDouble || inSingle {
		return nil, errors.New("unmatched quote")
	}

	flush()
	return args, nil
}
