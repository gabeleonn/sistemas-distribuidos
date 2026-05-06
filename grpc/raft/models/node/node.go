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

// NewNode creates a new Raft node with the given ID and address, and loads its persistent state from disk.
func NewNode(id int64, addr string) *Node {
	node := Node{
		currentTerm: 0,
		votedFor:    nil,
		log:         []log.LogEntry{},
		commitIndex: 0,
		lastApplied: 0,
		nextIndex:   make(map[int64]int64),
		matchIndex:  make(map[int64]int64),
		role:        Follower,
		leaderID:    nil,
		ID:          id,
		Addr:        addr,
	}

	node.loadPersistentState()

	return &node
}

// BeforeRemove is called before the node is removed from the cluster, allowing it to save its persistent state to disk.
func (n *Node) BeforeRemove() error {
	return n.savePersistentState()
}

func (n *Node) loadPersistentState() {
	path := storage.GetNodePath(n.ID)
	data, err := storage.LoadPersistentNodeState(path)
	if err == nil {
		n.currentTerm = data.CurrentTerm
		n.votedFor = data.VotedFor
		n.log = data.Log
	}
}

func (n *Node) savePersistentState() error {
	path := storage.GetNodePath(n.ID)
	data := storage.NodePersistentState{
		CurrentTerm: n.currentTerm,
		VotedFor:    n.votedFor,
		Log:         n.log,
	}
	return storage.SavePersistentNodeState(path, data)
}

func (n *Node) getLastLogIndex() int64 {
	if len(n.log) == 0 {
		return 0
	}
	return n.log[len(n.log)-1].Index
}

// ToNodeState converts the Node to a NodeState, which is used for streaming node states to clients.
func (n *Node) ToNodeState() *State {
	return &State{
		ID:           n.ID,
		Addr:         n.Addr,
		Role:         n.role,
		Term:         n.currentTerm,
		LeaderID:     n.leaderID,
		CommitIndex:  n.commitIndex,
		LastLogIndex: n.getLastLogIndex(),
		LastApplied:  n.lastApplied,
		Status:       Running,
	}
}
