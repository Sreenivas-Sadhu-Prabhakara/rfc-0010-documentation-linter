package lab

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestInteraction(t *testing.T) {
	es := []Event{{"insert", "HELLO"}, {"backspace", ""}, {"backspace", ""}, {"insert", "P"}, {"submit", ""}}
	a, e := Simulate(es, false, 200*time.Millisecond, 0, 20)
	if e != nil {
		t.Fatal(e)
	}
	b, e := Simulate(es, true, 200*time.Millisecond, 0, 20)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(a.Commands, []string{"HELP"}) || !reflect.DeepEqual(a.Commands, b.Commands) || a.Requests != 5 || b.Requests != 1 || a.CompletionMS != 1000 || b.CompletionMS != 200 {
		t.Fatalf("%+v %+v", a, b)
	}
}
func TestEditing(t *testing.T) {
	r, e := Simulate([]Event{{"backspace", ""}, {"insert", "é🙂"}, {"backspace", ""}, {"submit", ""}, {"submit", ""}}, true, 0, 0, 2)
	if e != nil || !reflect.DeepEqual(r.Commands, []string{"é", ""}) {
		t.Fatal(r, e)
	}
	_, e = Simulate([]Event{{"insert", "abc"}}, true, 0, 0, 2)
	if e == nil {
		t.Fatal("limit")
	}
	_, e = Simulate([]Event{{"invalid", ""}}, true, 0, 0, 2)
	if e == nil {
		t.Fatal("kind")
	}
}
func TestConversion(t *testing.T) {
	for _, tc := range []struct {
		s, p              string
		max               int
		out, status, code string
	}{{"abc", "strict", 3, "abc", "unchanged", ""}, {"José", "replacement", 4, "Jos?", "changed", ""}, {"José", "strict", 4, "", "rejected", "unsupported_character"}, {"abc", "strict", 2, "", "rejected", "output_overflow"}, {"", "strict", 0, "", "unchanged", ""}, {"\xff", "replacement", 4, "", "rejected", "invalid_utf8"}, {"\n", "strict", 1, "", "rejected", "unsupported_character"}} {
		r, e := Convert(tc.s, tc.p, tc.max)
		if r.Output != tc.out || r.Status != tc.status {
			t.Fatal(tc, r, e)
		}
		if tc.code != "" {
			var ce *ConversionError
			if !errors.As(e, &ce) || ce.Code != tc.code {
				t.Fatal(e)
			}
		} else if e != nil {
			t.Fatal(e)
		}
	}
}
func TestPipeline(t *testing.T) {
	for _, n := range []int{0, 1, 1006, 1007, 5000} {
		data := bytes.Repeat([]byte{42}, n)
		var got []byte
		seq := 0
		e := Pipeline(context.Background(), data, 1, 2, 1, func(f Frame) error {
			if f.Sequence != seq || len(f.Payload) > MaxPayload {
				t.Fatal(f)
			}
			seq++
			got = append(got, f.Payload...)
			return nil
		})
		if e != nil || !bytes.Equal(data, got) {
			t.Fatal(n, e)
		}
	}
}
func TestPipelineFailures(t *testing.T) {
	sentinel := errors.New("down")
	e := Pipeline(context.Background(), make([]byte, 100000), 1, 2, 1, func(Frame) error { return sentinel })
	if !errors.Is(e, sentinel) {
		t.Fatal(e)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	e = Pipeline(ctx, make([]byte, 100000), 1, 2, 1, func(Frame) error { return nil })
	if !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
func FuzzConvert(f *testing.F) {
	f.Add("José")
	f.Add("abc")
	f.Fuzz(func(t *testing.T, s string) {
		r, e := Convert(s, "replacement", 100)
		if e == nil && len(r.Output) > 100 {
			t.Fatal("overflow")
		}
		r, e = Convert(s, "strict", 100)
		if e == nil && r.Output != s {
			t.Fatal("strict changed input")
		}
	})
}
