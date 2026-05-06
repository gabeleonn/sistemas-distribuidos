package cluster

import (
	pb "raft/autogen"
	"raft/models/node"
	"testing"
)

func TestClusterStateRoundTrip(t *testing.T) {
	n1 := node.NewNode(0, "localhost:5001")
	n2 := node.NewNode(1, "localhost:5002")

	model := State{
		Nodes: map[int64]*node.Node{
			0: n1,
			1: n2,
		},
	}

	result := model.ToProto()

	if len(result.GetNodes()) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(result.GetNodes()))
	}

	resultMap := make(map[int64]*pb.StreamNodeStateReply)
	for _, n := range result.GetNodes() {
		resultMap[n.GetId()] = n
	}

	for _, n := range []*node.Node{n1, n2} {
		ns := n.ToNodeState()
		got, ok := resultMap[ns.ID]
		if !ok {
			t.Errorf("node %d not found in result", ns.ID)
			continue
		}
		if got.GetAddr() != ns.Addr {
			t.Errorf("addr mismatch for node %d: got %s, want %s", ns.ID, got.GetAddr(), ns.Addr)
		}
		if got.GetRole() != pb.NodeRole(ns.Role) {
			t.Errorf("role mismatch for node %d", ns.ID)
		}
		if got.GetTerm() != ns.Term {
			t.Errorf("term mismatch for node %d", ns.ID)
		}
		if got.GetCommitIndex() != ns.CommitIndex {
			t.Errorf("commitIndex mismatch for node %d", ns.ID)
		}
		if got.GetStatus() != pb.NodeStatus(ns.Status) {
			t.Errorf("status mismatch for node %d", ns.ID)
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
