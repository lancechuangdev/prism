package auth

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, username string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type memorySession struct {
	username  string
	expiresAt time.Time
}

type MemorySessionStore struct {
	mu       sync.Mutex
	sessions map[string]memorySession
	now      func() time.Time
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		sessions: make(map[string]memorySession),
		now:      time.Now,
	}
}

func (s *MemorySessionStore) Get(_ context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[key]
	if !ok {
		return "", ErrSessionNotFound
	}
	if !session.expiresAt.After(s.now().UTC()) {
		delete(s.sessions, key)
		return "", ErrSessionNotFound
	}
	return session.username, nil
}

func (s *MemorySessionStore) Set(_ context.Context, key string, username string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[key] = memorySession{
		username:  username,
		expiresAt: s.now().UTC().Add(ttl),
	}
	return nil
}

func (s *MemorySessionStore) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, key)
	return nil
}
