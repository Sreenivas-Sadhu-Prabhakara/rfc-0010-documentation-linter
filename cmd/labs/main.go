package main

import (
	"context"
	"encoding/json"
	"example.com/rfc-go-labs/lab"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func output(v any) error { return json.NewEncoder(os.Stdout).Encode(v) }
func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: labs interaction|convert|pipeline [flags]")
	}
	f := flag.NewFlagSet(os.Args[1], flag.ContinueOnError)
	switch os.Args[1] {
	case "interaction":
		rtt := f.Duration("rtt", 200*time.Millisecond, "simulated round trip")
		delay := f.Duration("processing", 10*time.Millisecond, "server work per request")
		path := f.String("events", "", "JSON event array (default HELP demo)")
		if err := f.Parse(os.Args[2:]); err != nil {
			return err
		}
		es := []lab.Event{}
		for _, c := range "HELLO" {
			es = append(es, lab.Event{Kind: "insert", Text: string(c)})
		}
		es = append(es, lab.Event{Kind: "backspace"}, lab.Event{Kind: "backspace"}, lab.Event{Kind: "insert", Text: "P"}, lab.Event{Kind: "submit"})
		if *path != "" {
			b, e := os.ReadFile(*path)
			if e != nil {
				return e
			}
			if e = json.Unmarshal(b, &es); e != nil {
				return e
			}
		}
		results := []lab.Interaction{}
		for _, local := range []bool{false, true} {
			r, e := lab.Simulate(es, local, *rtt, *delay, 4096)
			if e != nil {
				return e
			}
			results = append(results, r)
		}
		return output(results)
	case "convert":
		text := f.String("text", "José", "UTF-8 input")
		policy := f.String("policy", "strict", "strict or replacement")
		max := f.Int("max", 1006, "maximum output bytes")
		if err := f.Parse(os.Args[2:]); err != nil {
			return err
		}
		r, e := lab.Convert(*text, *policy, *max)
		if e != nil {
			_ = output(r)
			return e
		}
		return output(r)
	case "pipeline":
		text := f.String("text", "hello from RFC 7", "payload")
		path := f.String("file", "", "optional binary input")
		capacity := f.Int("capacity", 2, "queue capacity")
		link := f.Int("link", 1, "logical link")
		dest := f.Int("destination", 2, "host ID")
		fail := f.Int("fail-at", -1, "inject handler failure at sequence")
		if e := f.Parse(os.Args[2:]); e != nil {
			return e
		}
		data := []byte(*text)
		if *path != "" {
			b, e := os.ReadFile(*path)
			if e != nil {
				return e
			}
			data = b
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return lab.Pipeline(ctx, data, *link, *dest, *capacity, func(frame lab.Frame) error {
			if frame.Sequence == *fail {
				return fmt.Errorf("injected handler failure at %d", *fail)
			}
			return output(frame)
		})
	default:
		return fmt.Errorf("unknown command %q", os.Args[1])
	}
}
