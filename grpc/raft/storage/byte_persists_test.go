package storage

import (
	"os"
	"path/filepath"
	"raft/models/log"
	"raft/models/utils"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	filepath := filepath.Join(os.TempDir(), "test_state.gob")
	defer os.Remove(filepath)
	state := NodePersistentState{
		CurrentTerm: 1,
		VotedFor:    nil,
		Log: []log.LogEntry{
			{
				Index:   1,
				Term:    1,
				Command: log.Command{Type: log.SET, Key: "key1", Value: utils.GetPtrFromString("value1")},
			},
		},
	}
	err := SavePersistentNodeState(filepath, state)
	if err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}
	loadedState, err := LoadPersistentNodeState(filepath)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}
	if loadedState.CurrentTerm != state.CurrentTerm {
		t.Errorf("Expected CurrentTerm %d, got %d", state.CurrentTerm, loadedState.CurrentTerm)
	}
	if loadedState.VotedFor != state.VotedFor {
		t.Errorf("Expected VotedFor %v, got %v", state.VotedFor, loadedState.VotedFor)
	}
	if len(loadedState.Log) != len(state.Log) {
		t.Errorf("Expected Log length %d, got %d", len(state.Log), len(loadedState.Log))
	}
}
