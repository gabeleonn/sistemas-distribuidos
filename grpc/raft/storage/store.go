package storage

import (
	"errors"
	"raft/models/node"
	"sync"
)

// KeyValueStoreType is a simple in-memory key-value store with thread-safe access.
type KeyValueStoreType struct {
	data map[string]string
	mu   sync.RWMutex
}

// KeyValueStore is a global instance of KeyValueStoreType that can be used throughout the application.
var KeyValueStore = &KeyValueStoreType{
	data: make(map[string]string),
}

// Apply applies a committed log entry command to the state machine.
func (kv *KeyValueStoreType) Apply(cmd node.Command) error {
	switch cmd.Type {
	case node.SET:
		if cmd.Value == nil {
			return errors.New("SET command requires a value")
		}
		kv.set(cmd.Key, *cmd.Value)
	case node.DELETE:
		kv.delete(cmd.Key)
	default:
		return errors.New("unknown command type")
	}
	return nil
}

// Get retrieves the value associated with a key from the store.
func (kv *KeyValueStoreType) Get(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	val, ok := kv.data[key]
	return val, ok
}

func (kv *KeyValueStoreType) set(key, value string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key] = value
}

func (kv *KeyValueStoreType) delete(key string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.data, key)
}
