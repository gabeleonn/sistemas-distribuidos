package node

import (
	pb "raft/autogen"
)

// Metadata contains basic information about a node in the cluster.
type Metadata struct {
	ID   int64
	Addr string
}

// Role represents the role of a node in the Raft cluster.
type Role int

// Node roles
const (
	Follower Role = iota
	Candidate
	Leader
)

// CommandType represents the type of a command in the log entry.
type CommandType int

// Command types
const (
	SET CommandType = iota
	DELETE
)

// Command represents a command to be executed on the state machine.
type Command struct {
	Type  CommandType
	Key   string
	Value *string // Value is optional, only needed for SET commands
}

// LogEntry represents an entry in the Raft log.
type LogEntry struct {
	Index   int64
	Term    int64
	Command Command
}

// Node represents a Raft node with its state and metadata.
type Node struct {
	// Persistent state on all servers
	currentTerm int64
	votedFor    *int64 // optional
	log         []LogEntry

	// Volatile state on all servers
	commitIndex int64
	lastApplied int64

	// Volatile state on leaders
	nextIndex  map[int64]int64
	matchIndex map[int64]int64

	// Meta information
	role     Role
	leaderID *int64 // optional
	Metadata Metadata
}

// ToProto converts a Command to its protobuf representation.
func (c Command) ToProto() *pb.LogEntryCommand {
	cmd := &pb.LogEntryCommand{
		Type: pb.LogEntryCommandType(c.Type),
		Key:  c.Key,
	}
	if c.Value != nil {
		cmd.Value = c.Value
	}
	return cmd
}

// CommandFromProto converts a protobuf LogEntryCommand to a Command.
func CommandFromProto(p *pb.LogEntryCommand) Command {
	cmd := Command{
		Type: CommandType(p.GetType()),
		Key:  p.GetKey(),
	}
	if p.Value != nil {
		cmd.Value = p.Value
	}
	return cmd
}

// ToProto converts a LogEntry to its protobuf representation.
func (e LogEntry) ToProto() *pb.LogEntry {
	return &pb.LogEntry{
		Index:   e.Index,
		Term:    e.Term,
		Command: e.Command.ToProto(),
	}
}

// LogEntryFromProto converts a protobuf LogEntry to a LogEntry.
func LogEntryFromProto(p *pb.LogEntry) LogEntry {
	return LogEntry{
		Index:   p.GetIndex(),
		Term:    p.GetTerm(),
		Command: CommandFromProto(p.GetCommand()),
	}
}
