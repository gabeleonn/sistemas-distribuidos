package messages

import (
	pb "raft/autogen"
	"raft/models/node"
	"testing"
)

func TestCommandExecutionArgumentsToProto(t *testing.T) {
	cmd := CommandExecutionArguments{
		Command: node.Command{
			Type:  node.SET,
			Key:   "testKey",
			Value: func() *string { v := "testValue"; return &v }(),
		},
	}

	protoCmd := cmd.ToProto()

	if protoCmd.GetCommand() != pb.CommandType(node.SET) {
		t.Errorf("command type mismatch: got %v, want %v", protoCmd.GetCommand(), pb.CommandType(node.SET))
	}
	if protoCmd.GetKey() != "testKey" {
		t.Errorf("key mismatch: got %s, want %s", protoCmd.GetKey(), "testKey")
	}
	if protoCmd.GetValue() != "testValue" {
		t.Errorf("value mismatch: got %v, want %v", protoCmd.GetValue(), "testValue")
	}
}
