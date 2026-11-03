package lab

import (
	"fmt"
	"strings"
	"time"
)

// Note is our local proposal record, not an RFC submission format.
type Note struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Affiliation string `json:"affiliation"`
	Date        string `json:"date"`
	Body        string `json:"body"`
	Replaces    int    `json:"replaces,omitempty"`
}

func ValidateNotes(notes []Note) []string {
	problems := []string{}
	byID := map[int]Note{}
	for i, n := range notes {
		prefix := fmt.Sprintf("record %d", i+1)
		if n.ID < 1 {
			problems = append(problems, prefix+": positive id required")
		}
		if _, ok := byID[n.ID]; ok {
			problems = append(problems, prefix+": duplicate id")
		}
		byID[n.ID] = n
		for _, field := range []struct{ name, value string }{{"title", n.Title}, {"author", n.Author}, {"affiliation", n.Affiliation}, {"body", n.Body}} {
			if strings.TrimSpace(field.value) == "" {
				problems = append(problems, prefix+": "+field.name+" required")
			}
		}
		if _, e := time.Parse("2006-01-02", n.Date); e != nil {
			problems = append(problems, prefix+": date must be YYYY-MM-DD")
		}
	}
	// Only inspect graph structure after identity validation to avoid ambiguity.
	if len(problems) > 0 {
		return problems
	}
	for _, n := range notes {
		if n.Replaces != 0 {
			if _, ok := byID[n.Replaces]; !ok {
				problems = append(problems, fmt.Sprintf("note %d: replacement target %d missing", n.ID, n.Replaces))
			}
		}
	}
	if len(problems) > 0 {
		return problems
	}
	state := map[int]int{}
	var visit func(int) bool
	visit = func(id int) bool {
		if id == 0 {
			return false
		}
		if state[id] == 1 {
			return true
		}
		if state[id] == 2 {
			return false
		}
		state[id] = 1
		if visit(byID[id].Replaces) {
			return true
		}
		state[id] = 2
		return false
	}
	for _, n := range notes {
		if visit(n.ID) {
			problems = append(problems, "replacement graph contains a cycle")
			break
		}
	}
	return problems
}
