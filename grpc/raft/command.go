package raft

import (
	"fmt"
	"goraft/proto"
	"strings"
)

type CommandType string

const (
	SET CommandType = "SET"
	DEL CommandType = "DEL"
	GET CommandType = "GET"
)

type Command struct {
	Type  CommandType
	Key   string
	Value *string // Value is optional, only needed for SET commands
}

func (c *Command) ToProto() string {
	switch c.Type {
	case SET:
		return commandToString(*c)
	case DEL:
		return commandToString(*c)
	case GET:
		return commandToString(*c)
	default:
		return ""
	}
}

func CommandFromProto(s string) (Command, error) {
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
	case GET:
		return fmt.Sprintf("GET:%s", cmd.Key)
	default:
		return ""
	}
}

type ExecuteCommandRequest struct {
	Command Command
}

type ExecuteCommandResponse struct {
	Success bool
	Message *string // Value is returned for GET commands
}

func (r *ExecuteCommandResponse) ToProto() *proto.CommandResponse {
	resp := &proto.CommandResponse{
		Success: r.Success,
	}

	if r.Message != nil {
		resp.Message = *r.Message
	}

	return resp
}

func ExecuteCommandResponseFromProto(p *proto.CommandResponse) *ExecuteCommandResponse {
	var message *string
	if p.Message != "" {
		message = &p.Message
	}

	return &ExecuteCommandResponse{
		Success: p.Success,
		Message: message,
	}
}
