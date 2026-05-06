package node

import (
	"os"
	"path"
	"raft/models/log"
	"raft/models/utils"
	"testing"
)

func TestNodePersistentState(t *testing.T) {
	node := NewNode(2, "localhost:5002")

	node.currentTerm = 5
	node.votedFor = nil
	node.log = []log.LogEntry{}

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
	node := NewNode(2, "localhost:5002")

	node.currentTerm = 5
	node.votedFor = &votedFor
	node.log = []log.LogEntry{
		{Index: 1, Term: 1, Command: log.Command{Type: log.SET, Key: "x", Value: utils.GetPtrFromString("10")}},
		{Index: 2, Term: 2, Command: log.Command{Type: log.DELETE, Key: "y"}},
	}

	node.BeforeRemove()

	loadedNode := NewNode(2, "locahost:5000")

	if loadedNode.currentTerm != node.currentTerm {
		t.Errorf("Expected CurrentTerm %d, got %d", node.currentTerm, loadedNode.currentTerm)
	}

	if loadedNode.votedFor == nil || *loadedNode.votedFor != *node.votedFor {
		t.Errorf("Expected VotedFor %v, got %v", node.votedFor, loadedNode.votedFor)
	}

	if len(loadedNode.log) != len(node.log) {
		t.Fatalf("Expected Log length %d, got %d", len(node.log), len(loadedNode.log))
	}

	for i := range node.log {
		if loadedNode.log[i].Index != node.log[i].Index ||
			loadedNode.log[i].Term != node.log[i].Term ||
			loadedNode.log[i].Command.Type != node.log[i].Command.Type ||
			loadedNode.log[i].Command.Key != node.log[i].Command.Key ||
			((loadedNode.log[i].Command.Value == nil) != (node.log[i].Command.Value == nil)) ||
			(loadedNode.log[i].Command.Value != nil && *loadedNode.log[i].Command.Value != *node.log[i].Command.Value) {
			t.Errorf("Log entry %d does not match expected value", i)
		}
	}
}
