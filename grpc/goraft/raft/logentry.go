package raft

import (
	"fmt"
	"goraft/proto"
	"strings"
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
		Command: commandToString(le.Command),
	}
}

func LogEntryFromProto(p *proto.LogEntry) (*LogEntry, error) {
	cmd, err := stringToCommand(p.Command)
	if err != nil {
		return nil, err
	}

	return &LogEntry{
		Index:   p.Index,
		Term:    p.Term,
		Command: cmd,
	}, nil
}

func stringToCommand(s string) (Command, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) < 2 {
		return Command{}, fmt.Errorf("invalid command format: %s", s)
	}

	cmdType := CommandType(parts[0])
	key := parts[1]

	var value *string
	if cmdType == SET {
		if len(parts) != 3 {
			return Command{}, fmt.Errorf("SET command requires a value: %s", s)
		}
		value = &parts[2]
	}

	return Command{
		Type:  cmdType,
		Key:   key,
		Value: value,
	}, nil
}

func commandToString(cmd Command) string {
	switch cmd.Type {
	case SET:
		return fmt.Sprintf("SET:%s:%s", cmd.Key, *cmd.Value)
	case DEL:
		return fmt.Sprintf("DEL:%s", cmd.Key)
	default:
		return ""
	}
}
