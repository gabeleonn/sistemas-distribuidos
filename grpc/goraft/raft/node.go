package raft

import (
	"sync"
)

// Config holds the configuration for a Raft node, including its ID and the state machine it will use to
// apply commands from the log.
type Config struct {
	ID           int64
	StateMachine StateMachine
}

// Role represents the role of a Raft node, which can be Follower, Candidate, or Leader.
type Role string

// Node represents a Raft node in the cluster, containing its ID, and other relevant information.
const (
	Follower  Role = "Follower"
	Candidate Role = "Candidate"
	Leader    Role = "Leader"
)

// Node represents a Raft node in the cluster, containing its ID, role, term, and other relevant information.
type Node struct {
	mu sync.RWMutex

	id   int64
	role Role

	currentTerm int64
	votedFor    *int64
	log         []LogEntry

	commitIndex int64
	lastApplied int64

	nextIndex  map[int64]int64
	matchIndex map[int64]int64

	stateMachine StateMachine
}

// ID returns the ID of the node.
func (n *Node) ID() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.id
}

// Role returns the current role of the node (Follower, Candidate, or Leader).
func (n *Node) Role() Role {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.role
}

// CurrentTerm returns the current term of the node.
func (n *Node) CurrentTerm() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.currentTerm
}

// VotedFor returns the ID of the candidate that received the node's vote in the current term,
// or nil if the node has not voted for any candidate.
func (n *Node) VotedFor() (int64, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.votedFor == nil {
		return 0, false
	}
	return *n.votedFor, true
}

// Log returns the log entries of the node.
func (n *Node) Log() []LogEntry {
	n.mu.RLock()
	defer n.mu.RUnlock()
	logCopy := make([]LogEntry, len(n.log))
	copy(logCopy, n.log)

	return logCopy
}

// CommitIndex returns the index of the highest log entry known to be committed.
func (n *Node) CommitIndex() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.commitIndex
}

// LastApplied returns the index of the highest log entry applied to the state machine.
func (n *Node) LastApplied() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.lastApplied
}

// NextIndex returns the index of the next log entry to send to a given peer.
func (n *Node) NextIndex(peerID int64) int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.nextIndex[peerID]
}

// MatchIndex returns the index of the highest log entry known to be replicated on a given peer.
func (n *Node) MatchIndex(peerID int64) int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.matchIndex[peerID]
}

// RequestVote handles incoming RequestVote RPCs from other nodes in the cluster, allowing the node
// to participate in leader elections by granting or denying votes based on its current state and the
// information provided in the request.
func (n *Node) RequestVote(req RequestVoteRequest) RequestVoteResponse {
	n.mu.Lock()
	defer n.mu.Unlock()

	if req.Term < n.currentTerm {
		return RequestVoteResponse{
			Term:        n.currentTerm,
			VoteGranted: false,
		}
	}

	if req.Term > n.currentTerm {
		n.currentTerm = req.Term
		n.role = Follower
		n.votedFor = nil
	}

	if n.votedFor == nil || *n.votedFor == req.CandidateID {
		n.votedFor = &req.CandidateID

		return RequestVoteResponse{
			Term:        n.currentTerm,
			VoteGranted: true,
		}
	}

	return RequestVoteResponse{
		Term:        n.currentTerm,
		VoteGranted: false,
	}
}

// NewNode creates a new Node instance with the given ID and initializing it as a Follower with default
// values for other fields.
func NewNode(config Config) *Node {
	return &Node{
		id:           config.ID,
		role:         Follower,
		currentTerm:  0,
		votedFor:     nil,
		log:          []LogEntry{},
		commitIndex:  0,
		lastApplied:  0,
		nextIndex:    make(map[int64]int64),
		matchIndex:   make(map[int64]int64),
		stateMachine: config.StateMachine,
	}
}
