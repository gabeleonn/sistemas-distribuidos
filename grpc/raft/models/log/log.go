package log

import (
	pb "raft/autogen"
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

// LogEntry represents an entry in the Raft log.
type LogEntry struct {
	Index   int64
	Term    int64
	Command Command
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
