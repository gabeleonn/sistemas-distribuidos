package models

import (
	pb "raft/autogen"
)

type NodeStatus int

const (
	Stopped NodeStatus = iota
	Running
	Unreachable
)

type Cluster struct {
	Nodes map[int64]Metadata
}

type ClusterState struct {
	Nodes map[int64]NodeState
}

type NodeState struct {
	ID           int64
	Addr         string
	Role         NodeRole
	Term         int64
	LeaderID     *int64
	CommitIndex  int64
	LastLogIndex int64
	LastApplied  int64
	Status       NodeStatus
}

func NodeStateFromProto(p *pb.StreamNodeStateReply) NodeState {
	var leaderID *int64
	if p.LeaderId != nil {
		id := p.GetLeaderId()
		leaderID = &id
	}

	return NodeState{
		ID:           p.GetId(),
		Addr:         p.GetAddr(),
		Role:         NodeRole(p.GetRole()),
		Term:         p.GetTerm(),
		LeaderID:     leaderID,
		CommitIndex:  p.GetCommitIndex(),
		LastLogIndex: p.GetLastLogIndex(),
		LastApplied:  p.GetLastApplied(),
		Status:       NodeStatus(p.GetStatus()),
	}
}

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

func ClusterStateFromProto(p *pb.StreamNodesStatesReply) ClusterState {
	nodes := make(map[int64]NodeState)

	for _, node := range p.GetNodes() {
		nodes[node.GetId()] = NodeStateFromProto(node)
	}

	return ClusterState{
		Nodes: nodes,
	}
}

func (cs ClusterState) ToProto() *pb.StreamNodesStatesReply {
	nodes := make([]*pb.StreamNodeStateReply, 0, len(cs.Nodes))

	for _, node := range cs.Nodes {
		nodes = append(nodes, node.ToProto())
	}

	return &pb.StreamNodesStatesReply{
		Nodes: nodes,
	}
}
