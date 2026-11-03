package lab

import "fmt"

// AckEvent belongs to our educational business model, not the RFC 8 wire format.
type AckEvent struct {
	At   int64  `json:"at_ms"`
	Kind string `json:"kind"`
}
type AckResult struct {
	Transport string   `json:"transport"`
	Business  string   `json:"business"`
	Caller    string   `json:"caller"`
	Timeline  []string `json:"timeline"`
}

// Acknowledge processes ordered evidence, including evidence after caller timeout.
// Timeout does not decide business outcome. Repeated evidence is idempotent.
func Acknowledge(events []AckEvent) (AckResult, error) {
	r := AckResult{Transport: "pending", Business: "unknown", Caller: "waiting", Timeline: []string{}}
	var previous int64
	for i, e := range events {
		if e.At < 0 || (i > 0 && e.At < previous) {
			return r, fmt.Errorf("event %d: timestamps must be nonnegative and ordered", i)
		}
		previous = e.At
		switch e.Kind {
		case "received":
			r.Transport = "received"
		case "committed":
			if r.Business == "rejected" {
				return r, fmt.Errorf("event %d: conflicting terminal evidence", i)
			}
			r.Business = "committed"
			if r.Caller == "waiting" {
				r.Caller = "success"
			}
		case "rejected":
			if r.Business == "committed" {
				return r, fmt.Errorf("event %d: conflicting terminal evidence", i)
			}
			r.Business = "rejected"
			if r.Caller == "waiting" {
				r.Caller = "rejected"
			}
		case "timeout":
			if r.Caller == "waiting" {
				r.Caller = "timed_out"
			}
		default:
			return r, fmt.Errorf("event %d: unsupported event %q", i, e.Kind)
		}
		r.Timeline = append(r.Timeline, fmt.Sprintf("%dms %s: transport=%s business=%s caller=%s", e.At, e.Kind, r.Transport, r.Business, r.Caller))
	}
	return r, nil
}
