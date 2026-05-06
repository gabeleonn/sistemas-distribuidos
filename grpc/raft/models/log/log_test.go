package log

import (
	"raft/models/utils"
	"testing"
)

func TestCommandSetRoundTrip(t *testing.T) {
	original := Command{
		Type:  SET,
		Key:   "testKey",
		Value: utils.GetPtrFromString("testValue"),
	}

	protoCmd := original.ToProto()
	converted := CommandFromProto(protoCmd)

	if original.Type != converted.Type {
		t.Errorf("Expected Type %v, got %v", original.Type, converted.Type)
	}
	if original.Key != converted.Key {
		t.Errorf("Expected Key %s, got %s", original.Key, converted.Key)
	}
	if (original.Value == nil) != (converted.Value == nil) {
		t.Errorf("Expected Value nil status %v, got %v", original.Value == nil, converted.Value == nil)
	} else if original.Value != nil && *original.Value != *converted.Value {
		t.Errorf("Expected Value %s, got %s", *original.Value, *converted.Value)
	}
}

func TestCommandDeleteRoundTrip(t *testing.T) {
	original := Command{
		Type: DELETE,
		Key:  "testKey",
	}

	protoCmd := original.ToProto()
	converted := CommandFromProto(protoCmd)

	if original.Type != converted.Type {
		t.Errorf("Expected Type %v, got %v", original.Type, converted.Type)
	}
	if original.Key != converted.Key {
		t.Errorf("Expected Key %s, got %s", original.Key, converted.Key)
	}
	if converted.Value != nil {
		t.Errorf("Expected Value nil, got %v", *converted.Value)
	}
}
