package lab

import (
	"reflect"
	"strings"
	"testing"
)

func TestHandshakeStallsAndCompletion(t *testing.T) {
	r, e := RunHandshake([]HandshakeTick{
		{Valid: true, Data: "A"}, {Valid: true, Data: "A"},
		{Ready: true, Valid: true, Data: "A"},
		{Valid: true, Data: "B", Last: true},
		{Valid: true, Ready: true, Data: "B", Last: true},
	}, 4)
	if e != nil || r.Violations != 0 || !r.Finished || r.AcceptedWords != 2 || r.CompletedFrames != 1 {
		t.Fatal(r, e)
	}
	if r.Rows[0].AcceptedWords != 0 || r.Rows[3].CompletedFrames != 0 {
		t.Fatal("stall advanced counters")
	}
	if !reflect.DeepEqual(r.Rows[4].Completed, []string{"A", "B"}) {
		t.Fatal("reordered/duplicated words")
	}
	if r.Rows[0].Held.Data != "A" || len(r.Rows[2].Partial) != 1 {
		t.Fatal("snapshot mutation")
	}
}
func TestHandshakeMutations(t *testing.T) {
	for _, bad := range []HandshakeTick{{}, {Valid: true, Data: "changed"}, {Valid: true, Data: "A", Last: true}} {
		r, _ := RunHandshake([]HandshakeTick{{Valid: true, Data: "A"}, bad, {Ready: true, Valid: true, Data: "A", Last: false}, {Ready: true, Valid: true, Data: "B", Last: true}}, 4)
		if r.Violations != 1 || r.Rows[1].AcceptedWords != 0 || r.Rows[1].Held.Data != "A" || r.CompletedFrames != 1 {
			t.Fatal(r)
		}
	}
}
func TestHandshakeResetAndLimits(t *testing.T) {
	r, _ := RunHandshake([]HandshakeTick{
		{Ready: true, Valid: true, Data: "done", Last: true},
		{Ready: true, Valid: true, Data: "partial"},
		{Valid: true, Data: "held", Last: true},
		{Reset: true},
		{Ready: true, Valid: true, Data: "new", Last: true},
	}, 2)
	if r.Violations != 0 || r.DroppedWords != 1 || r.AcceptedWords != 3 || r.CompletedFrames != 2 || !r.Finished {
		t.Fatal(r)
	}
	if _, e := RunHandshake(nil, 0); e == nil {
		t.Fatal("zero limit")
	}
	if _, e := RunHandshake(make([]HandshakeTick, 1001), 4); e == nil {
		t.Fatal("tick limit")
	}
	r, _ = RunHandshake([]HandshakeTick{{Ready: true, Valid: true, Data: "A"}, {Ready: true, Valid: true, Data: "B"}, {Ready: true, Valid: true, Data: "B", Last: true}}, 2)
	if r.Violations != 1 || r.AcceptedWords != 2 || r.CompletedFrames != 1 {
		t.Fatal("frame limit", r)
	}
	r, _ = RunHandshake([]HandshakeTick{{Ready: true, Valid: true, Data: strings.Repeat("x", 129), Last: true}}, 2)
	if r.Violations != 1 || r.AcceptedWords != 0 {
		t.Fatal("payload limit")
	}
	r, _ = RunHandshake([]HandshakeTick{{Ready: true, Valid: true, Data: strings.Repeat("x", 128), Last: true}}, 1)
	if r.Violations != 0 || !r.Finished {
		t.Fatal("exact bound")
	}
}
func TestHandshakeUnfinishedAndLevelSemantics(t *testing.T) {
	for _, tick := range []HandshakeTick{{Valid: true, Data: "held", Last: true}, {Valid: true, Ready: true, Data: "partial"}} {
		r, _ := RunHandshake([]HandshakeTick{tick}, 4)
		if r.Finished {
			t.Fatal("false completion")
		}
	}
	// In a synchronous ready/valid model, each ready+valid tick transfers a word.
	// Holding both high for two ticks intentionally represents two transfers.
	r, _ := RunHandshake([]HandshakeTick{{Valid: true, Ready: true, Data: "same", Last: true}, {Valid: true, Ready: true, Data: "same", Last: true}}, 4)
	if r.CompletedFrames != 2 || r.AcceptedWords != 2 {
		t.Fatal("level semantics")
	}
}
func FuzzHandshakeAccounting(f *testing.F) {
	f.Add([]byte{3, 1, 0, 4, 7})
	f.Fuzz(func(t *testing.T, b []byte) {
		if len(b) > 1000 {
			b = b[:1000]
		}
		ticks := make([]HandshakeTick, len(b))
		for i, v := range b {
			ticks[i] = HandshakeTick{Ready: v&1 != 0, Valid: v&2 != 0, Last: v&4 != 0, Reset: v&8 != 0, Data: string([]byte{'a' + v%26})}
		}
		r, e := RunHandshake(ticks, 4)
		if e != nil {
			t.Fatal(e)
		}
		completedWords := 0
		for _, row := range r.Rows {
			completedWords += len(row.Completed)
			if row.AcceptedWords != completedWords+row.DroppedWords+len(row.Partial) {
				t.Fatal("lost/duplicated accepted word", row)
			}
			if len(row.Partial) >= 4 {
				t.Fatal("unbounded frame")
			}
		}
	})
}
