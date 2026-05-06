package cluster

import (
	pb "raft/autogen"
	"raft/models/node"
)

// Metadata represents the metadata of a node in the cluster, including its ID and address.
type Metadata struct {
	ID   int64
	Addr string
}

// State represents the current state of the cluster, including the state of each node.
type State struct {
	Nodes map[int64]*node.Node
}

// NodeStateFromProto converts a protobuf StreamNodeStateReply to a NodeState.
func NodeStateFromProto(p *pb.StreamNodeStateReply) node.State {
	var leaderID *int64
	if p.LeaderId != nil {
		id := p.GetLeaderId()
		leaderID = &id
	}

	return node.State{
		ID:           p.GetId(),
		Addr:         p.GetAddr(),
		Role:         node.Role(p.GetRole()),
		Term:         p.GetTerm(),
		LeaderID:     leaderID,
		CommitIndex:  p.GetCommitIndex(),
		LastLogIndex: p.GetLastLogIndex(),
		LastApplied:  p.GetLastApplied(),
		Status:       node.Status(p.GetStatus()),
	}
}

// ToProto converts a State to its protobuf representation.
func (s State) ToProto() *pb.StreamNodesStatesReply {
	nodes := make([]*pb.StreamNodeStateReply, 0, len(s.Nodes))

	for _, node := range s.Nodes {
		ns := node.ToNodeState()
		nodes = append(nodes, ns.ToProto())
	}

	return &pb.StreamNodesStatesReply{
		Nodes: nodes,
	}
}
