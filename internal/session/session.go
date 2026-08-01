package session

import (
	"net"
	"sync"
	"time"
)

type Session struct {
	ClientAddr string
	Key        []byte
	CreatedAt  time.Time
	LastSeen   time.Time
}

func (s *Session) Touch() {
	s.LastSeen = time.Now().UTC()
}

func (s *Session) Expired(timeout time.Duration) bool {
	return time.Since(s.LastSeen) > timeout
}

type Manager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	timeout  time.Duration
}

func NewManager(timeout time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		timeout:  timeout,
	}
}

func (m *Manager) Create(addr net.Addr, key []byte) *Session {
	now := time.Now().UTC()

	s := &Session{
		ClientAddr: addr.String(),
		Key:        append([]byte(nil), key...),
		CreatedAt:  now,
		LastSeen:   now,
	}

	m.mu.Lock()
	m.sessions[s.ClientAddr] = s
	m.mu.Unlock()

	return s
}

func (m *Manager) Get(addr net.Addr) (*Session, bool) {
	m.mu.RLock()
	s, ok := m.sessions[addr.String()]
	m.mu.RUnlock()

	if !ok {
		return nil, false
	}

	return s, true
}

func (m *Manager) Delete(addr net.Addr) {
	m.mu.Lock()
	delete(m.sessions, addr.String())
	m.mu.Unlock()
}

func (m *Manager) Touch(addr net.Addr) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[addr.String()]
	if !ok {
		return false
	}

	s.Touch()
	return true
}

func (m *Manager) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	for addr, s := range m.sessions {
		if now.Sub(s.LastSeen) > m.timeout {
			delete(m.sessions, addr)
		}
	}

}

func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.sessions)
}
