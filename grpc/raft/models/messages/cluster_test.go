package messages

import (
	"testing"
)

func TestStreamNodeStateArgumentsRoundTrip(t *testing.T) {
	args := StreamNodeStateArguments{ID: 42}

	proto := args.ToProto()
	result := StreamNodeStateArgumentsFromProto(proto)

	if result.ID != args.ID {
		t.Errorf("ID mismatch: got %d, want %d", result.ID, args.ID)
	}
}

func TestNodeInfoArgumentsRoundTripWithAddr(t *testing.T) {
	addr := "localhost:5001"
	args := NodeInfoArguments{
		ID:   7,
		Addr: &addr,
	}

	proto := args.ToProto()
	result := NodeInfoArgumentsFromProto(proto)

	if result.ID != args.ID {
		t.Errorf("ID mismatch: got %d, want %d", result.ID, args.ID)
	}
	if result.Addr == nil || *result.Addr != *args.Addr {
		t.Errorf("Addr mismatch: got %v, want %s", result.Addr, *args.Addr)
	}
}

func TestNodeInfoArgumentsRoundTripWithoutAddr(t *testing.T) {
	args := NodeInfoArguments{
		ID:   3,
		Addr: nil,
	}

	proto := args.ToProto()
	result := NodeInfoArgumentsFromProto(proto)

	if result.ID != args.ID {
		t.Errorf("ID mismatch: got %d, want %d", result.ID, args.ID)
	}
	if result.Addr != nil {
		t.Errorf("expected Addr to be nil, got %s", *result.Addr)
	}
}
