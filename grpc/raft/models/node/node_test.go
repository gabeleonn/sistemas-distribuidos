package node

import (
	"os"
	"path"
	"raft/models/log"
	"testing"
)

func TestNodePersistentState(t *testing.T) {
	node := Node{
		currentTerm: 5,
		votedFor:    nil,
		log:         []log.LogEntry{},
		ID:          1,
		Addr:        "localhost:8080",
	}

	node.LoadPersistentState()

	if node.currentTerm != 5 {
		t.Errorf("Expected CurrentTerm %d, got %d", 5, node.currentTerm)
	}
	if node.votedFor != nil {
		t.Errorf("Expected VotedFor %v, got %v", nil, node.votedFor)
	}
	if len(node.log) != 0 {
		t.Errorf("Expected Log length %d, got %d", 0, len(node.log))
	}
}

func TestNodeSaveAndLoadPersistentState(t *testing.T) {
	defer os.Remove(path.Join(os.TempDir(), "goraft_node_2.gob"))
	votedFor := int64(3)
	node := Node{
		currentTerm: 10,
		votedFor:    &votedFor,
		log: []log.LogEntry{
			{
				Index:   1,
				Term:    10,
				Command: log.Command{Type: log.SET, Key: "key1", Value: nil},
			},
		},
		ID:   2,
		Addr: "localhost:8081",
	}

	node.SavePersistentState()

	loadedNode := Node{ID: 2}
	loadedNode.LoadPersistentState()

	if loadedNode.currentTerm != node.currentTerm {
		t.Errorf("Expected CurrentTerm %d, got %d", node.currentTerm, loadedNode.currentTerm)
	}
}
