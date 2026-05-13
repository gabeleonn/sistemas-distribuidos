package raft

import (
	"fmt"
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

type State struct {
	ID       int64
	Role     Role
	VotedFor *int64

	CurrentTerm int64
	Log         []LogEntry
	CommitIndex int64
}

// Node represents a Raft node in the cluster, containing its ID, role, term, and other relevant information.
type Node struct {
	mu sync.RWMutex

	id       int64
	role     Role
	leaderID *int64

	currentTerm int64
	votedFor    *int64
	log         []LogEntry

	commitIndex int64
	lastApplied int64

	nextIndex  map[int64]int64
	matchIndex map[int64]int64

	stateMachine StateMachine
}

/*
========================================== Getters ==========================================
*/

func (n *Node) ID() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.id
}

func (n *Node) Role() Role {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.role
}

func (n *Node) CurrentTerm() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.currentTerm
}

func (n *Node) VotedFor() (int64, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.votedFor == nil {
		return 0, false
	}
	return *n.votedFor, true
}

func (n *Node) Log() []LogEntry {
	n.mu.RLock()
	defer n.mu.RUnlock()
	logCopy := make([]LogEntry, len(n.log))
	copy(logCopy, n.log)

	return logCopy
}

func (n *Node) CommitIndex() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.commitIndex
}

func (n *Node) LastApplied() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.lastApplied
}

func (n *Node) NextIndex(peerID int64) int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.nextIndex[peerID]
}

func (n *Node) MatchIndex(peerID int64) int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.matchIndex[peerID]
}

func (n *Node) PreviousLogInfo(index int64) (int64, int64, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	prevLogIndex := index - 1
	if prevLogIndex == 0 {
		return 0, 0, nil
	}

	if prevLogIndex < 0 || prevLogIndex > int64(len(n.log)) {
		return 0, 0, fmt.Errorf("invalid log index: %d", prevLogIndex)
	}

	prevLogTerm := n.log[prevLogIndex-1].Term
	return prevLogIndex, prevLogTerm, nil
}

/*
========================================== Setters ==========================================
*/
func (n *Node) AppendLog(cmd Command) (LogEntry, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	currentIndex := int64(len(n.log)) + 1
	logentry := NewLogEntry(currentIndex, n.currentTerm, cmd)
	n.log = append(n.log, logentry)
	return logentry, nil
}

/*
========================================== Snapshot ==========================================
*/
func (n *Node) GetState() State {
	n.mu.RLock()
	defer n.mu.RUnlock()

	logCopy := make([]LogEntry, len(n.log))
	copy(logCopy, n.log)

	var votedFor *int64
	if n.votedFor != nil {
		value := *n.votedFor
		votedFor = &value
	}

	return State{
		ID:          n.id,
		Role:        n.role,
		VotedFor:    votedFor,
		CurrentTerm: n.currentTerm,
		Log:         logCopy,
		CommitIndex: n.commitIndex,
	}
}

func (n *Node) EnsureLeader() (int64, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if n.role != Leader {
		leaderAddress := "unknown"
		if n.leaderID != nil {
			id := *n.leaderID
			port := 50050 + id
			leaderAddress = fmt.Sprintf("localhost:%d", port)
		}

		return 0, fmt.Errorf("not the leader, redirect to %s", leaderAddress)
	}

	return n.currentTerm, nil
}

/*
========================================== Lifecycle ==========================================
*/
func (n *Node) BecomeCandidate() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.role = Candidate
	n.currentTerm++

	candidateID := n.id
	n.votedFor = &candidateID

	return n.currentTerm
}

func (n *Node) BecomeLeader() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.role = Leader
	return n.currentTerm
}

func (n *Node) BecomeFollower(term int64) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if term <= n.currentTerm {
		return
	}

	n.role = Follower
	n.currentTerm = term
	n.votedFor = nil
}

/*
========================================== gRPC Handlers ==========================================
*/
func (n *Node) CandidateResponse(req RequestVoteRequest) RequestVoteResponse {
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

func (n *Node) CandidateRequest() RequestVoteRequest {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return RequestVoteRequest{
		Term:        n.currentTerm,
		CandidateID: n.id,
	}
}

func (n *Node) HeartbeatResponse(req AppendEntriesRequest) AppendEntriesResponse {
	n.mu.Lock()
	defer n.mu.Unlock()

	if req.Term < n.currentTerm {
		return AppendEntriesResponse{
			Term:    n.currentTerm,
			Success: false,
		}
	}

	if req.Term > n.currentTerm {
		n.currentTerm = req.Term
		n.votedFor = nil
	}

	n.role = Follower
	leaderID := req.LeaderID
	n.leaderID = &leaderID

	return AppendEntriesResponse{
		Term:    n.currentTerm,
		Success: true,
	}
}

func (n *Node) HeartbeatRequest() AppendEntriesRequest {
	return AppendEntriesRequest{
		Term:     n.currentTerm,
		LeaderID: n.id,
	}
}

func (n *Node) ExecuteGetCommandResponse(cmd Command) ExecuteCommandResponse {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if n.role != Leader {
		leaderAddress := "unknown"
		if n.leaderID != nil {
			id := *n.leaderID
			port := 50050 + id
			leaderAddress = fmt.Sprintf("localhost:%d", port)
		}

		msg := fmt.Sprintf("not the leader, redirect to %s", leaderAddress)

		return ExecuteCommandResponse{
			Success: false,
			Message: &msg,
		}
	}

	response, _ := n.stateMachine.Get(cmd.Key)
	return ExecuteCommandResponse{
		Success: true,
		Message: &response,
	}
}

func (n *Node) ExecuteCommandResponse(logentry LogEntry) ExecuteCommandResponse {
	n.mu.Lock()
	defer n.mu.Unlock()

	if logentry.Index > n.commitIndex {
		n.commitIndex = logentry.Index
	}

	for n.lastApplied < n.commitIndex {
		n.lastApplied++

		entry := n.log[n.lastApplied-1]

		if err := n.stateMachine.Apply(entry.Command); err != nil {
			msg := fmt.Sprintf("failed to apply command at index %d: %v", entry.Index, err)

			return ExecuteCommandResponse{
				Success: false,
				Message: &msg,
			}
		}
	}

	message := "ok"

	return ExecuteCommandResponse{
		Success: true,
		Message: &message,
	}
}

/*
========================================== Helpers ==========================================
*/
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
