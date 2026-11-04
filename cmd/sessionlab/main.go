package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"

	"example.com/rfc-go-labs/lab"
)

func main() {
	file := flag.String("file", "examples/sessions.json", "JSON scenario file")
	capacity := flag.Int("capacity", 8, "maximum pending+established connections (1..64)")
	depth := flag.Int("depth", 1, "messages per inbox (1..64)")
	html := flag.Bool("html", false, "emit a standalone static trace viewer instead of JSON")
	flag.Parse()
	if flag.NArg() != 0 {
		fail(fmt.Errorf("unexpected positional arguments"))
	}
	f, err := os.Open(*file)
	if err != nil {
		fail(err)
	}
	defer f.Close()
	input, err := io.ReadAll(io.LimitReader(f, (1<<20)+1))
	if err != nil {
		fail(err)
	}
	rows, err := decodeRun(input, *capacity, *depth)
	if err != nil {
		fail(err)
	}
	if *html {
		t := template.Must(template.New("trace").Funcs(template.FuncMap{"pretty": func(v any) string { b, _ := json.MarshalIndent(v, "", "  "); return string(b) }}).Parse(page))
		if err = t.Execute(os.Stdout, rows); err != nil {
			fail(err)
		}
	} else {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err = enc.Encode(rows); err != nil {
			fail(err)
		}
	}
}
func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
func decodeRun(input []byte, capacity, depth int) ([]lab.SessionResult, error) {
	if len(input) > 1<<20 {
		return nil, fmt.Errorf("input exceeds 1 MiB")
	}
	var events []lab.SessionEvent
	d := json.NewDecoder(bytes.NewReader(input))
	d.DisallowUnknownFields()
	if err := d.Decode(&events); err != nil {
		return nil, err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("expected one JSON array")
	}
	return lab.RunSessions(events, capacity, depth)
}

const page = `<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>RFC 11 · Session lifecycle lab</title>
<style>body{font:17px/1.6 system-ui;background:#101827;color:#e8edf5;max-width:1000px;margin:40px auto;padding:0 24px}h1{font-size:clamp(30px,5vw,52px);line-height:1.1}a{color:#86dafa}p{max-width:75ch}details{border:1px solid #43516a;border-radius:12px;margin:14px 0;padding:16px;background:#172238}summary{cursor:pointer;font-weight:650}.error{color:#ffb9ae}.ok{color:#99e3b1}pre{white-space:pre-wrap;overflow-wrap:anywhere;font:14px/1.5 ui-monospace,monospace}.tag{color:#a4b5cf}button{font:inherit;background:#86dafa;border:0;border-radius:8px;padding:10px 18px;cursor:pointer}footer{margin:40px 0;color:#a4b5cf}</style>
<p class="tag">INTERNET RFC SERIES / 011</p><h1>A connection has a lifecycle.</h1><p>Explore a Go-generated trace: reserve → acknowledge → route → drain → close. Expand any step to inspect all active connections and their bounded inboxes.</p>
<p>Inspired by G. Deloche's August 1969 GORDO design. <a href="https://www.rfc-editor.org/info/rfc11/">RFC 11</a> has status Unknown and is obsoleted by RFC 33. This is an educational model with modern generation IDs and explicit backpressure, not historical protocol compliance or a banking service.</p>
<p>This static viewer contains saved output from Go. It does not run a server or execute Go in your browser. To change the scenario, edit examples/sessions.json and regenerate with the command in the README.</p><button id="expand">Expand all steps</button>
{{range .}}<details><summary>Step {{.Step}} · {{.Event.Op}} · {{if .Error}}<span class="error">Rejected</span>{{else}}<span class="ok">Accepted</span>{{end}}</summary><pre>{{pretty .}}</pre></details>{{end}}
<footer>MIT-licensed source and tests · RFC covered: 11 · Next: 12 · No real transfers or credentials</footer><script>document.getElementById('expand').addEventListener('click',()=>{const all=[...document.querySelectorAll('details')];const open=all.some(x=>!x.open);all.forEach(x=>x.open=open);document.getElementById('expand').textContent=open?'Collapse all steps':'Expand all steps';});</script></html>`
