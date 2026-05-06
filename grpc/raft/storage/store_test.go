package storage

import (
	"raft/models/log"
	"testing"
)

func TestKeyValueStoreGet(t *testing.T) {
	store := KeyValueStore
	store.set("key1", "value1")

	val, ok := store.Get("key1")
	if !ok {
		t.Errorf("expected key to exist")
	}
	if val != "value1" {
		t.Errorf("expected value to be 'value1', got '%s'", val)
	}
}

func TestKeyValueStoreApplySet(t *testing.T) {
	store := KeyValueStore
	cmd := log.Command{
		Type:  log.SET,
		Key:   "key2",
		Value: func() *string { v := "value2"; return &v }(),
	}

	err := store.Apply(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	val, ok := store.Get("key2")
	if !ok {
		t.Errorf("expected key to exist")
	}
	if val != "value2" {
		t.Errorf("expected value to be 'value2', got '%s'", val)
	}
}

func TestKeyValueStoreApplyDelete(t *testing.T) {
	store := KeyValueStore
	store.set("key3", "value3")

	cmd := log.Command{
		Type: log.DELETE,
		Key:  "key3",
	}

	err := store.Apply(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, ok := store.Get("key3")
	if ok {
		t.Errorf("expected key to be deleted")
	}
}
