package messages

import (
	pb "raft/autogen"
	"raft/models/log"
	"testing"
)

func TestCommandExecutionArgumentsToProto(t *testing.T) {
	cmd := CommandExecutionArguments{
		Command: log.Command{
			Type:  log.SET,
			Key:   "testKey",
			Value: func() *string { v := "testValue"; return &v }(),
		},
	}

	protoCmd := cmd.ToProto()

	if protoCmd.GetCommand() != pb.CommandType(log.SET) {
		t.Errorf("command type mismatch: got %v, want %v", protoCmd.GetCommand(), pb.CommandType(log.SET))
	}
	if protoCmd.GetKey() != "testKey" {
		t.Errorf("key mismatch: got %s, want %s", protoCmd.GetKey(), "testKey")
	}
	if protoCmd.GetValue() != "testValue" {
		t.Errorf("value mismatch: got %v, want %v", protoCmd.GetValue(), "testValue")
	}
}

func TestCommandExecutionArgumentsRoundTrip(t *testing.T) {
	value := "testValue"
	cmd := CommandExecutionArguments{
		Command: log.Command{
			Type:  log.SET,
			Key:   "testKey",
			Value: &value,
		},
	}

	proto := cmd.ToProto()
	result := CommandExecutionArgumentsFromProto(proto)

	if result.Command.Type != cmd.Command.Type {
		t.Errorf("command type mismatch: got %v, want %v", result.Command.Type, cmd.Command.Type)
	}
	if result.Command.Key != cmd.Command.Key {
		t.Errorf("key mismatch: got %s, want %s", result.Command.Key, cmd.Command.Key)
	}
	if result.Command.Value == nil || *result.Command.Value != *cmd.Command.Value {
		t.Errorf("value mismatch: got %v, want %s", result.Command.Value, *cmd.Command.Value)
	}
}

func TestCommandExecutionArgumentsRoundTripWithoutValue(t *testing.T) {
	cmd := CommandExecutionArguments{
		Command: log.Command{
			Type:  log.DELETE,
			Key:   "testKey",
			Value: nil,
		},
	}

	proto := cmd.ToProto()
	result := CommandExecutionArgumentsFromProto(proto)

	if result.Command.Type != cmd.Command.Type {
		t.Errorf("command type mismatch: got %v, want %v", result.Command.Type, cmd.Command.Type)
	}
	if result.Command.Key != cmd.Command.Key {
		t.Errorf("key mismatch: got %s, want %s", result.Command.Key, cmd.Command.Key)
	}
	if result.Command.Value != nil {
		t.Errorf("expected value to be nil, got %s", *result.Command.Value)
	}
}

func TestCommandExecutionReplyRoundTrip(t *testing.T) {
	value := "resultValue"
	leaderID := int64(3)
	reply := CommandExecutionReply{
		Success:  true,
		Value:    &value,
		LeaderID: &leaderID,
	}

	proto := reply.ToProto()
	result := CommandExecutionReplyFromProto(proto)

	if result.Success != reply.Success {
		t.Errorf("success mismatch: got %v, want %v", result.Success, reply.Success)
	}
	if result.Value == nil || *result.Value != *reply.Value {
		t.Errorf("value mismatch: got %v, want %s", result.Value, *reply.Value)
	}
	if result.LeaderID == nil || *result.LeaderID != *reply.LeaderID {
		t.Errorf("leaderID mismatch: got %v, want %d", result.LeaderID, *reply.LeaderID)
	}
}

func TestCommandExecutionReplyRoundTripFailureWithLeaderRedirect(t *testing.T) {
	leaderID := int64(2)
	reply := CommandExecutionReply{
		Success:  false,
		Value:    nil,
		LeaderID: &leaderID,
	}

	proto := reply.ToProto()
	result := CommandExecutionReplyFromProto(proto)

	if result.Success {
		t.Errorf("expected success to be false")
	}
	if result.Value != nil {
		t.Errorf("expected value to be nil")
	}
	if result.LeaderID == nil || *result.LeaderID != leaderID {
		t.Errorf("leaderID mismatch: got %v, want %d", result.LeaderID, leaderID)
	}
}

func TestCommandExecutionReplyRoundTripWithoutOptionals(t *testing.T) {
	reply := CommandExecutionReply{
		Success:  true,
		Value:    nil,
		LeaderID: nil,
	}

	proto := reply.ToProto()
	result := CommandExecutionReplyFromProto(proto)

	if !result.Success {
		t.Errorf("expected success to be true")
	}
	if result.Value != nil {
		t.Errorf("expected value to be nil, got %s", *result.Value)
	}
	if result.LeaderID != nil {
		t.Errorf("expected leaderID to be nil, got %d", *result.LeaderID)
	}
}
