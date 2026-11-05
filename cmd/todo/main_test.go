package main

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

// TestSplitArgs checks the quote-aware tokenizer that turns one line of
// interactive input into the same []string shape os.Args provides.
func TestSplitArgs(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    []string
		wantErr bool
	}{
		{
			name: "plain",
			line: "add Work Buy",
			want: []string{"add", "Work", "Buy"},
		},
		{
			name: "double quotes",
			line: `add Work "Buy milk" high`,
			want: []string{"add", "Work", "Buy milk", "high"},
		},
		{
			name: "single quotes",
			line: `add Work 'Buy milk' high`,
			want: []string{"add", "Work", "Buy milk", "high"},
		},
		{
			name: "escaped quotes",
			line: `add Work "say \"hi\"" high`,
			want: []string{"add", "Work", `say "hi"`, "high"},
		},
		{
			name: "escaped backslash",
			line: `add Work "C:\\tmp"`,
			want: []string{"add", "Work", `C:\tmp`},
		},
		{
			name: "unmatched double quote",
			line: `add Work "Buy milk`,
			// The whole line is rejected instead of guessing where it ends.
			wantErr: true,
		},
		{
			name:    "unmatched single quote",
			line:    `add Work 'Buy milk`,
			wantErr: true,
		},
		{
			name:    "trailing escape inside quotes",
			line:    `add Work "Buy milk\`,
			wantErr: true,
		},
		{
			name: "extra spaces",
			line: "  add   \tWork    \"Buy milk\"   ",
			want: []string{"add", "Work", "Buy milk"},
		},
		{
			name: "empty line",
			line: "",
		},
		{
			name: "quoted empty text",
			line: `add Work ""`,
			want: []string{"add", "Work", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitArgs(tt.line)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("splitArgs(%q) = %#v, want an error", tt.line, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("splitArgs(%q) returned an unexpected error: %v", tt.line, err)
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("splitArgs(%q) = %#v, want %#v", tt.line, got, tt.want)
			}
		})
	}
}

// TestStorageRoundTrip saves a set of lists to a temporary file, loads them
// again, and checks that the in-memory map is unchanged. TODO_FILE points at
// the temp folder so the real user data file is never touched.
func TestStorageRoundTrip(t *testing.T) {
	t.Setenv("TODO_FILE", filepath.Join(t.TempDir(), "todos.txt"))
	t.Cleanup(func() {
		dataFile = ""
		lists = nil
	})

	if err := initStorage(); err != nil {
		t.Fatalf("initStorage() failed: %v", err)
	}

	want := map[string][]Task{
		"Work": {
			{Text: "Buy milk", Done: false, Priority: High},
			{Text: "Email boss", Done: true, Priority: Medium},
		},
		"Personal": {
			// A pipe and quotes must survive the list|done|priority|text
			// format, which splits on only the first three separators.
			{Text: `Gym | stretch`, Done: false, Priority: Low},
			{Text: `Read "Go" docs`, Done: true, Priority: High},
		},
	}

	lists = want
	if err := saveTasks(); err != nil {
		t.Fatalf("saveTasks() failed: %v", err)
	}

	lists = nil // make sure loadTasks really reads the file back
	loadTasks()

	if !reflect.DeepEqual(lists, want) {
		t.Fatalf("loaded lists = %#v, want %#v", lists, want)
	}
}

// TestMigration checks the one-time import of a legacy todos.txt: the first
// run copies it to the new location, and a later run leaves the existing data
// file alone.
func TestMigration(t *testing.T) {
	base := t.TempDir()
	legacyDir := filepath.Join(base, "legacy")
	if err := os.MkdirAll(legacyDir, 0o755); err != nil {
		t.Fatalf("creating legacy folder: %v", err)
	}

	legacy := "Work|false|high|Buy milk\nPersonal|true|medium|Gym\n"
	if err := os.WriteFile(filepath.Join(legacyDir, "todos.txt"), []byte(legacy), 0o644); err != nil {
		t.Fatalf("writing legacy file: %v", err)
	}

	// The target is in a folder that does not exist yet, so migration must
	// create it just like saveTasks does.
	target := filepath.Join(base, "config", "todo", "todos.txt")
	t.Setenv("TODO_FILE", target)
	chdirTo(t, legacyDir)
	t.Cleanup(func() {
		dataFile = ""
		lists = nil
	})

	if err := initStorage(); err != nil {
		t.Fatalf("first initStorage() failed: %v", err)
	}
	if dataFile != target {
		t.Fatalf("dataFile = %q, want %q", dataFile, target)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("reading imported file: %v", err)
	}
	if string(got) != legacy {
		t.Fatalf("imported content = %q, want %q", got, legacy)
	}

	loadTasks()
	if len(lists) != 2 {
		t.Fatalf("imported %d list(s), want 2: %#v", len(lists), lists)
	}

	// Simulate data the user created after the import. A second run must not
	// copy the old file over it again.
	existing := "Work|true|low|Changed\n"
	if err := os.WriteFile(target, []byte(existing), 0o644); err != nil {
		t.Fatalf("rewriting data file: %v", err)
	}
	if err := initStorage(); err != nil {
		t.Fatalf("second initStorage() failed: %v", err)
	}

	got, err = os.ReadFile(target)
	if err != nil {
		t.Fatalf("reading data file after second run: %v", err)
	}
	if string(got) != existing {
		t.Fatalf("second run overwrote the data file: %q, want %q", got, existing)
	}
}

// chdirTo changes the working directory for a test and restores it afterwards.
// Migration looks for "todos.txt" in the current folder, so the test has to
// run from the folder that holds the legacy file.
func chdirTo(t *testing.T, dir string) {
	t.Helper()

	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("changing to %q: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(old); err != nil {
			t.Errorf("restoring working directory to %q: %v", old, err)
		}
	})
}
