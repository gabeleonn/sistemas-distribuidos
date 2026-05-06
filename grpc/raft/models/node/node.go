package node

import (
	"raft/models/log"
	"raft/storage"
)

// Role represents the role of a node in the Raft cluster.
type Role int

// Node roles
const (
	Follower Role = iota
	Candidate
	Leader
)

// Node represents a Raft node with its state and metadata.
type Node struct {
	// Persistent state on all servers
	currentTerm int64
	votedFor    *int64 // optional
	log         []log.LogEntry

	// Volatile state on all servers
	commitIndex int64
	lastApplied int64

	// Volatile state on leaders
	nextIndex  map[int64]int64
	matchIndex map[int64]int64

	// Meta information
	role     Role
	leaderID *int64 // optional
	ID       int64
	Addr     string
}

// LoadPersistentState returns the persistent state of the node for storage.
func (n *Node) LoadPersistentState() {
	path := storage.GetNodePath(n.ID)
	data, err := storage.LoadPersistentNodeState(path)
	if err == nil {
		n.currentTerm = data.CurrentTerm
		n.votedFor = data.VotedFor
		n.log = data.Log
	}
}

// SavePersistentState saves the current persistent state of the node to disk.
func (n *Node) SavePersistentState() {
	path := storage.GetNodePath(n.ID)
	data := storage.NodePersistentState{
		CurrentTerm: n.currentTerm,
		VotedFor:    n.votedFor,
		Log:         n.log,
	}
	storage.SavePersistentNodeState(path, data)
}
