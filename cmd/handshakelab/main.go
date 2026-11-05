package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"

	"example.com/rfc-go-labs/lab"
)

//go:embed page.html
var page string

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
func run(args []string, out, errs io.Writer) int {
	f := flag.NewFlagSet("handshakelab", flag.ContinueOnError)
	f.SetOutput(errs)
	file := f.String("file", "examples/handshake.json", "JSON tick array")
	maxWords := f.Int("max-words", 4, "word limit per frame (1..64)")
	html := f.Bool("html", false, "emit a static HTML report")
	if e := f.Parse(args); e != nil {
		if e == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if f.NArg() != 0 {
		fmt.Fprintln(errs, "unexpected positional arguments")
		return 1
	}
	input, e := os.Open(*file)
	if e != nil {
		fmt.Fprintln(errs, e)
		return 1
	}
	defer input.Close()
	b, e := io.ReadAll(io.LimitReader(input, (1<<20)+1))
	if e != nil {
		fmt.Fprintln(errs, e)
		return 1
	}
	r, e := decodeRun(b, *maxWords)
	if e != nil {
		fmt.Fprintln(errs, e)
		return 1
	}
	if *html {
		e = renderHTML(out, r)
	} else {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		e = enc.Encode(r)
	}
	if e != nil {
		fmt.Fprintln(errs, e)
		return 1
	}
	if r.Violations > 0 || !r.Finished {
		return 2
	}
	return 0
}
func decodeRun(b []byte, maxWords int) (lab.HandshakeReport, error) {
	if len(b) > 1<<20 {
		return lab.HandshakeReport{}, fmt.Errorf("input exceeds 1 MiB")
	}
	var ticks []lab.HandshakeTick
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if e := d.Decode(&ticks); e != nil {
		return lab.HandshakeReport{}, e
	}
	if ticks == nil {
		return lab.HandshakeReport{}, fmt.Errorf("expected JSON tick array, not null")
	}
	var extra any
	if e := d.Decode(&extra); e != io.EOF {
		return lab.HandshakeReport{}, fmt.Errorf("expected one JSON array")
	}
	return lab.RunHandshake(ticks, maxWords)
}
func renderHTML(w io.Writer, r lab.HandshakeReport) error {
	t, e := template.New("trace").Funcs(template.FuncMap{"pretty": func(v any) string { b, _ := json.MarshalIndent(v, "", "  "); return string(b) }}).Parse(page)
	if e != nil {
		return e
	}
	return t.Execute(w, r)
}
