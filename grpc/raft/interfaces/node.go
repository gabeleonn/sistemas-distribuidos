package interfaces

import (
	"context"

	pb "raft/autogen"
)

type NodeServer struct {
	pb.UnimplementedRaftNodeServer
}

func (s *NodeServer) ExecuteCommand(
	ctx context.Context,
	req *pb.CommandExecutionArguments,
) (*pb.CommandExecutionReply, error) {
	return nil, nil
}

func (s *NodeServer) AppendEntries(
	ctx context.Context,
	req *pb.AppendEntriesArguments,
) (*pb.AppendEntriesReply, error) {
	return nil, nil
}

func (s *NodeServer) RequestVote(
	ctx context.Context,
	req *pb.RequestVoteArguments,
) (*pb.RequestVoteReply, error) {
	return nil, nil
}
