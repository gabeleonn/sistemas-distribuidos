package interfaces

import (
	"context"
	"raft/models/cluster"
	"raft/models/node"

	pb "raft/autogen"

	"google.golang.org/protobuf/types/known/emptypb"
)

// ClusterServer implements the RaftCluster gRPC service, providing methods to stream node states and manage cluster membership.
type ClusterServer struct {
	pb.UnimplementedRaftClusterServer
	State cluster.State
}

// StreamNodeState streams the state of a single node in the cluster to the client.
func (s *ClusterServer) StreamNodeState(
	req *pb.StreamNodeStateArguments,
	stream pb.RaftCluster_StreamNodeStateServer,
) error {
	return nil
}

// StreamNodesStates streams the states of all nodes in the cluster to the client.
func (s *ClusterServer) StreamNodesStates(
	req *emptypb.Empty,
	stream pb.RaftCluster_StreamNodesStatesServer,
) error {
	return nil
}

// AddNode adds a new node to the cluster with the given metadata.
func (s *ClusterServer) AddNode(
	ctx context.Context,
	req *pb.NodeInfoArguments,
) (*emptypb.Empty, error) {
	s.State.Nodes[req.GetId()] = node.NewNode(req.GetId(), req.GetAddr())

	return &emptypb.Empty{}, nil
}

// RemoveNode removes a node from the cluster based on the given metadata.
func (s *ClusterServer) RemoveNode(
	ctx context.Context,
	req *pb.NodeInfoArguments,
) (*emptypb.Empty, error) {
	node := s.State.Nodes[req.GetId()]

	if node != nil {
		node.BeforeRemove()
		delete(s.State.Nodes, node.ID)
	}

	return &emptypb.Empty{}, nil
}
