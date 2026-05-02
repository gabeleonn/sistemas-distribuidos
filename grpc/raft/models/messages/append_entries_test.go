package messages

import (
	"raft/models/node"
	"raft/models/utils"

	"testing"
)

func TestAppendEntriesRequestRoundTrip(t *testing.T) {
	req := AppendEntriesArguments{
		Term:         1,
		LeaderID:     123,
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		Entries: []node.LogEntry{
			{
				Index: 1,
				Term:  1,
				Command: node.Command{
					Type: node.SET,
					Key:  "key1", Value: utils.GetPtrFromString("value1"),
				},
			},
			{
				Index: 2,
				Term:  1,
				Command: node.Command{
					Type: node.DELETE,
					Key:  "key2", Value: nil,
				},
			},
		},
		LeaderCommit: 0,
	}

	protoReq := req.ToProto()
	modelReq := AppendEntriesArgumentsFromProto(protoReq)

	if modelReq.Term != req.Term {
		t.Errorf("term mismatch")
	}
	if modelReq.LeaderID != req.LeaderID {
		t.Errorf("leader ID mismatch")
	}
	if modelReq.PrevLogIndex != req.PrevLogIndex {
		t.Errorf("prev log index mismatch")
	}
	if modelReq.PrevLogTerm != req.PrevLogTerm {
		t.Errorf("prev log term mismatch")
	}
	if modelReq.LeaderCommit != req.LeaderCommit {
		t.Errorf("leader commit mismatch")
	}
	if len(modelReq.Entries) != len(req.Entries) {
		t.Errorf("entries length mismatch: got %d, want %d", len(modelReq.Entries), len(req.Entries))
	}

	for i, entry := range modelReq.Entries {
		if entry.Index != req.Entries[i].Index {
			t.Errorf("entry %d index mismatch: got %d, want %d", i, entry.Index, req.Entries[i].Index)
		}
		if entry.Term != req.Entries[i].Term {
			t.Errorf("entry %d term mismatch: got %d, want %d", i, entry.Term, req.Entries[i].Term)
		}
		if entry.Command.Type != req.Entries[i].Command.Type {
			t.Errorf("entry %d command type mismatch: got %v, want %v", i, entry.Command.Type, req.Entries[i].Command.Type)
		}
		if entry.Command.Key != req.Entries[i].Command.Key {
			t.Errorf("entry %d command key mismatch: got %s, want %s", i, entry.Command.Key, req.Entries[i].Command.Key)
		}
		if (entry.Command.Value == nil) != (req.Entries[i].Command.Value == nil) {
			t.Errorf("entry %d command value nil mismatch: got %v, want %v", i, entry.Command.Value == nil, req.Entries[i].Command.Value == nil)
		} else if entry.Command.Value != nil && *entry.Command.Value != *req.Entries[i].Command.Value {
			t.Errorf("entry %d command value mismatch: got %s, want %s", i, *entry.Command.Value, *req.Entries[i].Command.Value)
		}
	}
}
