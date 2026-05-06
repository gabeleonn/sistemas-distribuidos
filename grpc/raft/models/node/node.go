package node

import "raft/models/log"

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
