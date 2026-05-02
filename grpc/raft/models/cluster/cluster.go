package cluster

import (
	pb "raft/autogen"
	"raft/models/node"
)

// NodeStatus represents the status of a node in the cluster.
type NodeStatus int

// Node statuses
const (
	Stopped NodeStatus = iota
	Running
	Unreachable
)

// Cluster represents the state of the Raft cluster, including the metadata of all nodes and their current states.
type Cluster struct {
	Nodes map[int64]node.Metadata
}

// State represents the current state of the cluster, including the state of each node.
type State struct {
	Nodes map[int64]NodeState
}

// NodeState represents the state of an individual node in the cluster.
type NodeState struct {
	ID           int64
	Addr         string
	Role         node.Role
	Term         int64
	LeaderID     *int64
	CommitIndex  int64
	LastLogIndex int64
	LastApplied  int64
	Status       NodeStatus
}

// NodeStateFromProto converts a protobuf StreamNodeStateReply to a NodeState.
func NodeStateFromProto(p *pb.StreamNodeStateReply) NodeState {
	var leaderID *int64
	if p.LeaderId != nil {
		id := p.GetLeaderId()
		leaderID = &id
	}

	return NodeState{
		ID:           p.GetId(),
		Addr:         p.GetAddr(),
		Role:         node.Role(p.GetRole()),
		Term:         p.GetTerm(),
		LeaderID:     leaderID,
		CommitIndex:  p.GetCommitIndex(),
		LastLogIndex: p.GetLastLogIndex(),
		LastApplied:  p.GetLastApplied(),
		Status:       NodeStatus(p.GetStatus()),
	}
}

// ToProto converts a NodeState to its protobuf representation.
func (ns NodeState) ToProto() *pb.StreamNodeStateReply {
	role := pb.NodeRole(ns.Role)
	status := pb.NodeStatus(ns.Status)

	return &pb.StreamNodeStateReply{
		Id:           ns.ID,
		Addr:         ns.Addr,
		Role:         role,
		Term:         ns.Term,
		LeaderId:     ns.LeaderID,
		CommitIndex:  ns.CommitIndex,
		LastLogIndex: ns.LastLogIndex,
		LastApplied:  ns.LastApplied,
		Status:       status,
	}
}

// StateFromProto converts a protobuf StreamNodesStatesReply to a State.
func StateFromProto(p *pb.StreamNodesStatesReply) State {
	nodes := make(map[int64]NodeState)

	for _, node := range p.GetNodes() {
		nodes[node.GetId()] = NodeStateFromProto(node)
	}

	return State{
		Nodes: nodes,
	}
}

// ToProto converts a State to its protobuf representation.
func (s State) ToProto() *pb.StreamNodesStatesReply {
	nodes := make([]*pb.StreamNodeStateReply, 0, len(s.Nodes))

	for _, node := range s.Nodes {
		nodes = append(nodes, node.ToProto())
	}

	return &pb.StreamNodesStatesReply{
		Nodes: nodes,
	}
}
