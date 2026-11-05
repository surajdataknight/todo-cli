package main

import "strings"

// Priority is how important a task is.
// iota gives each constant a unique number, starting at 0:
//
//	Low  = 0
//	Medium = 1
//	High = 2
type Priority int

const (
	Low Priority = iota
	Medium
	High
)

// String turns a Priority into a readable word (e.g. High -> "high").
func (p Priority) String() string {
	switch p {
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	default:
		return "unknown"
	}
}

// parsePriority turns a word like "high" into a Priority.
// The second return value tells whether the word was recognized.
func parsePriority(s string) (Priority, bool) {
	switch strings.ToLower(s) {
	case "low":
		return Low, true
	case "medium":
		return Medium, true
	case "high":
		return High, true
	default:
		return Medium, false
	}
}

// Task is a single to-do item.
type Task struct {
	Text     string   // what the task says
	Done     bool     // true if finished, false if not
	Priority Priority // how important the task is (low/medium/high)
}
