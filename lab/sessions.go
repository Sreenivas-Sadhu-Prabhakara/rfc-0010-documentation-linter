package lab

import (
	"fmt"
	"sort"
	"sync"
)

// Sessions is an educational, in-memory lifecycle and inbound routing model.
// Caller-supplied user names are labels, not authenticated identities.
type Sessions struct {
	mu     sync.Mutex
	links  *Links
	depth  int
	active map[uint64]*session
	routes map[route]uint64
}
type route struct {
	host string
	link int
}
type session struct {
	lease    Lease
	incoming int
	inbox    []string
}
type SessionView struct {
	Lease
	Incoming int `json:"incoming"`
	Queued   int `json:"queued"`
}

func NewSessions(capacity, depth int) (*Sessions, error) {
	links, err := NewLinks(capacity)
	if err != nil {
		return nil, err
	}
	if depth < 1 || depth > 64 {
		return nil, fmt.Errorf("inbox depth must be 1..64")
	}
	return &Sessions{links: links, depth: depth, active: map[uint64]*session{}, routes: map[route]uint64{}}, nil
}

// Open reserves an outgoing link. An ACK is required before data may arrive.
// A primary is unique per (user, remote host); auxiliaries share its lifetime.
func (m *Sessions) Open(user, host, kind string, parent uint64) (Lease, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if kind == "primary" {
		for _, s := range m.active {
			if s.lease.User == user && s.lease.Host == host && s.lease.Kind == "primary" {
				return Lease{}, fmt.Errorf("primary already pending or established")
			}
		}
	}
	l, err := m.links.Open(user, host, kind, parent)
	if err != nil {
		return Lease{}, err
	}
	m.active[l.ID] = &session{lease: l}
	return l, nil
}

// Ack correlates an explicit modern generation ID and remote host with a request.
// Replaying the same ACK is harmless; conflicting ACKs never overwrite routing.
func (m *Sessions) Ack(user string, id uint64, host string, incoming int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.active[id]
	if !ok || s.lease.User != user || s.lease.Host != host {
		return fmt.Errorf("unknown or mismatched request")
	}
	if incoming < 1 || incoming > 255 {
		return fmt.Errorf("incoming link must be 1..255")
	}
	if s.lease.Established {
		if s.incoming == incoming {
			return nil
		}
		return fmt.Errorf("conflicting acknowledgment")
	}
	k := route{host, incoming}
	if _, exists := m.routes[k]; exists {
		return fmt.Errorf("incoming link already routed")
	}
	if err := m.links.Confirm(user, id); err != nil {
		return err
	}
	s.lease.Established = true
	s.incoming = incoming
	m.routes[k] = id
	return nil
}

// Reject abandons only a pending request; no automatic retries are attempted.
func (m *Sessions) Reject(user string, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.active[id]
	if !ok || s.lease.User != user || s.lease.Established {
		return fmt.Errorf("owned pending request required")
	}
	return m.close(user, id)
}
func (m *Sessions) Close(user string, id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.close(user, id)
}
func (m *Sessions) close(user string, id uint64) error {
	if err := m.links.Close(user, id); err != nil {
		return err
	}
	for key, s := range m.active {
		if key == id || s.lease.Parent == id {
			if s.lease.Established {
				delete(m.routes, route{s.lease.Host, s.incoming})
			}
			delete(m.active, key)
		}
	}
	return nil
}

// Deliver demonstrates host+link demultiplexing and explicit backpressure.
// Generation is a modern extension to reject delayed data after link reuse.
// 1024 bytes is a lab limit, unrelated to RFC 11's bit-level framing.
func (m *Sessions) Deliver(host string, incoming int, generation uint64, payload string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.routes[route{host, incoming}]
	if !ok || id != generation {
		return fmt.Errorf("unknown or stale route")
	}
	s := m.active[id]
	if len(payload) > 1024 {
		return fmt.Errorf("payload exceeds 1024 bytes")
	}
	if len(s.inbox) >= m.depth {
		return fmt.Errorf("inbox full: sender must retry explicitly")
	}
	s.inbox = append(s.inbox, payload)
	return nil
}
func (m *Sessions) Drain(user string, id uint64) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.active[id]
	if !ok || s.lease.User != user || !s.lease.Established {
		return nil, fmt.Errorf("owned established connection required")
	}
	out := append([]string{}, s.inbox...)
	s.inbox = nil
	return out, nil
}
func (m *Sessions) Snapshot() []SessionView {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]SessionView, 0, len(m.active))
	for _, s := range m.active {
		out = append(out, SessionView{s.lease, s.incoming, len(s.inbox)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
