package lab

import (
	"strings"
	"sync"
	"testing"
)

func TestSessionsHandshakeAndRouting(t *testing.T) {
	m, _ := NewSessions(8, 1)
	p, _ := m.Open("bank-a", "cloud-a", "primary", 0)
	if _, e := m.Open("bank-a", "cloud-a", "primary", 0); e == nil {
		t.Fatal("duplicate pending primary")
	}
	if _, e := m.Open("bank-a", "cloud-a", "auxiliary", p.ID); e == nil {
		t.Fatal("aux before ACK")
	}
	if e := m.Deliver("cloud-a", 15, p.ID, "early"); e == nil {
		t.Fatal("data before ACK")
	}
	for _, f := range []func() error{
		func() error { return m.Ack("intruder", p.ID, "cloud-a", 15) },
		func() error { return m.Ack("bank-a", p.ID, "wrong-host", 15) },
		func() error { return m.Ack("bank-a", p.ID, "cloud-a", 0) },
	} {
		if f() == nil {
			t.Fatal("invalid ACK")
		}
	}
	if e := m.Ack("bank-a", p.ID, "cloud-a", 15); e != nil {
		t.Fatal(e)
	}
	if e := m.Ack("bank-a", p.ID, "cloud-a", 15); e != nil {
		t.Fatal("replay", e)
	}
	if e := m.Ack("bank-a", p.ID, "cloud-a", 16); e == nil {
		t.Fatal("conflict")
	}
	a, _ := m.Open("bank-a", "cloud-a", "auxiliary", p.ID)
	if e := m.Ack("bank-a", a.ID, "cloud-a", 15); e == nil {
		t.Fatal("route collision")
	}
	if e := m.Ack("bank-a", a.ID, "cloud-a", 21); e != nil {
		t.Fatal(e)
	}
	q, _ := m.Open("bank-b", "cloud-b", "primary", 0)
	if e := m.Ack("bank-b", q.ID, "cloud-b", 15); e != nil {
		t.Fatal("host namespace", e)
	}
	if e := m.Deliver("cloud-a", 15, p.ID, "status"); e != nil {
		t.Fatal(e)
	}
	if e := m.Deliver("cloud-a", 15, p.ID, "overload"); e == nil {
		t.Fatal("unbounded inbox")
	}
	if e := m.Deliver("cloud-a", 21, a.ID, strings.Repeat("x", 1025)); e == nil {
		t.Fatal("oversize")
	}
	if e := m.Deliver("cloud-a", 21, a.ID, strings.Repeat("x", 1024)); e != nil {
		t.Fatal(e)
	}
	if _, e := m.Drain("bank-b", p.ID); e == nil {
		t.Fatal("cross-user read")
	}
	got, e := m.Drain("bank-a", p.ID)
	if e != nil || len(got) != 1 || got[0] != "status" {
		t.Fatal(got, e)
	}
	other, _ := m.Drain("bank-b", q.ID)
	if len(other) != 0 {
		t.Fatal("cross-host leak")
	}
	if e := m.Close("intruder", p.ID); e == nil {
		t.Fatal("cross-user close")
	}
	if e := m.Close("bank-a", p.ID); e != nil {
		t.Fatal(e)
	}
	if len(m.Snapshot()) != 1 {
		t.Fatal("cascade")
	}
	if e := m.Deliver("cloud-a", 21, a.ID, "late"); e == nil {
		t.Fatal("closed route")
	}
	r, _ := m.Open("bank-a", "cloud-a", "primary", 0)
	if r.Number != p.Number || r.ID == p.ID {
		t.Fatal("generation/reuse")
	}
	if e := m.Ack("bank-a", p.ID, "cloud-a", 15); e == nil {
		t.Fatal("stale ACK")
	}
	m.Ack("bank-a", r.ID, "cloud-a", 15)
	if e := m.Deliver("cloud-a", 15, p.ID, "old data"); e == nil {
		t.Fatal("stale data after reuse")
	}
}
func TestSessionsRejectAndCapacity(t *testing.T) {
	if _, e := NewSessions(0, 1); e == nil {
		t.Fatal("capacity")
	}
	if _, e := NewSessions(1, 0); e == nil {
		t.Fatal("depth")
	}
	m, _ := NewSessions(2, 1)
	p, _ := m.Open("u", "h", "primary", 0)
	if e := m.Reject("other", p.ID); e == nil {
		t.Fatal("ownership")
	}
	if e := m.Reject("u", p.ID); e != nil {
		t.Fatal(e)
	}
	q, _ := m.Open("u", "h", "primary", 0)
	if q.Number != p.Number {
		t.Fatal("reservation leaked")
	}
	m.Ack("u", q.ID, "h", 12)
	if e := m.Reject("u", q.ID); e == nil {
		t.Fatal("reject established")
	}
	a, _ := m.Open("u", "h", "auxiliary", q.ID)
	if _, e := m.Open("v", "z", "primary", 0); e == nil {
		t.Fatal("capacity")
	}
	m.Close("u", q.ID)
	if e := m.Ack("u", a.ID, "h", 13); e == nil {
		t.Fatal("pending child survived close")
	}
	if len(m.Snapshot()) != 0 {
		t.Fatal("pending cleanup")
	}
}
func TestSessionsConcurrentDuplicatePrimary(t *testing.T) {
	m, _ := NewSessions(64, 1)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); m.Open("u", "h", "primary", 0) }()
	}
	wg.Wait()
	if len(m.Snapshot()) != 1 {
		t.Fatal("primary uniqueness race")
	}
}
