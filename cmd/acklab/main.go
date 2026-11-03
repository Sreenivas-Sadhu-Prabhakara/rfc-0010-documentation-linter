package main

import (
	"encoding/json"
	"example.com/rfc-go-labs/lab"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	file := flag.String("events", "", "JSON event file (default: received, timeout, late commit)")
	flag.Parse()
	events := []lab.AckEvent{{At: 10, Kind: "received"}, {At: 100, Kind: "timeout"}, {At: 200, Kind: "committed"}}
	if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			return err
		}
		defer f.Close()
		d := json.NewDecoder(io.LimitReader(f, 1048577))
		d.DisallowUnknownFields()
		if err = d.Decode(&events); err != nil {
			return err
		}
		var extra any
		if err = d.Decode(&extra); err != io.EOF {
			return fmt.Errorf("expected exactly one JSON array")
		}
		if len(events) > 10000 {
			return fmt.Errorf("too many events")
		}
	}
	r, err := lab.Acknowledge(events)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
