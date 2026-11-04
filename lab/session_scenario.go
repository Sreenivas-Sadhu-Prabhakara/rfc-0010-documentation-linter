package lab

import "fmt"

type SessionEvent struct {
	Op       string `json:"op"`
	User     string `json:"user,omitempty"`
	Host     string `json:"host,omitempty"`
	Kind     string `json:"kind,omitempty"`
	ID       uint64 `json:"id,omitempty"`
	Parent   uint64 `json:"parent,omitempty"`
	Incoming int    `json:"incoming,omitempty"`
	Payload  string `json:"payload,omitempty"`
}
type SessionResult struct {
	Step     int           `json:"step"`
	Event    SessionEvent  `json:"event"`
	Error    string        `json:"error,omitempty"`
	Opened   *Lease        `json:"opened,omitempty"`
	Messages []string      `json:"messages,omitempty"`
	State    []SessionView `json:"state"`
}

// RunSessions uses a fresh manager per scenario. Semantic failures are recorded
// without stopping the trace, so cleanup after a failed action can be inspected.
func RunSessions(events []SessionEvent, capacity, depth int) ([]SessionResult, error) {
	if len(events) > 1000 {
		return nil, fmt.Errorf("maximum 1000 events")
	}
	m, err := NewSessions(capacity, depth)
	if err != nil {
		return nil, err
	}
	out := make([]SessionResult, 0, len(events))
	for i, e := range events {
		r := SessionResult{Step: i + 1, Event: e}
		var err error
		switch e.Op {
		case "open":
			l, openErr := m.Open(e.User, e.Host, e.Kind, e.Parent)
			err = openErr
			if err == nil {
				r.Opened = &l
			}
		case "ack":
			err = m.Ack(e.User, e.ID, e.Host, e.Incoming)
		case "reject":
			err = m.Reject(e.User, e.ID)
		case "close":
			err = m.Close(e.User, e.ID)
		case "deliver":
			err = m.Deliver(e.Host, e.Incoming, e.ID, e.Payload)
		case "drain":
			r.Messages, err = m.Drain(e.User, e.ID)
		default:
			err = fmt.Errorf("unknown operation %q", e.Op)
		}
		if err != nil {
			r.Error = err.Error()
		}
		r.State = m.Snapshot()
		out = append(out, r)
	}
	return out, nil
}
