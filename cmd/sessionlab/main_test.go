package main

import (
	"strings"
	"testing"
)

func TestDecodeRunRejectsMalformedInput(t *testing.T) {
	for _, input := range []string{`{}`, `[{"oops":1}]`, `[] []`, `[`, strings.Repeat(" ", (1<<20)+1)} {
		if _, e := decodeRun([]byte(input), 8, 1); e == nil {
			t.Fatal("accepted malformed input")
		}
	}
}
func TestTraceContinuesAfterRejectedAction(t *testing.T) {
	rows, e := decodeRun([]byte(`[{"op":"open","user":"u","host":"h","kind":"primary"},{"op":"deliver","host":"h","incoming":15,"id":1,"payload":"too early"},{"op":"ack","user":"u","host":"h","id":1,"incoming":15},{"op":"close","user":"u","id":1}]`), 8, 1)
	if e != nil || len(rows) != 4 {
		t.Fatal(rows, e)
	}
	if rows[1].Error == "" || rows[2].Error != "" || len(rows[3].State) != 0 {
		t.Fatal(rows)
	}
	if rows[0].State[0].Established {
		t.Fatal("snapshot mutated by later steps")
	}
}
