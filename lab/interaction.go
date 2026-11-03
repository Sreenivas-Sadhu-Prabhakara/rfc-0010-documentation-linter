package lab

import (
	"fmt"
	"time"
	"unicode/utf8"
)

type Event struct {
	Kind string `json:"kind"`
	Text string `json:"text,omitempty"`
}
type Interaction struct {
	Mode           string   `json:"mode"`
	Commands       []string `json:"commands"`
	Requests       int      `json:"requests"`
	DisplayDelayMS int64    `json:"display_delay_ms"`
	CompletionMS   int64    `json:"completion_ms"`
}

// Simulate uses a stop-and-wait remote model; input events are queued at t=0.
// Local editing costs zero simulated time. One round trip is paid per submission.
func Simulate(events []Event, local bool, rtt, process time.Duration, limit int) (Interaction, error) {
	r := Interaction{Mode: "remote", Commands: []string{}}
	if local {
		r.Mode = "local"
	}
	if rtt < 0 || process < 0 || limit < 1 {
		return r, fmt.Errorf("invalid timing or input limit")
	}
	var text []rune
	for _, e := range events {
		switch e.Kind {
		case "insert":
			if !utf8.ValidString(e.Text) {
				return r, fmt.Errorf("invalid UTF-8")
			}
			if len(text)+utf8.RuneCountInString(e.Text) > limit {
				return r, fmt.Errorf("input limit exceeded")
			}
			text = append(text, []rune(e.Text)...)
		case "backspace":
			if len(text) > 0 {
				text = text[:len(text)-1]
			}
		case "submit":
			r.Commands = append(r.Commands, string(text))
			text = nil
		default:
			return r, fmt.Errorf("unsupported event %q", e.Kind)
		}
		if !local || e.Kind == "submit" {
			r.Requests++
			r.CompletionMS += (rtt + process).Milliseconds()
		}
		if !local && e.Kind != "submit" {
			r.DisplayDelayMS = r.CompletionMS
		}
	}
	return r, nil
}
