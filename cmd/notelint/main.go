package main

import (
	"bytes"
	"encoding/json"
	"example.com/rfc-go-labs/lab"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func run() error {
	file := flag.String("file", "examples/notes.json", "JSON proposal registry")
	flag.Parse()
	f, e := os.Open(*file)
	if e != nil {
		return e
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, 1048577))
	if e != nil {
		return e
	}
	if len(b) > 1048576 {
		return fmt.Errorf("registry exceeds 1 MiB")
	}
	var notes []lab.Note
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if e = d.Decode(&notes); e != nil {
		return e
	}
	if notes == nil {
		return fmt.Errorf("registry must be an array, not null")
	}
	var extra any
	if e = d.Decode(&extra); e != io.EOF {
		return fmt.Errorf("expected one JSON array")
	}
	if len(notes) > 10000 {
		return fmt.Errorf("registry exceeds 10000 notes")
	}
	problems := lab.ValidateNotes(notes)
	if e = json.NewEncoder(os.Stdout).Encode(map[string]any{"valid": len(problems) == 0, "notes": len(notes), "problems": problems}); e != nil {
		return e
	}
	if len(problems) > 0 {
		return fmt.Errorf("registry validation failed")
	}
	return nil
}
