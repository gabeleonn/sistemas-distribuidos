package store

import (
	"fmt"
	"goraft/raft"
	"sync"
)

// Store is a simple in-memory key-value store with thread-safe access.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// Apply applies a command to the store, modifying its state accordingly. It returns an error if the command is invalid or cannot be applied.
func (s *Store) Apply(cmd raft.Command) error {
	switch cmd.Type {
	case raft.SET:
		if cmd.Value == nil {
			return fmt.Errorf("SET command requires a value")
		}
		s.set(cmd.Key, *cmd.Value)
	case raft.DELETE:
		s.delete(cmd.Key)
	default:
		return fmt.Errorf("unknown command type: %v", cmd.Type)
	}

	return nil
}

// Get retrieves the value associated with a key from the store. It returns the value and a boolean indicating whether the key exists in the store.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

func (s *Store) set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// NewStore creates a new instance of the Store.
func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}
