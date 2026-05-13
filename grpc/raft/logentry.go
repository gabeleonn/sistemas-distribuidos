package raft

import (
	"goraft/proto"
)

type LogEntry struct {
	Index   int64
	Term    int64
	Command Command
}

func (le *LogEntry) ToProto() *proto.LogEntry {
	return &proto.LogEntry{
		Index:   le.Index,
		Term:    le.Term,
		Command: le.Command.ToProto(),
	}
}

func LogEntryFromProto(p *proto.LogEntry) (*LogEntry, error) {
	cmd, err := CommandFromProto(p.Command)
	if err != nil {
		return nil, err
	}

	return &LogEntry{
		Index:   p.Index,
		Term:    p.Term,
		Command: cmd,
	}, nil
}

func NewLogEntry(index, term int64, command Command) LogEntry {
	return LogEntry{
		Index:   index,
		Term:    term,
		Command: command,
	}
}
