package service

import (
	"context"
	"goraft/proto"
	"goraft/raft"
	"log/slog"
)

type NodeService struct {
	proto.UnimplementedNodeServer

	node        *raft.Node
	logger      *slog.Logger
	onHeartbeat func()
}

func NewNodeService(node *raft.Node, logger *slog.Logger, onHeartbeat func()) *NodeService {
	return &NodeService{
		node:        node,
		logger:      logger,
		onHeartbeat: onHeartbeat,
	}
}

func (s *NodeService) RequestVote(
	ctx context.Context,
	req *proto.RequestVoteRequest,
) (*proto.RequestVoteResponse, error) {
	request := raft.RequestVoteRequestFromProto(req)
	vote := s.node.CandidateResponse(*request)

	return vote.ToProto(), nil
}

func (s *NodeService) AppendEntries(
	ctx context.Context,
	req *proto.AppendEntriesRequest,
) (*proto.AppendEntriesResponse, error) {
	request := raft.AppendEntriesRequestFromProto(req)
	response := s.node.HeartbeatResponse(*request)

	if response.Success {
		s.onHeartbeat()
	}

	return response.ToProto(), nil
}
