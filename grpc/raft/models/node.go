package models

import (
	pb "raft/autogen"
)

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

func (e LogEntry) ToProto() *pb.LogEntry {
	return &pb.LogEntry{
		Index:   e.Index,
		Term:    e.Term,
		Command: e.Command.ToProto(),
	}
}

func LogEntryFromProto(p *pb.LogEntry) LogEntry {
	return LogEntry{
		Index:   p.GetIndex(),
		Term:    p.GetTerm(),
		Command: CommandFromProto(p.GetCommand()),
	}
}
