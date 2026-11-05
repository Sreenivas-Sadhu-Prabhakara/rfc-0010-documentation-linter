package lab

import "fmt"

// HandshakeTick is a modern word-level ready/valid teaching abstraction.
// It is not the asynchronous bit-level IMP/HOST electrical interface.
type HandshakeTick struct {
	Ready bool   `json:"ready"`
	Valid bool   `json:"valid"`
	Data  string `json:"data,omitempty"`
	Last  bool   `json:"last,omitempty"`
	Reset bool   `json:"reset,omitempty"`
}
type HeldWord struct {
	Data string `json:"data"`
	Last bool   `json:"last"`
}
type HandshakeRow struct {
	Tick            int           `json:"tick"`
	Signals         HandshakeTick `json:"signals"`
	Action          string        `json:"action"`
	Error           string        `json:"error,omitempty"`
	Held            *HeldWord     `json:"held,omitempty"`
	Partial         []string      `json:"partial"`
	Completed       []string      `json:"completed,omitempty"`
	AcceptedWords   int           `json:"accepted_words"`
	CompletedFrames int           `json:"completed_frames"`
	DroppedWords    int           `json:"dropped_words"`
}
type HandshakeReport struct {
	Rows            []HandshakeRow `json:"rows"`
	Violations      int            `json:"violations"`
	Finished        bool           `json:"finished"`
	AcceptedWords   int            `json:"accepted_words"`
	CompletedFrames int            `json:"completed_frames"`
	DroppedWords    int            `json:"dropped_words"`
}

// RunHandshake records invalid ticks without mutating receiver state. A valid
// recovery or reset can follow, but the trace retains its violation count.
// Reset discards a partial frame and an unaccepted held offer. Completed frame
// counts survive reset. Held offers are not included in accepted/dropped words.
func RunHandshake(ticks []HandshakeTick, maxWords int) (HandshakeReport, error) {
	if maxWords < 1 || maxWords > 64 {
		return HandshakeReport{}, fmt.Errorf("max words must be 1..64")
	}
	if len(ticks) > 1000 {
		return HandshakeReport{}, fmt.Errorf("maximum 1000 ticks")
	}
	r := HandshakeReport{Rows: make([]HandshakeRow, 0, len(ticks))}
	partial := []string{}
	var held *HeldWord
	for i, t := range ticks {
		row := HandshakeRow{Tick: i + 1, Signals: t, Action: "idle"}
		switch {
		case t.Reset:
			r.DroppedWords += len(partial)
			partial = []string{}
			held = nil
			row.Action = "reset"
		case len(t.Data) > 128:
			row.Error = "data exceeds 128 bytes"
		case held != nil && (!t.Valid || t.Data != held.Data || t.Last != held.Last):
			row.Error = "producer changed or withdrew a stalled offer"
		case !t.Valid:
			// Data and Last have no meaning while Valid is false.
		case !t.Ready:
			held = &HeldWord{t.Data, t.Last}
			row.Action = "stall"
		case len(partial)+1 > maxWords || (len(partial)+1 == maxWords && !t.Last):
			row.Error = "frame requires Last by the configured word limit"
		default:
			r.AcceptedWords++
			held = nil
			partial = append(partial, t.Data)
			row.Action = "transfer"
			if t.Last {
				r.CompletedFrames++
				row.Completed = append([]string{}, partial...)
				partial = []string{}
				row.Action = "complete"
			}
		}
		if row.Error != "" {
			row.Action = "violation"
			r.Violations++
		}
		row.Partial = append([]string{}, partial...)
		if held != nil {
			h := *held
			row.Held = &h
		}
		row.AcceptedWords = r.AcceptedWords
		row.CompletedFrames = r.CompletedFrames
		row.DroppedWords = r.DroppedWords
		r.Rows = append(r.Rows, row)
	}
	r.Finished = held == nil && len(partial) == 0
	return r, nil
}
