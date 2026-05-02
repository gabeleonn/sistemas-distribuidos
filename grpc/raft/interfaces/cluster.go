package interfaces

import (
	"context"

	pb "raft/autogen"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterServer struct {
	pb.UnimplementedRaftClusterServer
}

func (s *ClusterServer) StreamNodeState(
	req *pb.StreamNodeStateArguments,
	stream pb.RaftCluster_StreamNodeStateServer,
) error {
	return nil
}

func (s *ClusterServer) StreamNodesStates(
	req *emptypb.Empty,
	stream pb.RaftCluster_StreamNodesStatesServer,
) error {
	return nil
}

func (s *ClusterServer) AddNode(
	ctx context.Context,
	req *pb.NodeInfoArguments,
) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *ClusterServer) RemoveNode(
	ctx context.Context,
	req *pb.NodeInfoArguments,
) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}