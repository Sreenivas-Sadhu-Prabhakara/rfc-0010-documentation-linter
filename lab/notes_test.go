package lab

import (
	"strings"
	"testing"
)

func goodNote(id int) Note {
	return Note{ID: id, Title: "Repeated titles allowed", Author: "Engineer", Affiliation: "Lab", Date: "2026-11-03", Body: "Could we validate proposals before publishing?"}
}
func TestNotesValid(t *testing.T) {
	a, b := goodNote(1), goodNote(2)
	b.Replaces = 1
	if p := ValidateNotes([]Note{a, b}); len(p) != 0 {
		t.Fatal(p)
	}
}
func TestNotesInvalid(t *testing.T) {
	a := goodNote(1)
	cases := []struct {
		name  string
		notes []Note
		want  string
	}{{"duplicate", []Note{a, a}, "duplicate"}, {"missing target", []Note{func() Note { n := a; n.Replaces = 2; return n }()}, "missing"}, {"self cycle", []Note{func() Note { n := a; n.Replaces = 1; return n }()}, "cycle"}, {"two cycle", []Note{func() Note { n := a; n.Replaces = 2; return n }(), func() Note { n := goodNote(2); n.Replaces = 1; return n }()}, "cycle"}, {"bad date", []Note{func() Note { n := a; n.Date = "2026-02-30"; return n }()}, "date"}, {"blank", []Note{func() Note { n := a; n.Author = " "; return n }()}, "author"}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := ValidateNotes(c.notes)
			if !strings.Contains(strings.Join(p, ";"), c.want) {
				t.Fatal(p)
			}
		})
	}
}
