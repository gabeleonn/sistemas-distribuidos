package cluster

import (
	pb "raft/autogen"
	"testing"
)

func TestClusterStateRoundTrip(t *testing.T) {
	p1id := int64(0)
	p2id := int64(1)

	p1 := createStreamNodeStateReply(&p1id)
	p2 := createStreamNodeStateReply(&p2id)

	// Proto → Model
	nodes := &pb.StreamNodesStatesReply{
		Nodes: []*pb.StreamNodeStateReply{p1, p2},
	}

	model := StateFromProto(nodes)

	// Model → Proto
	result := model.ToProto()

	if len(result.GetNodes()) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(result.GetNodes()))
	}

	for i, node := range result.GetNodes() {
		expected := nodes.GetNodes()[i]

		if node.GetId() != expected.GetId() {
			t.Errorf("id mismatch for node %d: got %d, want %d", i, node.GetId(), expected.GetId())
		}
	}
}

func TestNodeStateRoundTrip(t *testing.T) {
	p := createStreamNodeStateReply(nil)

	// Proto → Model
	model := NodeStateFromProto(p)

	// Model → Proto
	result := model.ToProto()

	// Validações
	if result.GetId() != p.GetId() {
		t.Errorf("id mismatch: got %d, want %d", result.GetId(), p.GetId())
	}

	if result.GetAddr() != p.GetAddr() {
		t.Errorf("addr mismatch: got %s, want %s", result.GetAddr(), p.GetAddr())
	}

	if result.GetTerm() != p.GetTerm() {
		t.Errorf("term mismatch")
	}

	if result.GetCommitIndex() != p.GetCommitIndex() {
		t.Errorf("commitIndex mismatch")
	}

	if result.GetLastLogIndex() != p.GetLastLogIndex() {
		t.Errorf("lastLogIndex mismatch")
	}

	if result.GetLastApplied() != p.GetLastApplied() {
		t.Errorf("lastApplied mismatch")
	}

	if result.GetStatus() != p.GetStatus() {
		t.Errorf("status mismatch")
	}

	if result.GetRole() != p.GetRole() {
		t.Errorf("role mismatch")
	}

	if result.GetLeaderId() != p.GetLeaderId() {
		t.Errorf("leaderId mismatch")
	}
}

func createStreamNodeStateReply(id *int64) *pb.StreamNodeStateReply {
	leaderID := int64(2)

	return &pb.StreamNodeStateReply{
		Id: *id,
		Addr: func() string {
			if id != nil {
				return "localhost:5001"
			}

			return ""
		}(),
		Role:         pb.NodeRole_ROLE_LEADER,
		Term:         3,
		LeaderId:     &leaderID,
		CommitIndex:  10,
		LastLogIndex: 12,
		LastApplied:  10,
		Status:       pb.NodeStatus_STATUS_RUNNING,
	}
}
