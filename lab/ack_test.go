package lab

import "testing"

func TestAckEvidence(t *testing.T) {
	cases := []struct {
		name             string
		events           []AckEvent
		business, caller string
	}{
		{"receipt only", []AckEvent{{1, "received"}}, "unknown", "waiting"},
		{"timeout", []AckEvent{{1, "received"}, {2, "timeout"}}, "unknown", "timed_out"},
		{"late commit", []AckEvent{{1, "timeout"}, {2, "committed"}}, "committed", "timed_out"},
		{"early commit", []AckEvent{{1, "committed"}, {2, "timeout"}}, "committed", "success"},
		{"reject", []AckEvent{{1, "rejected"}}, "rejected", "rejected"},
		{"duplicate", []AckEvent{{1, "committed"}, {1, "committed"}}, "committed", "success"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, e := Acknowledge(c.events)
			if e != nil || r.Business != c.business || r.Caller != c.caller {
				t.Fatal(r, e)
			}
		})
	}
}
func TestAckInvalid(t *testing.T) {
	for _, es := range [][]AckEvent{{{-1, "received"}}, {{2, "received"}, {1, "timeout"}}, {{0, "invalid"}}, {{1, "committed"}, {2, "rejected"}}, {{1, "rejected"}, {2, "committed"}}} {
		if _, e := Acknowledge(es); e == nil {
			t.Fatal("accepted invalid evidence", es)
		}
	}
}
