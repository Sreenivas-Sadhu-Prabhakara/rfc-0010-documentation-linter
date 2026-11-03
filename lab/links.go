package lab

import (
	"fmt"
	"sort"
	"sync"
)

// Lease is a modern opaque handle, not a historical wire identifier.
type Lease struct {
	ID          uint64 `json:"id"`
	User        string `json:"user"`
	Host        string `json:"host"`
	Number      int    `json:"number"`
	Kind        string `json:"kind"`
	Parent      uint64 `json:"parent,omitempty"`
	Established bool   `json:"established"`
}
type Links struct {
	mu       sync.Mutex
	capacity int
	next     uint64
	active   map[uint64]Lease
}

func NewLinks(capacity int) (*Links, error) {
	if capacity < 1 || capacity > 64 {
		return nil, fmt.Errorf("capacity must be 1..64")
	}
	return &Links{capacity: capacity, active: make(map[uint64]Lease)}, nil
}
func (m *Links) Open(user, host, kind string, parent uint64) (Lease, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if user == "" || host == "" {
		return Lease{}, fmt.Errorf("user and host required")
	}
	if kind != "primary" && kind != "auxiliary" {
		return Lease{}, fmt.Errorf("unknown kind")
	}
	if kind == "primary" && parent != 0 {
		return Lease{}, fmt.Errorf("primary cannot have parent")
	}
	if kind == "auxiliary" {
		p, ok := m.active[parent]
		if !ok || p.User != user || p.Host != host || p.Kind != "primary" || !p.Established {
			return Lease{}, fmt.Errorf("established owned primary required")
		}
	}
	if len(m.active) >= m.capacity {
		return Lease{}, fmt.Errorf("capacity exhausted")
	}
	used := make(map[int]bool)
	for _, l := range m.active {
		if l.Host == host {
			used[l.Number] = true
		}
	}
	number := 1
	for used[number] {
		number++
	}
	if number > 255 {
		return Lease{}, fmt.Errorf("host link numbers exhausted")
	}
	m.next++
	l := Lease{ID: m.next, User: user, Host: host, Number: number, Kind: kind, Parent: parent}
	m.active[l.ID] = l
	return l, nil
}
func (m *Links) Confirm(user string, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.active[id]
	if !ok || l.User != user {
		return fmt.Errorf("unknown or unowned lease")
	}
	l.Established = true
	m.active[id] = l
	return nil
}
func (m *Links) Close(user string, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.active[id]
	if !ok || l.User != user {
		return fmt.Errorf("unknown or unowned lease")
	}
	delete(m.active, id)
	for k, v := range m.active {
		if v.Parent == id {
			delete(m.active, k)
		}
	}
	return nil
}
func (m *Links) Snapshot() []Lease {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Lease, 0, len(m.active))
	for _, l := range m.active {
		out = append(out, l)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
