package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIExitAndOutput(t *testing.T) {
	for _, tc := range []struct {
		input string
		code  int
	}{
		{`[{"ready":true,"valid":true,"data":"A","last":true}]`, 0},
		{`[{"valid":true,"data":"A"}]`, 2},
		{`[{"valid":true,"data":"A"},{"valid":true,"data":"B"}]`, 2},
		{`[] []`, 1}, {`[{"unknown":true}]`, 1}, {`null`, 1},
	} {
		p := filepath.Join(t.TempDir(), "ticks.json")
		os.WriteFile(p, []byte(tc.input), 0600)
		var out, errs bytes.Buffer
		if code := run([]string{"-file", p}, &out, &errs); code != tc.code {
			t.Fatal(tc, code, errs.String())
		}
		if tc.code == 1 && out.Len() != 0 {
			t.Fatal("malformed input emitted output")
		}
		if tc.code != 1 && !strings.Contains(out.String(), `"finished"`) {
			t.Fatal("missing report")
		}
	}
}
func TestHTMLIsEscaped(t *testing.T) {
	r, e := decodeRun([]byte(`[{"ready":true,"valid":true,"data":"</pre><script>alert(1)</script>","last":true}]`), 4)
	if e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e = renderHTML(&out, r); e != nil {
		t.Fatal(e)
	}
	if strings.Contains(out.String(), "</pre><script>alert(1)</script>") {
		t.Fatal("HTML injection")
	}
	if !strings.Contains(out.String(), "Accepted words") {
		t.Fatal("missing report")
	}
}
func TestDecoderBounds(t *testing.T) {
	if _, e := decodeRun(bytes.Repeat([]byte{' '}, (1<<20)+1), 4); e == nil {
		t.Fatal("size bound")
	}
	if _, e := decodeRun([]byte(`[]`), 65); e == nil {
		t.Fatal("word limit")
	}
}
