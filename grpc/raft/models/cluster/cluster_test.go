package cluster

import (
	pb "raft/autogen"
	"sort"
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

	// Sort both slices by ID before comparing (map iteration is non-deterministic)
	resultNodes := result.GetNodes()
	sort.Slice(resultNodes, func(i, j int) bool {
		return resultNodes[i].GetId() < resultNodes[j].GetId()
	})

	expectedNodes := []*pb.StreamNodeStateReply{p1, p2}
	sort.Slice(expectedNodes, func(i, j int) bool {
		return expectedNodes[i].GetId() < expectedNodes[j].GetId()
	})

	for i, node := range resultNodes {
		expected := expectedNodes[i]

		if node.GetId() != expected.GetId() {
			t.Errorf("id mismatch for node %d: got %d, want %d", i, node.GetId(), expected.GetId())
		}
		if node.GetAddr() != expected.GetAddr() {
			t.Errorf("addr mismatch for node %d: got %s, want %s", i, node.GetAddr(), expected.GetAddr())
		}
		if node.GetRole() != expected.GetRole() {
			t.Errorf("role mismatch for node %d", i)
		}
		if node.GetTerm() != expected.GetTerm() {
			t.Errorf("term mismatch for node %d", i)
		}
		if node.GetCommitIndex() != expected.GetCommitIndex() {
			t.Errorf("commitIndex mismatch for node %d", i)
		}
		if node.GetLastLogIndex() != expected.GetLastLogIndex() {
			t.Errorf("lastLogIndex mismatch for node %d", i)
		}
		if node.GetLastApplied() != expected.GetLastApplied() {
			t.Errorf("lastApplied mismatch for node %d", i)
		}
		if node.GetStatus() != expected.GetStatus() {
			t.Errorf("status mismatch for node %d", i)
		}
	}
}

func TestNodeStateRoundTrip(t *testing.T) {
	id := int64(42)
	p := createStreamNodeStateReply(&id)

	// Proto → Model
	model := NodeStateFromProto(p)

	// Model → Proto
	result := model.ToProto()

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

func TestNodeStateRoundTripWithoutLeader(t *testing.T) {
	id := int64(1)
	p := &pb.StreamNodeStateReply{
		Id:           id,
		Addr:         "localhost:5001",
		Role:         pb.NodeRole_ROLE_FOLLOWER,
		Term:         1,
		LeaderId:     nil,
		CommitIndex:  0,
		LastLogIndex: 0,
		LastApplied:  0,
		Status:       pb.NodeStatus_STATUS_RUNNING,
	}

	model := NodeStateFromProto(p)
	result := model.ToProto()

	if model.LeaderID != nil {
		t.Errorf("expected LeaderID to be nil, got %d", *model.LeaderID)
	}
	if result.LeaderId != nil {
		t.Errorf("expected proto LeaderId to be nil, got %d", *result.LeaderId)
	}
}

func createStreamNodeStateReply(id *int64) *pb.StreamNodeStateReply {
	leaderID := int64(2)

	return &pb.StreamNodeStateReply{
		Id:           *id,
		Addr:         "localhost:5001",
		Role:         pb.NodeRole_ROLE_LEADER,
		Term:         3,
		LeaderId:     &leaderID,
		CommitIndex:  10,
		LastLogIndex: 12,
		LastApplied:  10,
		Status:       pb.NodeStatus_STATUS_RUNNING,
	}
}
