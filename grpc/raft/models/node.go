package models

type Metadata struct {
	ID   int64
	Addr string	
}

type NodeRole int

const (
	Follower NodeRole = iota
	Candidate
	Leader
)

type CommandType int

const (
	SET CommandType = iota
	DELETE
)

type Command struct {
	Type  CommandType
	Key   string
	Value *string // Value is optional, only needed for SET commands
}

type LogEntry struct {
	Index int64
	Term int64
	Command Command
}

type Node struct {
	// Persistent state on all servers
	currentTerm int64
	votedFor *int64 // optional
	log []LogEntry

	// Volatile state on all servers
	commitIndex int64
	lastApplied int64
	
	// Volatile state on leaders
	nextIndex map[int64]int64
	matchIndex map[int64]int64
	
	// Meta information
	role NodeRole
	leaderId *int64 // optional
	Metadata Metadata
}
