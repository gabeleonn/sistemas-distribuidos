package node

import (
	pb "raft/autogen"
	"raft/models/log"
	"raft/models/utils"
	"testing"
)

func TestCommandRoundTrip(t *testing.T) {
	value := "myvalue"

	p := &pb.LogEntryCommand{
		Type:  pb.LogEntryCommandType_LOG_ENTRY_COMMAND_SET,
		Key:   "mykey",
		Value: &value,
	}

	model := log.CommandFromProto(p)

	result := model.ToProto()

	if result.GetType() != p.GetType() {
		t.Errorf("type mismatch")
	}

	if result.GetKey() != p.GetKey() {
		t.Errorf("key mismatch: got %s, want %s", result.GetKey(), p.GetKey())
	}

	if result.GetValue() != utils.GetStringFromPtr(p.Value) {
		t.Errorf("value mismatch: got empty, want %s", *p.Value)
	}
}

func TestCommandRoundTripWithoutValue(t *testing.T) {
	p := &pb.LogEntryCommand{
		Type: pb.LogEntryCommandType_LOG_ENTRY_COMMAND_DELETE,
		Key:  "mykey",
	}

	model := log.CommandFromProto(p)

	result := model.ToProto()

	if result.GetType() != p.GetType() {
		t.Errorf("type mismatch")
	}

	if result.GetKey() != p.GetKey() {
		t.Errorf("key mismatch: got %s, want %s", result.GetKey(), p.GetKey())
	}

	if result.Value != nil {
		t.Errorf("expected value to be nil, got %s", *result.Value)
	}
}

func TestCommandRoundTripWithNilValue(t *testing.T) {
	p := &pb.LogEntryCommand{
		Type:  pb.LogEntryCommandType_LOG_ENTRY_COMMAND_SET,
		Key:   "mykey",
		Value: nil,
	}

	model := log.CommandFromProto(p)

	result := model.ToProto()

	if result.GetType() != p.GetType() {
		t.Errorf("type mismatch")
	}

	if result.GetKey() != p.GetKey() {
		t.Errorf("key mismatch: got %s, want %s", result.GetKey(), p.GetKey())
	}

	if result.Value != nil {
		t.Errorf("expected value to be nil, got %s", *result.Value)
	}
}

func TestCommandTypeMapping(t *testing.T) {
	tests := []struct {
		proto pb.LogEntryCommandType
		model log.CommandType
	}{
		{pb.LogEntryCommandType_LOG_ENTRY_COMMAND_SET, log.SET},
		{pb.LogEntryCommandType_LOG_ENTRY_COMMAND_DELETE, log.DELETE},
	}

	for _, tt := range tests {
		p := &pb.LogEntryCommand{
			Type: tt.proto,
			Key:  "testkey",
		}

		model := log.CommandFromProto(p)

		if model.Type != tt.model {
			t.Errorf("mapping mismatch: got %v, want %v", model.Type, tt.model)
		}
	}
}

func TestLogEntryRoundTrip(t *testing.T) {
	value := "myvalue"

	p := &pb.LogEntry{
		Index: 1,
		Term:  2,
		Command: &pb.LogEntryCommand{
			Type:  pb.LogEntryCommandType_LOG_ENTRY_COMMAND_SET,
			Key:   "mykey",
			Value: &value,
		},
	}

	model := log.LogEntryFromProto(p)

	result := model.ToProto()

	if result.GetIndex() != p.GetIndex() {
		t.Errorf("index mismatch")
	}

	if result.GetTerm() != p.GetTerm() {
		t.Errorf("term mismatch")
	}

	if result.GetCommand().GetType() != p.GetCommand().GetType() {
		t.Errorf("command type mismatch")
	}

	if result.GetCommand().GetKey() != p.GetCommand().GetKey() {
		t.Errorf("command key mismatch: got %s, want %s", result.GetCommand().GetKey(), p.GetCommand().GetKey())
	}

	if result.GetCommand().GetValue() != utils.GetStringFromPtr(p.GetCommand().Value) {
		t.Errorf("command value mismatch: got empty, want %s", *p.GetCommand().Value)
	}
}
