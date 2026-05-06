package node

import (
	pb "raft/autogen"
)

// Status represents the status of a node in the cluster.
type Status int

// Node statuses
const (
	Stopped Status = iota
	Running
	Unreachable
)

// State represents the state of an individual node in the cluster.
type State struct {
	ID           int64
	Addr         string
	Role         Role
	Term         int64
	LeaderID     *int64
	CommitIndex  int64
	LastLogIndex int64
	LastApplied  int64
	Status       Status
}

// ToProto converts a State to its protobuf representation.
func (s State) ToProto() *pb.StreamNodeStateReply {
	role := pb.NodeRole(s.Role)
	status := pb.NodeStatus(s.Status)

	return &pb.StreamNodeStateReply{
		Id:           s.ID,
		Addr:         s.Addr,
		Role:         role,
		Term:         s.Term,
		LeaderId:     s.LeaderID,
		CommitIndex:  s.CommitIndex,
		LastLogIndex: s.LastLogIndex,
		LastApplied:  s.LastApplied,
		Status:       status,
	}
}

// StateFromProto converts a protobuf StreamNodesStatesReply to a State.
func StateFromProto(p *pb.StreamNodesStatesReply) map[int64]*Node {
	nodes := make(map[int64]*Node)

	for _, node := range p.GetNodes() {
		nodes[node.GetId()] = NewNode(node.GetId(), node.GetAddr())
	}

	return nodes
}
