package lab

import (
	"sync"
	"testing"
)

func TestLinksOwnershipAndReuse(t *testing.T) {
	m, _ := NewLinks(4)
	p, _ := m.Open("a", "h", "primary", 0)
	if _, e := m.Open("a", "h", "auxiliary", p.ID); e == nil {
		t.Fatal("unconfirmed primary")
	}
	if e := m.Confirm("b", p.ID); e == nil {
		t.Fatal("cross user confirm")
	}
	m.Confirm("a", p.ID)
	a, e := m.Open("a", "h", "auxiliary", p.ID)
	if e != nil || a.Number == 0 {
		t.Fatal(a, e)
	}
	if _, e = m.Open("b", "h", "auxiliary", p.ID); e == nil {
		t.Fatal("cross user")
	}
	if e = m.Close("b", p.ID); e == nil {
		t.Fatal("cross user close")
	}
	if e = m.Close("a", p.ID); e != nil || len(m.Snapshot()) != 0 {
		t.Fatal(e)
	}
	q, _ := m.Open("a", "h", "primary", 0)
	if q.Number != p.Number || q.ID == p.ID {
		t.Fatal("bad reuse")
	}
	if e = m.Confirm("a", p.ID); e == nil {
		t.Fatal("stale confirmation")
	}
}
func TestLinksCapacityConcurrent(t *testing.T) {
	m, _ := NewLinks(8)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); m.Open("a", "h", "primary", 0) }()
	}
	wg.Wait()
	ls := m.Snapshot()
	if len(ls) != 8 {
		t.Fatal(len(ls))
	}
	seen := map[int]bool{}
	for _, l := range ls {
		if seen[l.Number] || l.Number == 0 {
			t.Fatal(l)
		}
		seen[l.Number] = true
	}
}
func TestLinksInvalid(t *testing.T) {
	if _, e := NewLinks(65); e == nil {
		t.Fatal("capacity")
	}
	m, _ := NewLinks(1)
	if _, e := m.Open("", "h", "primary", 0); e == nil {
		t.Fatal("empty")
	}
	if _, e := m.Open("a", "h", "bad", 0); e == nil {
		t.Fatal("kind")
	}
	m.Open("a", "h", "primary", 0)
	if _, e := m.Open("b", "h", "primary", 0); e == nil {
		t.Fatal("full")
	}
}
