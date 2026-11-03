package main

import (
	"encoding/json"
	"example.com/rfc-go-labs/lab"
	"fmt"
	"os"
)

func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func run() error {
	m, e := lab.NewLinks(4)
	if e != nil {
		return e
	}
	p, e := m.Open("alice", "bank-sim", "primary", 0)
	if e != nil {
		return e
	}
	if e = m.Confirm("alice", p.ID); e != nil {
		return e
	}
	a, e := m.Open("alice", "bank-sim", "auxiliary", p.ID)
	if e != nil {
		return e
	}
	if e = m.Confirm("alice", a.ID); e != nil {
		return e
	}
	before := m.Snapshot()
	if e = m.Close("alice", p.ID); e != nil {
		return e
	}
	return json.NewEncoder(os.Stdout).Encode(map[string]any{"before_close": before, "after_close": m.Snapshot()})
}
